/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"me.weldnor/swede/core/formatter"
	"me.weldnor/swede/core/lexer"
	"me.weldnor/swede/core/parser"
)

// formatCmd represents the format command.
var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "Format one or more files",
	Run: func(cmd *cobra.Command, args []string) {
		formatFilesParallel(args)
	},
}

func formatFilesParallel(paths []string) {
	statusChan := make(chan string)

	for _, path := range paths {
		path := path
		go func() {
			err := formatFile(path)

			if err != nil {
				statusChan <- "can`t format file" + path
			} else {
				statusChan <- "file " + path + " formatted successfully"
			}
		}()
	}

	for range len(paths) {
		var status string = <-statusChan
		fmt.Println(status)
	}
}

func formatFile(path string) error {
	code, err := os.ReadFile(path)
	if err != nil {
		return errors.New("cant read file: " + path)
	}

	lexer := lexer.NewLexer(string(code))
	parser := parser.NewParser(lexer.Scan())
	parserResult := parser.Parse()

	if len(parserResult.Errors) > 0 {
		return errors.New("found errors while parsing file: " + path)
	}

	_formatter := formatter.NewFormatter(&parserResult.RootNode)
	formattedCode, err := _formatter.FormatParallel()
	if err != nil {
		return err
	}

	os.WriteFile(path, []byte(formattedCode), 0o644)
	return nil
}

func init() {
	rootCmd.AddCommand(formatCmd)
}
