//
// (C) Copyright 2021 Intel Corporation.
//
// SPDX-License-Identifier: BSD-2-Clause-Patent
//

package pretty

import (
	"errors"
	"io"

	"github.com/daos-stack/daos/src/control/lib/control"
	"github.com/daos-stack/daos/src/control/lib/txtfmt"
)

// PrintMetricsListResp formats a list of metrics for display and writes them to
// the output.
func PrintMetricsListResp(resp *control.MetricsListResp, out io.Writer) error {
	if resp == nil {
		return errors.New("nil response")
	}

	if len(resp.AvailableMetrics) == 0 {
		return nil
	}

	nameTitle := "Metric Name"
	descTitle := "Description"

	tablePrint := txtfmt.NewTableFormatter(nameTitle, descTitle)
	tablePrint.InitWriter(out)
	table := []txtfmt.TableRow{}

	for _, m := range resp.AvailableMetrics {
		row := txtfmt.TableRow{
			nameTitle: m.Name,
			descTitle: m.Description,
		}

		table = append(table, row)
	}

	tablePrint.Format(table)

	return nil
}
