package main

import (
	"compiler/codegen"
	"compiler/lexer"
	"compiler/parser"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("请指定源文件路径")
		return
	}

	sourceCode := readSourceFile(os.Args[1])

	// 词法分析
	l := lexer.NewLexer(sourceCode)


	// 检查词法错误
	if l.HasErrors() {
		fmt.Println("词法分析错误：")
		for _, err := range l.GetErrors() {
			fmt.Println(err)
		}
		return
	}

	// 语法分析
	p := parser.NewParser(l)
	ast, err := p.Parse()
	if err != nil {
		fmt.Println("语法分析错误：")
		fmt.Println(err)
		return
	}

	// 代码生成
	cg := codegen.NewCodeGenerator()
	output := cg.Generate(ast)
	
	// 输出目标代码
	if err := writeOutput("output.asm", output); err != nil {
		fmt.Printf("代码生成错误：%v\n", err)
		return
	}

	fmt.Println("编译成功！输出文件：output.asm")
	fmt.Println("您可以使用emu8086打开并运行此文件")
}

func readSourceFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("读取源文件错误：%v\n", err)
		os.Exit(1)
	}
	return string(data)
}

func writeOutput(path string, code []string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("创建输出文件失败：%v", err)
	}
	defer f.Close()

	for _, line := range code {
		if _, err := fmt.Fprintln(f, line); err != nil {
			return fmt.Errorf("写入输出文件失败：%v", err)
		}
	}
	return nil
}
