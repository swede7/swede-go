/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"sync"

	"github.com/spf13/cobra"
	"github.com/swede7/swede-go/core/lang/swede/linter"
	"github.com/swede7/swede-go/core/lang/swede/parser"
)

// lintCmd represents the lint command.
var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Lint one or more files",

	Run: func(cmd *cobra.Command, args []string) {
		lintFilesParallel(args)
	},
}

func lintFilesParallel(paths []string) {
	wg := &sync.WaitGroup{}
	mutex := &sync.Mutex{}

	for _, path := range paths {
		wg.Add(1)

		path := path
		go func() {
			defer wg.Done()
			linterErrors, err := lintFile(path, mutex)

			mutex.Lock()
			defer mutex.Unlock()

			fmt.Println("processing " + path + "...")

			if err != nil {
				fmt.Println(err.Error())
				return
			}

			if len(linterErrors) == 0 {
				fmt.Println("no errors found. GOOD")
				return
			}

			for _, linterError := range linterErrors {
				fmt.Printf("%s [%d:%d] %s\n", linterError.Severity, linterError.StartPosition.Line+1, linterError.StartPosition.Column+1, linterError.Message)
			}
		}()
	}

	wg.Wait()
}

func lintFile(path string, mutex *sync.Mutex) ([]linter.LinterError, error) {
	parserResult := parser.ParseFile(path)

	if len(parserResult.Errors) > 0 {
		return nil, errors.New("found errors while parsing file: " + path)
	}

	linter := linter.NewLinter(&parserResult.RootNode)
	return linter.Lint(), nil
}

func init() {
	rootCmd.AddCommand(lintCmd)
}
