package cmd

import (
	"errors"
	"github.com/CloudForecast/barometer/pkg"
	"github.com/CloudForecast/barometer/pkg/barometerApi"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

var dryRun bool

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "run the CloudForecast barometer agent",
	RunE:  run,
}

func run(cmd *cobra.Command, args []string) error {
	apiKey := viper.GetString("apiKey")
	clusterUUID := viper.GetString("clusterUUID")
	if apiKey == "" {
		return errors.New("a Cloudforecast api key is required")
	}

	client := barometerApi.NewBarometerApi(apiKey, clusterUUID)

	gracefulExit, err := pkg.RunAll(client)
	if err != nil {
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-c:
		log.Info().Msgf("Received signal %s to exit", sig.String())
		gracefulExit()
	}

	return nil
}

func init() {
	agentCmd.Flags().BoolVar(&dryRun, "dryrun", false, "Dry run, do not actually send POST requests")
	agentCmd.Flags().String("prometheus-url", "", "Prometheus service address")
	agentCmd.Flags().String("schedule", "auto-generated", "Cron schedule for fetching and sending metrics")
	_ = viper.BindPFlag("dryrun", agentCmd.Flags().Lookup("dryrun"))
	_ = viper.BindEnv("prometheusUrl", "CLOUDFORECAST_PROMETHEUS_HTTP_API_URL")
	_ = viper.BindPFlag("prometheusUrl", agentCmd.Flags().Lookup("prometheus-url"))
	_ = viper.BindEnv("schedule", "CLOUDFORECAST_BAROMETER_CRON_SCHEDULE")
	_ = viper.BindPFlag("schedule", agentCmd.Flags().Lookup("schedule"))
	RootCmd.AddCommand(agentCmd)
}
