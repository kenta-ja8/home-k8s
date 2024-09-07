package main

import (
	"os"

	"github.com/kenta-ja8/home-k8s-app/pkg/client"
	"github.com/kenta-ja8/home-k8s-app/pkg/entity"
	"github.com/kenta-ja8/home-k8s-app/pkg/logger"
	"github.com/kenta-ja8/home-k8s-app/pkg/usecase"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{}

var yourCommand = &cobra.Command{
	Use:   "sample",
	Short: "sample command",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Info("Start sample command")
		cfg := entity.LoadConfig()
		db, err := client.NewPostgresClient(cfg)
		if err != nil {
			return err
		}
		sampleUsecase := usecase.NewSampleUsecase(cfg, db)
		return sampleUsecase.InsertRecord()
	},
}

func init() {
	rootCmd.AddCommand(yourCommand)
	yourCommand.Flags().StringP("flagname", "f", "defaultValue", "Flag description")
}

func main() {
	logger.Info("Hello World!")
	defer logger.Info("Goodbye World!")

	if err := rootCmd.Execute(); err != nil {
		logger.Error("Error:", err)
		os.Exit(1)
	}
}
