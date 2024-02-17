package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// lspCmd represents the lsp command
var lspCmd = &cobra.Command{
	Use:   "lsp",
	Short: "Start lsp server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("lsp called")
	},
}

func init() {
	rootCmd.AddCommand(lspCmd)
}
