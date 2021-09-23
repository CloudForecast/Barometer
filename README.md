CloudForecast Barometer Agent
==========

Welcome to CloudForecast's Barometer Agent! 

This repository provides a Helm chart to automatically install the Barometer Agent to 
monitor your Kubernetes cluster, as well as optionally installing Prometheus, `kube-state-metrics`,
and `node-exporter` if you do not already have them installed.

## Usage

To install the agent, you can use our Helm chart.

First, add the repository:

    helm repo add cloudforecast-barometer https://cloudforecast.github.io/barometer

If you already have a monitoring stack running (Prometheus, `kube-state-metrics`, and `node-exporter`), you can just point our agent to your prometheus endpoint (e.g: _http://prometheus-server.default.svc.cluster.local_):  

    helm upgrade --install cloudforecast-barometer cloudforecast-barometer/cloudforecast-agent \
        --set barometerAgent.clusterUuid=<provided uuid> \
        --set barometerAgent.apiKey=<provided api key> \
        --set barometerAgent.prometheusUrlOverride=<prometheus endpoint>

If you , install the chart:

    helm upgrade --install cloudforecast-barometer cloudforecast-barometer/cloudforecast-agent \
        --set barometerAgent.clusterUuid=<provided uuid> \
        --set barometerAgent.apiKey=<provided api key> \
        --set prometheus.enabled=true


To see what configuration options are available, see [`values.yaml`](charts/cloudforecast-agent/values.yaml).

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
