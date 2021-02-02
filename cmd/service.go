package cmd

import (
	"BitCoinProfitStrategy/app"
	"BitCoinProfitStrategy/config"
	"context"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var configFile string

var serviceCmd = &cobra.Command{
	Use:   "server",
	Short: "start server",

	PreRun: func(cmd *cobra.Command, args []string) {
		setup()
	},

	RunE: func(cmd *cobra.Command, args []string) error {

		return run()
	},
}

func setup() {
	config.ConfigServer(configFile)
}

func run() error {

	router := app.InitRouter()

	srv := &http.Server{
		Addr: viper.GetString("server.address"),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("server listen failed %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
	log.Println("server shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdonw %v:", err)
	}
	log.Println("server exiting")


	return nil
}

func init() {

	serviceCmd.PersistentFlags().StringVarP(&configFile, "config",
		"c", "", "service App file")
}
