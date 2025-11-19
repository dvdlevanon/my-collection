package main

import (
	"context"
	"my-collection/server/pkg/app"
	"my-collection/server/pkg/utils"
	"os"
	"strings"

	"github.com/op/go-logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var logger = logging.MustGetLogger("main")
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
	config := app.MyCollectionConfig{
		RootDir:                     viper.GetString("root-directory"),
		ListenAddress:               viper.GetString("address"),
		FilesFilter:                 utils.VideoFilter{},
		AutoMixItemsCount:           viper.GetInt("auto-mix-items-count"),
		MixOnDemandItemsCount:       viper.GetInt("mix-on-demand-items-count"),
		ItemsOptimizerMaxResolution: viper.GetInt("items-optimizer-max-resolution"),
		ProcessorPaused:             viper.GetBool("processor-paused"),
		CoversCount:                 viper.GetInt("covers-count"),
		PreviewSceneCount:           viper.GetInt("preview-scene-count"),
		PreviewSceneDuration:        viper.GetInt("preview-scene-duration"),
		OpenSubtitleApiKeys:         viper.GetStringSlice("open-subtitle-api-keys"),
	}

	config.DebugPrint()

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
	if err := utils.ConfigureLogger(); err != nil {
		panic(err)
	}

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: config.yaml in current directory or $HOME/.my-collection.yaml)")

	// Server configuration flags
	rootCmd.Flags().String("root-directory", "", "Server root directory")
	rootCmd.Flags().String("address", ":8080", "Server listen address")

	// Processing configuration flags
	rootCmd.Flags().Int("auto-mix-items-count", 0, "Number of items for auto mix")
	rootCmd.Flags().Int("mix-on-demand-items-count", 30, "Number of items for mix on demand")
	rootCmd.Flags().Int("items-optimizer-max-resolution", 1080, "Maximum resolution for items optimizer")
	rootCmd.Flags().Bool("processor-paused", false, "Whether the processor is paused")

	// Media configuration flags
	rootCmd.Flags().Int("covers-count", 0, "Number of covers to generate")
	rootCmd.Flags().Int("preview-scene-count", 0, "Number of preview scenes to generate")
	rootCmd.Flags().Int("preview-scene-duration", 0, "Duration of preview scenes in seconds")

	// API keys flag (comma-separated or repeated)
	rootCmd.Flags().StringSlice("open-subtitle-api-keys", []string{}, "OpenSubtitles API keys (comma-separated)")
}

func initConfig() {
	// 1. First, set up config file reading
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigType("yaml")

		// Look for config.yaml in current directory first
		if _, err := os.Stat("config.yaml"); err == nil {
			viper.SetConfigFile("config.yaml")
		} else {
			// Fall back to .my-collection.yaml in home directory
			home, err := os.UserHomeDir()
			cobra.CheckErr(err)
			viper.AddConfigPath(home)
			viper.SetConfigName(".my-collection")
		}
	}

	// 2. Read config file (lowest priority)
	if err := viper.ReadInConfig(); err == nil {
		logger.Infof("Using config file: %s", viper.ConfigFileUsed())
	}

	// 3. Set up environment variables (middle priority)
	viper.SetEnvPrefix("my_collection")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.AutomaticEnv()
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "MY_COLLECTION_") {
			logger.Debugf(env)
		}
	}

	// 4. Bind CLI flags LAST (highest priority)
	// This must happen after ReadInConfig() so flags override config file values
	viper.BindPFlags(rootCmd.Flags())
	logger.Debugf("Command Line Args: %v", os.Args)
}
