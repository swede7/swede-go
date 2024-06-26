/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/swede7/swede-go/core/lang/swede/formatter"
	"github.com/swede7/swede-go/core/lang/swede/parser"
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
	parserResult := parser.ParseFile(path)

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
