package main

import (
	"context"
	"fmt"
	"my-collection/server/pkg/app"
	"my-collection/server/pkg/utils"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "my-collection",
	Short: "TBD",
	Long:  `TBD`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := run(); err != nil {
			utils.LogError("Error in main", err)
			os.Exit(1)
		}
	},
}

func run() error {
	if err := utils.ConfigureLogger(); err != nil {
		return err
	}

	config := app.MyCollectionConfig{
		RootDir:                     viper.GetString("root-directory"),
		ListenAddress:               viper.GetString("address"),
		FilesFilter:                 utils.VideoFilter{},
		AutoMixItemsCount:           40,
		MixOnDemandItemsCount:       30,
		ItemsOptimizerMaxResolution: 1080,
		ProcessorPaused:             false,
		CoversCount:                 3,
		PreviewSceneCount:           4,
		PreviewSceneDuration:        3,
	}

	mc, err := app.New(config)
	if err != nil {
		return err
	}

	return mc.Run(context.Background())
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func main() {
	Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.my-collection.yaml)")
	rootCmd.Flags().String("root-directory", "", "Server root directory")
	rootCmd.Flags().String("address", ":6969", "Server listen address")

	viper.BindPFlags(rootCmd.Flags())
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".my-collection")
	}

	viper.SetEnvPrefix("my_collection")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
