package generator

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"text/template"

	"golang.org/x/tools/go/ast/astutil"
)

type Generator struct {
	filepath string

	featureFiles []string

	//parser data
	fset    *token.FileSet
	astFile *ast.File

	stepDefinitions []stepDefinitionWithFuncName
}

type GeneratorOptions struct {
	FeatureFiles []string
}

func NewGenerator(options GeneratorOptions) *Generator {
	g := &Generator{
		featureFiles: options.FeatureFiles,
	}

	g.filepath = getProcessedFilePath()
	g.parseSourceFile()

	return g
}

func (g *Generator) Generate() {
	g.stepDefinitions = g.findStepDefinitionFuncs()

	fmt.Println(g.stepDefinitions)

	runnerFuncDecl := g.findTestRunnerFuncDecl()
	newRunnerFuncDecl := g.generateTestRunnerFuncDecl()

	if runnerFuncDecl == nil {
		g.insertTestRunnerFuncDecl(newRunnerFuncDecl)
	} else {
		fmt.Println("found test runner")
		g.updateTestRunnerFuncDecl(newRunnerFuncDecl)
	}

	g.addRequiredImports()

	g.saveGeneratedFile()
}

func (g *Generator) addRequiredImports() {
	astutil.AddImport(g.fset, g.astFile, "testing")
	astutil.AddImport(g.fset, g.astFile, "me.weldnor/swede/runner")
}

func (g *Generator) parseSourceFile() {
	g.fset = token.NewFileSet()

	astFile, err := parser.ParseFile(g.fset, getProcessedFilePath(), nil, parser.ParseComments)

	if err != nil {
		panic(errors.New("can't parse file"))
	}

	g.astFile = astFile
}

func (g *Generator) saveGeneratedFile() {
	outFile, err := os.OpenFile(g.filepath, os.O_WRONLY, 0666)
	defer outFile.Close()

	if err != nil {
		panic("oops")
	}

	if err := format.Node(outFile, g.fset, g.astFile); err != nil {
		panic(err)
	}
}

type stepDefinitionWithFuncName struct {
	StepDefinition string
	FuncName       string
}

func (g *Generator) findStepDefinitionFuncs() []stepDefinitionWithFuncName {
	result := make([]stepDefinitionWithFuncName, 0)

	for _, decl := range g.astFile.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)

		if !ok {
			continue
		}

		if !g.checkStepFuncSignature(funcDecl) {
			continue
		}

		stepDefinitionComments := g.getStepCommentsFromFuncDecl(funcDecl)

		//todo > 1?
		if len(stepDefinitionComments) != 1 {
			continue
		}

		stepDefinitionComment := stepDefinitionComments[0]
		stepDefinition := g.parseStepDefinitionFromComment(stepDefinitionComment)
		funcName := funcDecl.Name.String()

		result = append(result, stepDefinitionWithFuncName{
			StepDefinition: stepDefinition,
			FuncName:       funcName,
		})
	}

	return result
}

func (g *Generator) parseStepDefinitionFromComment(commentText string) string {
	stepDefinition := strings.ReplaceAll(commentText, "swede:step", "")
	stepDefinition = strings.ReplaceAll(stepDefinition, "//", "")

	return strings.TrimSpace(stepDefinition)
}

func (g *Generator) checkStepFuncSignature(funcDecl *ast.FuncDecl) bool {
	//todo implement
	return true
}

func (g *Generator) getStepCommentsFromFuncDecl(funcDecl *ast.FuncDecl) []string {
	stepDefinitionComments := make([]string, 0)

	if funcDecl.Doc == nil {
		return nil
	}

	for _, comment := range funcDecl.Doc.List {
		if strings.Contains(comment.Text, "swede:step") {
			stepDefinitionComments = append(stepDefinitionComments, comment.Text)
		}
	}

	return stepDefinitionComments
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
    
    // Load feature files
{{range $key, $value := .featureFiles}} 
    _runner.LoadFeatureFile("{{$value}}")
{{end}}

    // Register step functions
{{range $key, $value := .stepDefinitions}} 
    _runner.RegisterFunc("{{$value.StepDefinition}}", {{$value.FuncName}})
{{end}}
    _runner.Run()
}

`

func (g *Generator) generateTestRunnerFuncDecl() *ast.FuncDecl {
	templateData := map[string]interface{}{
		"featureFiles":    g.featureFiles,
		"stepDefinitions": g.stepDefinitions,
	}

	buf := &bytes.Buffer{}

	t := template.New("t")
	t = template.Must(t.Parse(testRunnerTemplate))
	t.Execute(buf, templateData)

	templateAst, err := parser.ParseFile(
		token.NewFileSet(),
		//Источник для парсинга лежит не в файле,
		"",
		buf,
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
