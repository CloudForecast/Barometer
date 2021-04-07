package cmd

import (
	"errors"
	"github.com/CloudForecast/Barometer/pkg"
	"github.com/CloudForecast/Barometer/pkg/barometerApi"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

var dryRun bool

var agentCmd = &cobra.Command{
	Use: "agent",
	Short: "run the Cloudforecast barometer agent",
	RunE: run,
}

func run(cmd *cobra.Command, args []string) error {
	apiKey := viper.GetString("apiKey")
	if apiKey == "" {
		return errors.New("a Cloudforecast api key is required")
	}

	client := barometerApi.NewBarometerApi(apiKey)

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
	_ = viper.BindPFlag("dryrun", agentCmd.Flags().Lookup("dryrun"))
	RootCmd.AddCommand(agentCmd)
}
