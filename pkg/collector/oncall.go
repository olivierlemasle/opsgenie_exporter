package collector

import (
	"context"
	"sync"

	"github.com/opsgenie/opsgenie-go-sdk-v2/schedule"
	"github.com/prometheus/client_golang/prometheus"
)

func (c opsgenieCollector) collectOnCalls(ctx context.Context, ch chan<- prometheus.Metric, done chan<- bool) {
	schedules, _ := c.listSchedules(ctx)

	wg := sync.WaitGroup{}
	wg.Add(len(schedules))

	for _, schedule := range schedules {
		go func(schedule string, ch chan<- prometheus.Metric) {
			recipients, _ := c.getOnCallRecipients(ctx, schedule)
			for _, recipient := range recipients {
				ch <- prometheus.MustNewConstMetric(
					prometheus.NewDesc("opsgenie_oncall_recipient", "Opsgenie oncall recipient", []string{"schedule_name", "recipient"}, nil),
					prometheus.GaugeValue,
					float64(1),
					schedule, recipient,
				)
			}
			wg.Done()
		}(schedule, ch)
	}

	wg.Wait()
	done <- true
}

func (c opsgenieCollector) listSchedules(ctx context.Context) ([]string, error) {
	res := []string{}
	c.logger.Info("Listing schedules")
	expand := false
	result, err := c.schedule.List(ctx, &schedule.ListRequest{Expand: &expand})
	if err != nil {
		return res, err
	}
	for _, s := range result.Schedule {
		res = append(res, s.Name)
	}
	return res, nil
}

func (c opsgenieCollector) getOnCallRecipients(ctx context.Context, scheduleName string) ([]string, error) {
	flat := true
	c.logger.Infof("Getting oncall recipients for schedule %v", scheduleName)
	result, err := c.schedule.GetOnCalls(ctx, &schedule.GetOnCallsRequest{
		Flat:                   &flat,
		ScheduleIdentifierType: schedule.Name,
		ScheduleIdentifier:     scheduleName,
	})
	if err != nil {
		return []string{}, err
	}
	return result.OnCallRecipients, nil
}
