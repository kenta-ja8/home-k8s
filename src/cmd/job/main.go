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
	rootCmd.AddCommand(&cobra.Command{
		Use: "pantry-order-reminder",
		RunE: func(cmd *cobra.Command, args []string) error {
			uc := usecase.NewPantryOrderReminderUsecase()
			return uc.Run(cmd.Context())
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use: "collect-natureremo",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			cfg := entity.LoadConfig()
			db, err := client.NewPostgresClient(cfg)
			if err != nil {
				return err
			}
			natureremoClient := client.NewNatureremoClient(cfg)
			uc := usecase.NewNatureremo(natureremoClient, db)
			return uc.Collect(ctx)
		},
	})
}

func main() {
	logger.Info("start job")
	defer logger.Info("end job")

	cfg := entity.LoadConfig()
	logger.Init(cfg)

	if err := rootCmd.Execute(); err != nil {
		logger.Error("Error: %v", err)
		os.Exit(1)
	}
}
