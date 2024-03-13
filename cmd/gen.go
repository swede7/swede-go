/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"me.weldnor/swede/generator"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate test file",
	Run: func(cmd *cobra.Command, args []string) {
		g := generator.NewGenerator()
		g.Generate()
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
}
