package generator

import (
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path"
	"strings"
)

type Generator struct {
	filepath string
	//parser data
	fset    *token.FileSet
	astFile *ast.File
}

func NewGenerator() *Generator {
	g := &Generator{}

	g.filepath = getFilePath()
	g.parseFile()

	return g
}

func getFilePath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(errors.New("can't get current working directory"))
	}

	return path.Join(wd, os.Getenv("GOFILE"))
}

func (g *Generator) parseFile() {
	g.fset = token.NewFileSet()

	astFile, err := parser.ParseFile(g.fset, getFilePath(), nil, parser.ParseComments)

	if err != nil {
		panic(errors.New("can't parse file"))
	}

	g.astFile = astFile
}

func (g *Generator) Generate() {
	g.findStepDefinitionFuncs()

	runnerFuncDecl := g.findTestRunnerFuncDecl()
	newRunnerFuncDecl := g.generateTestRunnerFuncDecl()

	if runnerFuncDecl == nil {
		g.insertTestRunnerFuncDecl(newRunnerFuncDecl)
	} else {
		fmt.Println("found test runner")
		g.updateTestRunnerFuncDecl(newRunnerFuncDecl)
	}

	g.saveToFile()
}

func (g *Generator) saveToFile() {
	outFile, err := os.OpenFile(g.filepath, os.O_WRONLY, 0666)
	if err != nil {
		panic("oops")
	}
	//Не забываем прибраться
	defer outFile.Close()

	if err := format.Node(outFile, g.fset, g.astFile); err != nil {
		panic(err)
	}
}

func (g *Generator) findStepDefinitionFuncs() {
	for _, decl := range g.astFile.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		fmt.Println("scanning function " + funcDecl.Name.String() + "...")

		if g.checkFunction(funcDecl) {
			fmt.Println("found step definition function: " + funcDecl.Name.String())
		}
	}
}

func (g *Generator) checkFunction(funcDecl *ast.FuncDecl) bool {
	if funcDecl.Doc == nil {
		return false
	}

	comments := funcDecl.Doc.List

	for _, comment := range comments {
		if checkComment(comment.Text) {
			return true
		}
	}
	return false
}

func checkComment(comment string) bool {
	if strings.Contains(comment, "swede:step") {
		return true
	}
	return false
}

func (g *Generator) isTestRunnerFuncDecl(decl ast.Decl) bool {
	funcDecl, ok := decl.(*ast.FuncDecl)

	if !ok {
		return false
	}

	return funcDecl.Name.String() == testRunnerFunctionName
}

func (g *Generator) insertTestRunnerFuncDecl(funcDecl *ast.FuncDecl) {
	g.astFile.Decls = append(g.astFile.Decls, funcDecl)
}

func (g *Generator) updateTestRunnerFuncDecl(funcDecl *ast.FuncDecl) {
	for i, decl := range g.astFile.Decls {
		if g.isTestRunnerFuncDecl(decl) {
			g.astFile.Decls[i] = funcDecl
			return
		}
	}
}

const testRunnerFunctionName = "TestSwedeRunner"

func (g *Generator) findTestRunnerFuncDecl() *ast.FuncDecl {
	for _, decl := range g.astFile.Decls {
		if g.isTestRunnerFuncDecl(decl) {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				panic("oops")
			}

			return funcDecl
		}
	}
	return nil
}

const testRunnerTemplate = `
package main

import(
    "testing"
    "me.weldnor/swede/runner"
)

func TestSwedeRunner(t *testing.T) {
    _runner := runner.NewRunner()
	_runner.LoadFeatureFile("./feature/formatted.swede")
}

`

func (g *Generator) generateTestRunnerFuncDecl() *ast.FuncDecl {

	templateAst, err := parser.ParseFile(
		token.NewFileSet(),
		//Источник для парсинга лежит не в файле,
		"",
		[]byte(testRunnerTemplate),
		parser.ParseComments,
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(templateAst.Decls)

	decl := templateAst.Decls[1] //skip import declaration
	funcDecl, ok := decl.(*ast.FuncDecl)

	if !ok {
		panic("oops")
	}

	return funcDecl
}
