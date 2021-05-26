package cmd

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var GlobalViper *viper.Viper

var RootCmd = &cobra.Command{
	PersistentPreRunE: initializeConfig,
	Use:               "cloudforecast-agent",
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
	RootCmd.PersistentFlags().String("loglevel", "info", "minimum log level to print")
	RootCmd.PersistentFlags().String("api-key", "", "Cloudforecast Barometer API key")
	RootCmd.PersistentFlags().String("cluster-uuid", "", "Cloudforecast-provided UUID for the Kubernetes cluster")
	RootCmd.PersistentFlags().String("kubeconfig", "", "Path to kubeconfig file")
	RootCmd.PersistentFlags().String("api-host", "", "Barometer API host")
}

func initializeConfig(cmd *cobra.Command, args []string) error {
	_ = viper.BindPFlag("loglevel", cmd.Flags().Lookup("loglevel"))
	zerolog.SetGlobalLevel(convertLogLevelToZerolog(viper.GetString("loglevel")))

	_ = viper.BindPFlag("apiKey", cmd.Flags().Lookup("api-key"))
	_ = viper.BindEnv("apiKey", "CLOUDFORECAST_BAROMETER_API_KEY")
	_ = viper.BindPFlag("clusterUUID", cmd.Flags().Lookup("cluster-uuid"))
	_ = viper.BindEnv("clusterUUID", "CLOUDFORECAST_BAROMETER_CLUSTER_UUID")
	_ = viper.BindPFlag("kubeconfig", cmd.Flags().Lookup("kubeconfig"))
	_ = viper.BindEnv("apiHost", "CLOUDFORECAST_BAROMETER_API_ENDPOINT")
	_ = viper.BindPFlag("apiHost", cmd.Flags().Lookup("api-host"))

	return nil
}
