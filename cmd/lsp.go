package cmd

import (
	"github.com/spf13/cobra"
	"github.com/swede7/swede-go/lsp"
)

var stdioFlag bool

// lspCmd represents the lsp command.
var lspCmd = &cobra.Command{
	Use:   "lsp",
	Short: "Start lsp server",
	Run: func(cmd *cobra.Command, args []string) {
		server := lsp.NewLspServer()
		server.Start()
	},
}

func init() {
	rootCmd.AddCommand(lspCmd)
	lspCmd.Flags().BoolVarP(&stdioFlag, "stdio", "s", false, "use stdio as communication protocol")
}
