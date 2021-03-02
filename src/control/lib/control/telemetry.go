//
// (C) Copyright 2021 Intel Corporation.
//
// SPDX-License-Identifier: BSD-2-Clause-Patent
//

package control

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	pclient "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

func scrapeMetrics(ctx context.Context, host string, port uint32) (map[string]*pclient.MetricFamily, error) {
	addr := &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", host, port),
		Path:   "metrics",
	}
	body, err := httpGetBody(ctx, addr, http.Get)
	if err != nil {
		return nil, err
	}

	parser := expfmt.TextParser{}
	reader := strings.NewReader(string(body))
	return parser.TextToMetricFamilies(reader)
}

type (
	// MetricsListReq is used to request the list of metrics.
	MetricsListReq struct {
		request
		Port uint32 // Port to use for collecting telemetry data
	}

	// MetricsListResp contains the list of available metrics.
	MetricsListResp struct {
		HostErrorsResp
		AvailableMetrics []*Metric `json:"available_metrics"`
	}

	// Metric describes a specific metric available from the server.
	Metric struct {
		Name        string `json:"metric_name"`
		Description string `json:"metric_description"`
	}
)

// MetricsList fetches the list of available metric types from the DAOS nodes.
func MetricsList(ctx context.Context, req *MetricsListReq) (*MetricsListResp, error) {
	scraped, err := scrapeMetrics(ctx, req.HostList[0], req.Port)
	if err != nil {
		return nil, errors.Wrap(err, "unable to query metrics")
	}

	resp := new(MetricsListResp)

	list := make([]*Metric, 0)
	for name, mf := range scraped {
		newMetric := &Metric{
			Name:        name,
			Description: mf.GetHelp(),
		}
		list = append(list, newMetric)
	}

	resp.AvailableMetrics = list
	return resp, nil
}

type (
	// MetricsQueryReq is used to query telemetry values.
	MetricsQueryReq struct {
		request
		Port        uint32   // port to use for collecting telemetry data
		MetricNames []string // if empty, collects all metrics
	}

	// MetricsQueryResp contains the list of telemetry values per host.
	MetricsQueryResp struct {
		MetricsByHost map[string]MetricValue `json:"metrics_by_host"`
	}

	// MetricValue represents a single item of telemetry data.
	MetricValue struct {
		Metric
		Value string `json:"metric_value"`
	}
)

// MetricsQuery fetches the requested metrics values from the DAOS nodes.
func MetricsQuery(ctx context.Context, req *MetricsQueryReq) (*MetricsQueryResp, error) {
	return nil, nil
}
