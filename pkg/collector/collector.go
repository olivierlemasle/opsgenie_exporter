package collector

import (
	"context"
	"errors"

	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/schedule"
	"github.com/opsgenie/opsgenie-go-sdk-v2/team"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"github.com/olivierlemasle/opsgenie_exporter/pkg/config"
)

const pageSize = 100

type opsgenieCollector struct {
	alert    *alert.Client
	schedule *schedule.Client
	team     *team.Client
	logger   *logrus.Logger
}

func NewCollector(logger *logrus.Logger) (prometheus.Collector, error) {
	apiKey := config.Get("apiKey")
	if apiKey == "" {
		return nil, errors.New("No Opsgenie API key provided")
	}

	var apiURL client.ApiUrl
	if url := config.Get("apiUrl"); url != "" {
		apiURL = client.ApiUrl(url)
	} else {
		apiURL = client.API_URL
	}

	config := &client.Config{
		ApiKey:         apiKey,
		OpsGenieAPIURL: apiURL,
		Logger:         logger,
	}

	// TODO: configure proxy
	// TODO: configure timeouts

	alertClient, err := alert.NewClient(config)
	if err != nil {
		return nil, err
	}

	scheduleClient, err := schedule.NewClient(config)
	if err != nil {
		return nil, err
	}

	teamClient, err := team.NewClient(config)
	if err != nil {
		return nil, err
	}

	return opsgenieCollector{
		alert:    alertClient,
		schedule: scheduleClient,
		team:     teamClient,
		logger:   logger,
	}, nil
}

func (c opsgenieCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func (c opsgenieCollector) Collect(ch chan<- prometheus.Metric) {
	ctx := context.Background()

	c.logger.Info("Collect metrics...")

	doneCollectOnCalls := make(chan bool)
	doneCollectAlerts := make(chan bool)

	go c.collectOnCalls(ctx, ch, doneCollectOnCalls)
	go c.collectAlerts(ctx, ch, doneCollectAlerts)

	<-doneCollectOnCalls
	<-doneCollectAlerts

	c.logger.Info("Metrics collected")
}
