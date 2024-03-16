/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"me.weldnor/swede/generator"
)

var featureFiles []string

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate test file",
	Run: func(cmd *cobra.Command, args []string) {

		generatorOptions := generator.GeneratorOptions{}

		if featureFiles != nil {
			generatorOptions.FeatureFiles = featureFiles
		}

		fmt.Println(featureFiles)

		g := generator.NewGenerator(generatorOptions)
		g.Generate()
	},
}

func init() {
	genCmd.Flags().StringSliceVarP(&featureFiles, "feature-file", "f", []string{}, "Feature files")
	rootCmd.AddCommand(genCmd)
}
