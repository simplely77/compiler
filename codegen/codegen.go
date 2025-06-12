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
	// Initial COM header and jump to main execution
	cg.code = append(cg.code,
		"#make_COM#",
		"ORG 100h",
		"",
		"jmp main_start", // Jump over data and procedures
		"",
		".DATA",
		"    msg_input db 'Enter a number: $'",
		"    msg_output db 'Result: $'",
		"    msg_div_by_zero db 'Error: Division by zero!$'",
		"    newline db 13, 10, '$'",
	)

	// Declare variables in .DATA section
	cg.collectVars(ast)
	for name := range cg.varMap {
		cg.code = append(cg.code, fmt.Sprintf("    %s dw 0", name))
	}
	cg.code = append(cg.code, "")

	// Add helper procedures to the .CODE section before main_start
	cg.code = append(cg.code, ".CODE")
	cg.AddHelperFunctions()

	cg.code = append(cg.code,
		"main_start:", // Actual start of the main program logic
		"    mov ax, @data",
		"    mov ds, ax",
	)

	// Generate all statements (main program logic)
	for _, stmt := range ast.Statements {
		cg.genStatement(stmt)
	}

	cg.code = append(cg.code,
		"    mov ah, 4Ch", // Program exit
		"    int 21h",
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
	cg.genExpr(a.Value, "ax")
	cg.code = append(cg.code, fmt.Sprintf("    mov %s, ax", a.Ident))
}

func (cg *CodeGenerator) genPrint(p *parser.PrintStatement) {
	cg.genExpr(p.Expr, "ax")
	cg.code = append(cg.code,
		"    push ax",
		"    mov dx, offset msg_output",
		"    mov ah, 9",
		"    int 21h",
		"    pop ax",
		"    call print_number",
		"    mov dx, offset newline",
		"    mov ah, 9",
		"    int 21h",
	)
}

func (cg *CodeGenerator) genInput(i *parser.InputStatement) {
	cg.code = append(cg.code,
		"    mov dx, offset msg_input",
		"    mov ah, 9",
		"    int 21h",
		"    call read_number",
		fmt.Sprintf("    mov %s, ax", i.Ident),
		"    mov dx, offset newline",
		"    mov ah, 9",
		"    int 21h",
	)
}

func (cg *CodeGenerator) genIf(i *parser.IfStatement) {
    cg.genExpr(i.Condition, "ax")
    elseLabel := cg.newLabel()
    endLabel := cg.newLabel()

    cg.code = append(cg.code,
        "    cmp ax, 0",
        fmt.Sprintf("    je %s", elseLabel),
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

	cg.genExpr(w.Condition, "ax")
	cg.code = append(cg.code,
		"    cmp ax, 0",
		fmt.Sprintf("    je %s", endLabel),
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
		cg.code = append(cg.code, fmt.Sprintf("    mov %s, %s", target, e.Name))
	case *parser.BooleanExpr:
		if e.Value {
			cg.code = append(cg.code, fmt.Sprintf("    mov %s, 1", target))
		} else {
			cg.code = append(cg.code, fmt.Sprintf("    mov %s, 0", target))
		}
	case *parser.BinaryExpr:
		cg.genExpr(e.Left, target)
		cg.code = append(cg.code, "    push ax")
		cg.genExpr(e.Right, "ax")
		cg.code = append(cg.code, "    pop bx")
		switch e.Op {
		case "+":
			cg.code = append(cg.code, "    add ax, bx")
		case "-":
			cg.code = append(cg.code, "    sub bx, ax", "    mov ax, bx")
		case "*":
			cg.code = append(cg.code, "    mul bx")
		case "/":
			cg.genExpr(e.Left, target)
			cg.code = append(cg.code, "    push ax")
			cg.genExpr(e.Right, "ax")
			cg.code = append(cg.code, "    pop bx")

			// 现在：AX = 右操作数 (除数), BX = 左操作数 (被除数)
			// 我们想计算 BX / AX (被除数 / 除数)
			// IDIV 指令期望被除数在 AX (或 DX:AX) 中，除数作为其操作数。
			cg.code = append(cg.code,
				"    mov cx, ax", // 将除数 (右操作数) 从 AX 移动到 CX
				"    mov ax, bx", // 将被除数 (左操作数) 从 BX 移动到 AX
			)

			divOkLabel := cg.newLabel()
			cg.code = append(cg.code,
				"    cmp cx, 0",    // 检查除数 (现在在 CX 中) 是否为0
				fmt.Sprintf("    jne %s", divOkLabel),   // 如果不为0，继续执行
				"    mov dx, offset msg_div_by_zero", // 除数为0，显示错误信息
				"    mov ah, 9",
				"    int 21h",
				"    mov ah, 4Ch",  // 程序退出
				"    int 21h",
				fmt.Sprintf("%s:", divOkLabel),
				"    cwd",          // 符号扩展 AX (被除数) 到 DX:AX
				"    idiv cx",      // 有符号除法 DX:AX / CX (除数)
			)
		}
	case *parser.ComparisonExpr:
		cg.genExpr(e.Left, "ax")
		cg.code = append(cg.code, "    push ax")
		cg.genExpr(e.Right, "ax")
		cg.code = append(cg.code, "    pop bx")
		switch e.Op {
		case "==":
			label := cg.newLabel()
			cg.code = append(cg.code,
				"    cmp bx, ax",
				"    mov ax, 0",
				fmt.Sprintf("    jne %s", label),
				"    mov ax, 1",
				fmt.Sprintf("%s:", label),
			)
		case "!=":
			label := cg.newLabel()
			cg.code = append(cg.code,
				"    cmp bx, ax",
				"    mov ax, 0",
				fmt.Sprintf("    je %s", label),
				"    mov ax, 1",
				fmt.Sprintf("%s:", label),
			)
		case "<":
			label := cg.newLabel()
			cg.code = append(cg.code,
				"    cmp bx, ax",
				"    mov ax, 0",
				fmt.Sprintf("    jge %s", label),
				"    mov ax, 1",
				fmt.Sprintf("%s:", label),
			)
		case ">":
			label := cg.newLabel()
			cg.code = append(cg.code,
				"    cmp bx, ax",
				"    mov ax, 0",
				fmt.Sprintf("    jle %s", label),
				"    mov ax, 1",
				fmt.Sprintf("%s:", label),
			)
		case "<=":
			label := cg.newLabel()
			cg.code = append(cg.code,
				"    cmp bx, ax",
				"    mov ax, 0",
				fmt.Sprintf("    jg %s", label),
				"    mov ax, 1",
				fmt.Sprintf("%s:", label),
			)
		case ">=":
			label := cg.newLabel()
			cg.code = append(cg.code,
				"    cmp bx, ax",
				"    mov ax, 0",
				fmt.Sprintf("    jl %s", label),
				"    mov ax, 1",
				fmt.Sprintf("%s:", label),
			)
		}
	}
}

func (cg *CodeGenerator) newLabel() string {
	label := fmt.Sprintf("label_%d", cg.labelCount)
	cg.labelCount++
	return label
}

// AddHelperFunctions 添加辅助函数
func (cg *CodeGenerator) AddHelperFunctions() {
	cg.code = append(cg.code,
		"print_number PROC",
		"    push ax",
		"    push bx",
		"    push cx",
		"    push dx",
		"    push si",
		"    mov si, 0",
		"    cmp ax, 0",
		"    jge print_positive",
		"    mov si, 1",
		"    neg ax",
	"print_positive:",
		"    mov bx, 10",
		"    mov cx, 0",
	"print_number_loop:",
		"    mov dx, 0",
		"    div bx",
		"    push dx",
		"    inc cx",
		"    test ax, ax",
		"    jnz print_number_loop",
		"    cmp si, 1",
		"    je print_minus",
		"    jmp print_number_output",
	"print_minus:",
		"    mov ah, 2",
		"    mov dl, '-'",
		"    int 21h",
	"print_number_output:",
		"    pop dx",
		"    add dl, '0'",
		"    mov ah, 2",
		"    int 21h",
		"    loop print_number_output",
		"    pop si",
		"    pop dx",
		"    pop cx",
		"    pop bx",
		"    pop ax",
		"    ret",
		"print_number ENDP",
		"",
		"read_number PROC",
		"    push bx",
		"    push cx",
		"    push dx",
		"    mov bx, 0",
		"    mov cx, 0",
		"    mov si, 0",
	"read_number_loop:",
		"    mov ah, 1",
		"    int 21h",
		"    cmp al, 13",
		"    je read_number_done",
		"    cmp al, '-'",
		"    je handle_negative",
		"    cmp al, '0'",
		"    jb read_number_loop",
		"    cmp al, '9'",
		"    ja read_number_loop",
		"    sub al, '0'",
		"    mov cl, al",
		"    mov ax, bx",
		"    mov bx, 10",
		"    mul bx",
		"    mov bx, ax",
		"    add bx, cx",
		"    jmp read_number_loop",
	"handle_negative:",
		"    mov si, 1",
		"    jmp read_number_loop",
	"read_number_done:",
		"    cmp si, 1",
		"    jne not_negative",
		"    neg bx",
	"not_negative:",
		"    mov ax, bx",
		"    pop dx",
		"    pop cx",
		"    pop bx",
		"    ret",
		"read_number ENDP",
	)
}
