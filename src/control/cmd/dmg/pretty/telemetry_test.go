//
// (C) Copyright 2021 Intel Corporation.
//
// SPDX-License-Identifier: BSD-2-Clause-Patent
//

package pretty

import (
	"errors"
	"strings"
	"testing"

	"github.com/daos-stack/daos/src/control/common"
	"github.com/daos-stack/daos/src/control/lib/control"
)

func TestPretty_PrintMetricsListResp(t *testing.T) {
	for name, tc := range map[string]struct {
		resp      *control.MetricsListResp
		writeErr  error
		expOutput string
		expErr    error
	}{
		"nil resp": {
			expErr: errors.New("nil response"),
		},
		"nil list": {
			resp: &control.MetricsListResp{},
		},
		"empty list": {
			resp: &control.MetricsListResp{
				AvailableMetrics: []*control.Metric{},
			},
		},
		"one item": {
			resp: &control.MetricsListResp{
				AvailableMetrics: []*control.Metric{
					{
						Name:        "test_metric_1",
						Description: "Test Metric",
					},
				},
			},
			expOutput: `
Metric Name   Description 
-----------   ----------- 
test_metric_1 Test Metric 
`,
		},
		"multi item": {
			resp: &control.MetricsListResp{
				AvailableMetrics: []*control.Metric{
					{
						Name:        "test_metric_1",
						Description: "Test metric",
					},
					{
						Name:        "test_metric_2",
						Description: "Another test metric",
					},
					{
						Name:        "funny_hats",
						Description: "Hilarious headwear",
					},
				},
			},
			expOutput: `
Metric Name   Description         
-----------   -----------         
test_metric_1 Test metric         
test_metric_2 Another test metric 
funny_hats    Hilarious headwear  
`,
		},
		"write failure": {
			resp: &control.MetricsListResp{
				AvailableMetrics: []*control.Metric{
					{
						Name:        "test_metric_1",
						Description: "Test Metric",
					},
				},
			},
			writeErr: errors.New("mock write"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			writer := &common.MockWriter{
				WriteErr: tc.writeErr,
			}

			err := PrintMetricsListResp(tc.resp, writer)

			common.CmpErr(t, tc.expErr, err)
			common.AssertEqual(t, strings.TrimLeft(tc.expOutput, "\n"), writer.GetWritten(), "")
		})
	}
}
