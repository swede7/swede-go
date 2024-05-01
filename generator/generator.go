package generator

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/dstutil"
	"go/token"
	"strings"
	"text/template"
)

type Generator struct {
	featureFiles []string
	source       string

	//parser data
	dstFile *dst.File

	stepDefinitions []struct {
		StepDefinition string
		FuncName       string
	}

	beforeFeatureFuncName  string
	afterFeatureFuncName   string
	beforeScenarioFuncName string
	afterScenarioFuncName  string

	debug bool
}

type Options struct {
	FeatureFiles []string
	Source       string
	Debug        bool
}

func NewGenerator(options Options) *Generator {
	g := &Generator{
		featureFiles: options.FeatureFiles,
		source:       options.Source,
	}

	g.parseSourceFile()
	return g
}

func (g *Generator) Generate() string {
	g.findStepDefinitionFuncs()
	g.findHandlersFunc()

	g.printDebugMessage(g.stepDefinitions)

	runnerFuncDecl := g.findTestRunnerFuncDecl()
	g.printDebugMessage(runnerFuncDecl)

	newRunnerFuncDecl := g.generateTestRunnerFuncDecl()

	if runnerFuncDecl == nil {
		g.insertTestRunnerFuncDecl(newRunnerFuncDecl)
	} else {
		g.printDebugMessage("found test runner")
		g.updateTestRunnerFuncDecl(newRunnerFuncDecl)
	}

	g.addRequiredImports()

	return g.formatGeneratedFile()
}

func (g *Generator) addRequiredImports() {
	if g.containsTestingModule() {
		return
	}

	inserted := false

	testingImportSpec := &dst.ImportSpec{
		Path: &dst.BasicLit{
			Kind:  token.STRING,
			Value: "\"testing\"",
		},
	}

	newDstFile := dstutil.Apply(g.dstFile, nil, func(cursor *dstutil.Cursor) bool {
		if inserted {
			return false
		}

		_, ok := cursor.Node().(*dst.ImportSpec)

		if !ok {
			return true
		}

		cursor.InsertAfter(testingImportSpec)
		inserted = true

		return true
	})

	g.dstFile = newDstFile.(*dst.File)
	g.dstFile.Imports = append(g.dstFile.Imports, testingImportSpec)
}

func (g *Generator) containsTestingModule() bool {
	for _, importSpec := range g.dstFile.Imports {
		if importSpec.Path.Value == "testing" {
			return true
		}
	}
	return false
}

func (g *Generator) parseSourceFile() {
	g.printDebugMessage("parsing source file...")

	dstFile, err := decorator.Parse(g.source)

	if err != nil {
		panic(errors.New("can't parse file"))
	}

	g.dstFile = dstFile
}

func (g *Generator) formatGeneratedFile() string {
	buf := bytes.Buffer{}

	err := decorator.Fprint(&buf, g.dstFile)

	if err != nil {
		panic(err)
	}

	return buf.String()
}

func (g *Generator) findStepDefinitionFuncs() {
	decls := g.dstFile.Decls

	funcDeclVisitor(decls, func(fd *dst.FuncDecl) {
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

func isStepDefinitionFunc(funcDecl *dst.FuncDecl) bool {
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
	decls := g.dstFile.Decls

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

	funcDeclVisitor(decls, func(fd *dst.FuncDecl) {
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

func (g *Generator) getStepCommentsFromFuncDecl(funcDecl *dst.FuncDecl) []string {
	stepDefinitionComments := make([]string, 0)

	if funcDecl.Decs.Start == nil {
		return nil
	}

	for _, comment := range funcDecl.Decs.Start.All() {
		if strings.Contains(comment, "swede:step") {
			stepDefinitionComments = append(stepDefinitionComments, comment)
		}
	}

	return stepDefinitionComments
}

func (g *Generator) isTestRunnerFuncDecl(decl dst.Decl) bool {
	funcDecl, ok := decl.(*dst.FuncDecl)

	if !ok {
		return false
	}

	return funcDecl.Name.String() == testRunnerFunctionName
}

func (g *Generator) insertTestRunnerFuncDecl(funcDecl *dst.FuncDecl) {

	g.dstFile.Decls = append(g.dstFile.Decls, funcDecl)
}

func (g *Generator) updateTestRunnerFuncDecl(funcDecl *dst.FuncDecl) {
	updateFunc := func(c *dstutil.Cursor) bool {
		node := c.Node()
		funcNode, ok := node.(*dst.FuncDecl)

		if !ok {
			return true
		}

		if funcNode.Name.String() != testRunnerFunctionName {
			return true
		}

		c.Replace(funcDecl)
		g.printDebugMessage("HELLO WORLD")
		return true
	}

	dstutil.Apply(g.dstFile, updateFunc, nil)
}

const testRunnerFunctionName = "TestSwedeRunner"

func (g *Generator) findTestRunnerFuncDecl() *dst.FuncDecl {
	for _, decl := range g.dstFile.Decls {
		if g.isTestRunnerFuncDecl(decl) {
			funcDecl, ok := decl.(*dst.FuncDecl)
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
{{range $key, $value := .featureFiles}} 
    _runner.LoadFeatureFile("{{$value}}")
{{end}}
{{range $key, $value := .stepDefinitions}} 
    _runner.RegisterFunc("{{$value.StepDefinition}}", {{$value.FuncName}})
{{end}}
    _runner.Run()
}

`

func (g *Generator) generateTestRunnerFuncDecl() *dst.FuncDecl {
	templateData := map[string]interface{}{
		"featureFiles":           g.featureFiles,
		"stepDefinitions":        g.stepDefinitions,
		"beforeFeatureFuncName":  g.beforeFeatureFuncName,
		"afterFeatureFuncName":   g.afterFeatureFuncName,
		"beforeScenarioFuncName": g.beforeScenarioFuncName,
		"afterScenarioFuncName":  g.afterScenarioFuncName,
	}

	buf := &bytes.Buffer{}

	t := template.New("")
	t = template.Must(t.Parse(testRunnerTemplate))

	err := t.Execute(buf, templateData)
	if err != nil {
		panic(err)
	}

	g.printDebugMessage(buf)
	g.printDebugMessage(g.stepDefinitions)

	templateAst, err := decorator.Parse(buf.String())

	if err != nil {
		panic(err)
	}

	g.printDebugMessage(templateAst.Decls)

	decl := templateAst.Decls[1] //skip import declaration
	funcDecl, ok := decl.(*dst.FuncDecl)

	if !ok {
		panic("oops")
	}

	return funcDecl
}

func (g *Generator) printDebugMessage(message any) {
	if g.debug {
		fmt.Println(message)
	}
}
