package generator

import (
	"errors"
	"fmt"
	"go/ast"
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

func (g *Generator) getStepDefinitionFuncs() {
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

	// fmt.Print("scanning comments...")
	// commentGroups := g.astFile.Comments

	// for _, commentGroup := range commentGroups {
	// 	for _, comment := range commentGroup.List {
	// 		if strings.Contains(comment.Text, "swede:step") {
	// 			endPosition := g.fset.PositionFor(comment.End(), false) // todo second parameter?
	// 			fmt.Println("found swede comment on pos " + endPosition.String())
	// 		}
	// 	}
	// }
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

func (g *Generator) Generate() {

	g.getStepDefinitionFuncs()
}
