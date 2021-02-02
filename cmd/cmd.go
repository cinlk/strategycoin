package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{}

func Run() error {

	return rootCmd.Execute()

}

func init() {

	rootCmd.AddCommand(serviceCmd)
}
