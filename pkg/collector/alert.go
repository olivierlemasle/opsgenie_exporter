package collector

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
	"github.com/opsgenie/opsgenie-go-sdk-v2/team"
	"github.com/prometheus/client_golang/prometheus"
)

func (c opsgenieCollector) collectAlerts(ctx context.Context, ch chan<- prometheus.Metric, done chan<- bool) {
	teams, _ := c.listTeams(ctx)

	wg := sync.WaitGroup{}
	wg.Add(len(teams))

	for _, team := range teams {
		go func(team string, ch chan<- prometheus.Metric) {
			totalAlertsCount, _ := c.countAlerts(ctx, team)
			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc("opsgenie_alerts_total", "Number of Opsgenie alerts", []string{"team"}, nil),
				prometheus.CounterValue,
				float64(totalAlertsCount),
				team,
			)

			openAlertsCountByGroup, _ := c.listOpenAlerts(ctx, team)
			for g, count := range openAlertsCountByGroup {
				ch <- prometheus.MustNewConstMetric(
					prometheus.NewDesc("opsgenie_alerts_open", "Number of open Opsgenie alerts", []string{"team", "priority", "acknowledged"}, nil),
					prometheus.GaugeValue,
					float64(count),
					team, g.priority, strconv.FormatBool(g.acknowledged),
				)
			}
			wg.Done()
		}(team, ch)
	}

	wg.Wait()
	done <- true
}

func (c opsgenieCollector) listTeams(ctx context.Context) ([]string, error) {
	res := []string{}
	c.logger.Info("Listing teams")
	result, err := c.team.List(ctx, &team.ListTeamRequest{})
	if err != nil {
		return res, err
	}
	for _, team := range result.Teams {
		res = append(res, team.Name)
	}
	return res, nil
}

type alertGroup struct {
	priority     string
	acknowledged bool
}

func (c opsgenieCollector) listOpenAlerts(ctx context.Context, team string) (map[alertGroup]int, error) {
	offset := 0
	res := make(map[alertGroup]int)
	for {
		c.logger.Infof("Listing active alerts for team %v (offset %v)", team, offset)
		result, err := c.alert.List(ctx, &alert.ListAlertRequest{
			Limit:  pageSize,
			Offset: offset,
			Query:  fmt.Sprintf("teams:%s status:open", team),
		})
		if err != nil {
			return res, err
		}
		if len(result.Alerts) == 0 {
			return res, nil
		}
		for _, alert := range result.Alerts {
			priority := string(alert.Priority)
			g := alertGroup{priority, alert.Acknowledged}
			count, ok := res[g]
			if ok {
				res[g] = count + 1
			} else {
				res[g] = 1
			}
		}
		offset = offset + pageSize
	}
}

func (c opsgenieCollector) countAlerts(ctx context.Context, team string) (int, error) {
	c.logger.Infof("Counting all alerts for team %v", team)
	result, err := c.alert.CountAlerts(ctx, &alert.CountAlertsRequest{
		Query: fmt.Sprintf("teams:%v", team),
	})
	if err != nil {
		return 0, nil
	}
	return result.Count, nil
}
