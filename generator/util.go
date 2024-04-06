package generator

import (
	"errors"
	"go/ast"
	"os"
	"path"
)

func getProcessedFilePath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(errors.New("can't get current working directory"))
	}

	return path.Join(wd, os.Getenv("GOFILE"))
}

func funcDeclVisitor(decls []ast.Decl, f func(*ast.FuncDecl)) {
	for _, decl := range decls {
		funcDecl, ok := decl.(*ast.FuncDecl)

		if !ok {
			continue
		}

		f(funcDecl)
	}
}

func funcDeclHasComment(funcDecl *ast.FuncDecl, checkFunc func(string) bool) bool {
	if funcDecl.Doc == nil {
		return false
	}

	for _, comment := range funcDecl.Doc.List {
		if checkFunc(comment.Text) {
			return true
		}
	}

	return false
}

func funcDeclGetComment(funcDecl *ast.FuncDecl, checkFunc func(string) bool) (string, error) {
	if funcDecl.Doc == nil {
		return "", errors.New("no comments")
	}

	for _, comment := range funcDecl.Doc.List {
		if checkFunc(comment.Text) {
			return comment.Text, nil
		}
	}

	return "", errors.New("no comment found by condition")
}
