package cmd

import "github.com/spf13/cobra"


var configFile string

var serviceCmd = &cobra.Command{
	Use: "server",
	Short:"start server",
	
	PreRun: func(cmd *cobra.Command, args []string) {
		setup()
	},
	
	RunE: func(cmd *cobra.Command, args []string) error {

		return run()
	},

}



func setup(){

}

func run() error{

	return nil
}

func init(){

	serviceCmd.PersistentFlags().StringVarP(&configFile, "config",
		"c", "", "service App file")
}



