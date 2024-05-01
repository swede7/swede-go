package generator

import (
	"errors"
	"github.com/dave/dst"
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

func funcDeclVisitor(decls []dst.Decl, f func(*dst.FuncDecl)) {
	for _, decl := range decls {
		funcDecl, ok := decl.(*dst.FuncDecl)

		if !ok {
			continue
		}

		f(funcDecl)
	}
}

func funcDeclHasComment(funcDecl *dst.FuncDecl, checkFunc func(string) bool) bool {
	if funcDecl.Decs.Start == nil {
		return false
	}

	for _, _comment := range funcDecl.Decs.Start.All() {
		if checkFunc(_comment) {
			return true
		}
	}

	return false
}

func funcDeclGetComment(funcDecl *dst.FuncDecl, checkFunc func(string) bool) (string, error) {
	if funcDecl.Decs.Start == nil {
		return "", errors.New("no comments")
	}

	for _, _comment := range funcDecl.Decs.Start.All() {
		if checkFunc(_comment) {
			return _comment, nil
		}
	}

	return "", errors.New("no _comment found by condition")
}
