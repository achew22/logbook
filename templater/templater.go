/**
 * Copyright 2018 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the license.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specified language governing permissions and
 * limitations under the License.
 */
package templater

import (
	"fmt"
	"os"
	"strings"

	"github.com/achew22/logbook/config"
	"github.com/achew22/logbook/parser"
)

type extractedEntry struct {
	originDate parser.Date
	text       string
}

func Print(c *config.Config, entries map[parser.Date]*parser.LogEntry, today parser.Date) string {
	buf := &strings.Builder{}

	fmt.Fprintf(buf, "# %s - %s\n\n", c.Name, today.ToYmd())

	todayLog, ok := entries[today]
	if !ok {
		fmt.Fprintf(buf, "There are no reminders for today\n\n")
	} else {
		if len(todayLog.PastReferences) > 0 {
			fmt.Fprintf(buf, "## Reminders:\n\n")

			for originDate, messages := range todayLog.PastReferences {
				fmt.Fprintf(buf, "From %s:\n\n", originDate.ToYmd())
				for _, m := range messages {
					fmt.Fprintf(buf, " *  %s\n", m)
				}
			}
		}
		fmt.Fprintf(buf, "\n")
	}

	parseErrors := map[parser.Date][]*parser.ParseError{}
	for d, entry := range entries {
		if len(entry.Errors) > 0 {
			parseErrors[d] = entry.Errors
		}
	}

	if len(parseErrors) > 0 {
		fmt.Fprintf(buf, "## Parse errors\n\n")
		for d, errors := range parseErrors {
			fmt.Fprintf(buf, "From %s:\n\n", d.ToYmd())
			for _, err := range errors {
				fmt.Fprintf(buf, " *  %s\n", err.Message)
			}
			fmt.Fprintf(buf, "\n")
		}
	}

	fmt.Fprintf(buf, "\n")

	fmt.Fprintf(os.Stderr, "Wrote file\n")
	return buf.String()
}
