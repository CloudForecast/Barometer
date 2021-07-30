CloudForecast Barometer Agent
==========

Welcome to CloudForecast's Barometer Agent! 

This repository provides a Helm chart to automatically install the Barometer Agent to 
monitor your Kubernetes cluster, as well as optionally installing Prometheus, `kube-state-metrics`,
and `node-exporter` if you do not already have them installed.

## Usage

To install the agent, use the Helm chart.

First, add the repository:

    helm repo add cloudforecast-barometer https://cloudforecast.github.io/Barometer

Then, install the chart:

    helm upgrade --install RELEASE_NAME cloudforecast-barometer/cloudforecast-agent \
        --set barometerAgent.apiKey=<provided api key> \
        --set barometerAgent.clusterUuid=<provided uuid>

To see what configuration options are available, see [`values.yaml`](charts/cloudforecast-agent/values.yaml). As part of the configuration options, you can skip the Prometheus, `kube-state-metrics`, and `node-exporter` install if you already have them installed.

## Contributing

### Updating the version number

When you've made changes to the chart, update the `version` key in 
[`Chart.yaml`](charts/cloudforecast-agent/Chart.yaml) to a new version number. When
ready, push to the `main` branch. The chart will be automatically packaged,
tagged with the correct version number, and published to the Helm repo at
`https://cloudforecast.github.io/barometer`.

If you have made changes to the application itself, update both `appVersion`
and `version` in [`Chart.yaml`](charts/cloudforecast-agent/Chart.yaml). 
A new Docker container with the specified `appVersion` number will be automatically published to
Github Container Registry and included in the chart.
