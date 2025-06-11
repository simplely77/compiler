package codegen

import (
	"compiler/parser"
	"fmt"
)

type CodeGenerator struct {
	code    []string
	varCount int
	varMap  map[string]string
	labelCount int
}

func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{
		code:    make([]string, 0),
		varMap:  make(map[string]string),
		labelCount: 0,
	}
}

func (cg *CodeGenerator) Generate(ast *parser.AST) []string {
	cg.code = append(cg.code,
		"section .data",
		"    format db '%d', 10, 0",
		"    input_format db '%d', 0",
		"",
		"section .bss",
	)
	// 声明所有变量
	cg.collectVars(ast)
	for name := range cg.varMap {
		cg.code = append(cg.code, fmt.Sprintf("    %s resq 1", name))
	}
	cg.code = append(cg.code,
		"",
		"section .text",
		"    global main",
		"    extern printf",
		"    extern scanf",
		"",
		"main:",
		"    sub rsp, 32",
	)
	// 生成所有语句
	for _, stmt := range ast.Statements {
		cg.genStatement(stmt)
	}
	cg.code = append(cg.code,
		"    add rsp, 32",
		"    xor rax, rax",
		"    ret",
	)
	return cg.code
}

func (cg *CodeGenerator) collectVars(ast *parser.AST) {
	for _, stmt := range ast.Statements {
		switch s := stmt.(type) {
		case *parser.Assignment:
			cg.varMap[s.Ident] = s.Ident
			cg.collectVarsFromExpr(s.Value)
		case *parser.PrintStatement:
			cg.collectVarsFromExpr(s.Expr)
		case *parser.InputStatement:
			cg.varMap[s.Ident] = s.Ident
		case *parser.IfStatement:
			cg.collectVarsFromExpr(s.Condition)
			for _, stmt := range s.Then {
				cg.collectVarsFromStatement(stmt)
			}
			for _, stmt := range s.Else {
				cg.collectVarsFromStatement(stmt)
			}
		case *parser.WhileStatement:
			cg.collectVarsFromExpr(s.Condition)
			for _, stmt := range s.Body {
				cg.collectVarsFromStatement(stmt)
			}
		}
	}
}

func (cg *CodeGenerator) collectVarsFromExpr(expr parser.Expr) {
	switch e := expr.(type) {
	case *parser.IdentExpr:
		cg.varMap[e.Name] = e.Name
	case *parser.BinaryExpr:
		cg.collectVarsFromExpr(e.Left)
		cg.collectVarsFromExpr(e.Right)
	case *parser.ComparisonExpr:
		cg.collectVarsFromExpr(e.Left)
		cg.collectVarsFromExpr(e.Right)
	}
}

func (cg *CodeGenerator) collectVarsFromStatement(stmt parser.Statement) {
	switch s := stmt.(type) {
	case *parser.Assignment:
		cg.varMap[s.Ident] = s.Ident
		cg.collectVarsFromExpr(s.Value)
	case *parser.PrintStatement:
		cg.collectVarsFromExpr(s.Expr)
	case *parser.InputStatement:
		cg.varMap[s.Ident] = s.Ident
	}
}

func (cg *CodeGenerator) genStatement(stmt parser.Statement) {
	switch s := stmt.(type) {
	case *parser.Assignment:
		cg.genAssignment(s)
	case *parser.PrintStatement:
		cg.genPrint(s)
	case *parser.InputStatement:
		cg.genInput(s)
	case *parser.IfStatement:
		cg.genIf(s)
	case *parser.WhileStatement:
		cg.genWhile(s)
	}
}

func (cg *CodeGenerator) genAssignment(a *parser.Assignment) {
	cg.genExpr(a.Value, "rax")
	cg.code = append(cg.code, fmt.Sprintf("    mov [rel %s], rax", a.Ident))
}

func (cg *CodeGenerator) genPrint(p *parser.PrintStatement) {
	cg.genExpr(p.Expr, "rax")
	cg.code = append(cg.code,
		"    mov rcx, format",
		"    mov rdx, rax",
		"    call printf",
	)
}

func (cg *CodeGenerator) genInput(i *parser.InputStatement) {
	cg.code = append(cg.code,
		"    mov rcx, input_format",
		fmt.Sprintf("    lea rdx, [rel %s]", i.Ident),
		"    call scanf",
	)
}

func (cg *CodeGenerator) genIf(i *parser.IfStatement) {
	cg.genExpr(i.Condition, "rax")
	elseLabel := cg.newLabel()
	endLabel := cg.newLabel()
	
	cg.code = append(cg.code,
		"    test rax, rax",
		fmt.Sprintf("    jz %s", elseLabel),
	)
	
	for _, stmt := range i.Then {
		cg.genStatement(stmt)
	}
	
	cg.code = append(cg.code,
		fmt.Sprintf("    jmp %s", endLabel),
		fmt.Sprintf("%s:", elseLabel),
	)
	
	for _, stmt := range i.Else {
		cg.genStatement(stmt)
	}
	
	cg.code = append(cg.code, fmt.Sprintf("%s:", endLabel))
}

func (cg *CodeGenerator) genWhile(w *parser.WhileStatement) {
	startLabel := cg.newLabel()
	endLabel := cg.newLabel()
	
	cg.code = append(cg.code, fmt.Sprintf("%s:", startLabel))
	
	cg.genExpr(w.Condition, "rax")
	cg.code = append(cg.code,
		"    test rax, rax",
		fmt.Sprintf("    jz %s", endLabel),
	)
	
	for _, stmt := range w.Body {
		cg.genStatement(stmt)
	}
	
	cg.code = append(cg.code,
		fmt.Sprintf("    jmp %s", startLabel),
		fmt.Sprintf("%s:", endLabel),
	)
}

func (cg *CodeGenerator) genExpr(expr parser.Expr, target string) {
	switch e := expr.(type) {
	case *parser.NumberExpr:
		cg.code = append(cg.code, fmt.Sprintf("    mov %s, %s", target, e.Value))
	case *parser.IdentExpr:
		cg.code = append(cg.code, fmt.Sprintf("    mov %s, [rel %s]", target, e.Name))
	case *parser.BooleanExpr:
		if e.Value {
			cg.code = append(cg.code, fmt.Sprintf("    mov %s, 1", target))
		} else {
			cg.code = append(cg.code, fmt.Sprintf("    mov %s, 0", target))
		}
	case *parser.BinaryExpr:
		cg.genExpr(e.Left, target)
		cg.code = append(cg.code, "    push rax")
		cg.genExpr(e.Right, "rax")
		cg.code = append(cg.code, "    pop rbx")
		switch e.Op {
		case "+":
			cg.code = append(cg.code, "    add rax, rbx")
		case "-":
			cg.code = append(cg.code, "    sub rbx, rax", "    mov rax, rbx")
		case "*":
			cg.code = append(cg.code, "    imul rax, rbx")
		}
	case *parser.ComparisonExpr:
		cg.genExpr(e.Left, "rax")
		cg.code = append(cg.code, "    push rax")
		cg.genExpr(e.Right, "rax")
		cg.code = append(cg.code, "    pop rbx")
		switch e.Op {
		case "==":
			cg.code = append(cg.code,
				"    cmp rbx, rax",
				"    sete al",
				"    movzx rax, al",
			)
		case "!=":
			cg.code = append(cg.code,
				"    cmp rbx, rax",
				"    setne al",
				"    movzx rax, al",
			)
		case "<":
			cg.code = append(cg.code,
				"    cmp rbx, rax",
				"    setl al",
				"    movzx rax, al",
			)
		case ">":
			cg.code = append(cg.code,
				"    cmp rbx, rax",
				"    setg al",
				"    movzx rax, al",
			)
		case "<=":
			cg.code = append(cg.code,
				"    cmp rbx, rax",
				"    setle al",
				"    movzx rax, al",
			)
		case ">=":
			cg.code = append(cg.code,
				"    cmp rbx, rax",
				"    setge al",
				"    movzx rax, al",
			)
		}
		cg.code = append(cg.code, "    xor rax, 1")  // 反转结果
	}
}

func (cg *CodeGenerator) newLabel() string {
	label := fmt.Sprintf("label_%d", cg.labelCount)
	cg.labelCount++
	return label
}
