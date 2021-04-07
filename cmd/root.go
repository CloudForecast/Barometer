package cmd

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var LogLevel string
var ApiKey string
var GlobalViper *viper.Viper

var RootCmd = &cobra.Command{
	PersistentPreRunE: initializeConfig,
	Use: "cloudforecast-agent",
}

func convertLogLevelToZerolog(input string) zerolog.Level {
	switch input {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		panic(fmt.Sprintf("invalid loglevel %s provided", input))
	}
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&LogLevel, "loglevel", "info", "minimum log level to print")
	RootCmd.PersistentFlags().StringVar(&ApiKey, "apikey", "", "Cloudforecast Barometer API key")
}

func initializeConfig(cmd *cobra.Command, args []string) error {
	viper.SetEnvPrefix("cloudforecast")

	_ = viper.BindPFlag("loglevel", cmd.Flags().Lookup("loglevel"))
	zerolog.SetGlobalLevel(convertLogLevelToZerolog(viper.GetString("loglevel")))

	_ = viper.BindEnv("prometheusUrl", "PROMETHEUS_HTTP_API_URL")
	_ = viper.BindEnv("apiKey", "BAROMETER_API_KEY")
	_ = viper.BindPFlag("apiKey", cmd.Flags().Lookup("apikey"))

	return nil
}