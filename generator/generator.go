package generator

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
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

	stepDefinitions []struct {
		StepDefinition string
		FuncName       string
	}

	beforeFeatureFuncName  string
	afterFeatureFuncName   string
	beforeScenarioFuncName string
	afterScenarioFuncName  string
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
	g.findStepDefinitionFuncs()
	g.findHandlersFunc()

	fmt.Println(g.stepDefinitions)

	runnerFuncDecl := g.findTestRunnerFuncDecl()
	fmt.Println(runnerFuncDecl)

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

	ast.Print(nil, g.astFile)

	outFile, err := os.OpenFile(g.filepath, os.O_WRONLY, 0666)
	defer outFile.Close()

	if err != nil {
		panic("oops")
	}

	printer.Fprint(outFile, token.NewFileSet(), g.astFile)
}

func (g *Generator) findStepDefinitionFuncs() {
	decls := g.astFile.Decls

	funcDeclVisitor(decls, func(fd *ast.FuncDecl) {
		if !isStepDefinitionFunc(fd) {
			return
		}

		comment, err := funcDeclGetComment(fd, isSwedeStepDefinitionComment)

		if err != nil {
			panic(err)
		}

		stepDefinition := parseSwedeStepDefinitionFromComment(comment)

		g.stepDefinitions = append(g.stepDefinitions, struct {
			StepDefinition string
			FuncName       string
		}{StepDefinition: stepDefinition, FuncName: fd.Name.String()})
	})
}

func isStepDefinitionFunc(funcDecl *ast.FuncDecl) bool {
	return funcDeclHasComment(funcDecl, isSwedeStepDefinitionComment)
}

func isSwedeStepDefinitionComment(comment string) bool {
	return strings.Contains(comment, "swede:step")
}

func parseSwedeStepDefinitionFromComment(comment string) string {
	stepDefinition := strings.ReplaceAll(comment, "swede:step", "")
	stepDefinition = strings.ReplaceAll(stepDefinition, "//", "")
	return strings.TrimSpace(stepDefinition)
}

func (g *Generator) findHandlersFunc() {
	decls := g.astFile.Decls

	hasSwedeBeforeScenarioComment := func(comment string) bool {
		return strings.Contains(comment, "swede:beforeScenario")
	}

	hasSwedeAfterScenarioComment := func(comment string) bool {
		return strings.Contains(comment, "swede:afterScenario")
	}

	hasSwedeBeforeFeatureComment := func(comment string) bool {
		return strings.Contains(comment, "swede:beforeFeature")
	}

	hasSwedeAfterFeatureComment := func(comment string) bool {
		return strings.Contains(comment, "swede:afterFeature")
	}

	funcDeclVisitor(decls, func(fd *ast.FuncDecl) {
		funcName := fd.Name.String()

		if funcDeclHasComment(fd, hasSwedeBeforeScenarioComment) {
			g.beforeScenarioFuncName = funcName
		}

		if funcDeclHasComment(fd, hasSwedeAfterScenarioComment) {
			g.afterScenarioFuncName = funcName
		}

		if funcDeclHasComment(fd, hasSwedeBeforeFeatureComment) {
			g.beforeFeatureFuncName = funcName

		}

		if funcDeclHasComment(fd, hasSwedeAfterFeatureComment) {
			g.afterFeatureFuncName = funcName
		}
	})
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
	updateFunc := func(c *astutil.Cursor) bool {
		node := c.Node()
		funcNode, ok := node.(*ast.FuncDecl)

		if !ok {
			return true
		}

		if funcNode.Name.String() != testRunnerFunctionName {
			return true
		}

		c.Replace(funcDecl)
		fmt.Println("HELLO WORLD")
		return true
	}

	astutil.Apply(g.astFile, updateFunc, nil)
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
		"featureFiles":           g.featureFiles,
		"stepDefinitions":        g.stepDefinitions,
		"beforeFeatureFuncName":  g.beforeFeatureFuncName,
		"afterFeatureFuncName":   g.afterFeatureFuncName,
		"beforeScenarioFuncName": g.beforeScenarioFuncName,
		"afterScenarioFuncName":  g.afterScenarioFuncName,
	}

	buf := &bytes.Buffer{}

	t := template.New("t")
	t = template.Must(t.Parse(testRunnerTemplate))
	t.Execute(buf, templateData)

	fmt.Println(buf)
	fmt.Println(g.stepDefinitions)

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
