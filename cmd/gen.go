/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/swede7/swede-go/generator"
)

var featureFiles []string

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate test file",
	Run: func(cmd *cobra.Command, args []string) {
		generatorOptions := generator.Options{}

		generatorOptions.Source = getSource()

		if featureFiles != nil {
			generatorOptions.FeatureFiles = featureFiles
		}

		g := generator.NewGenerator(generatorOptions)
		newSource := g.Generate()

		writeNewSource(newSource)
	},
}

func getSource() string {
	path := getProcessedFilePath()

	bs, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	return string(bs)
}

func writeNewSource(source string) {
	path := getProcessedFilePath()

	err := os.WriteFile(path, []byte(source), 0644)
	if err != nil {
		return
	}
}

func getProcessedFilePath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(errors.New("can't get current working directory"))
	}

	return path.Join(wd, os.Getenv("GOFILE"))
}

func init() {
	genCmd.Flags().StringSliceVarP(&featureFiles, "feature-file", "f", []string{}, "Feature files")
	rootCmd.AddCommand(genCmd)
}
