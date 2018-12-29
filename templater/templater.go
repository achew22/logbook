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
		fmt.Fprintf(buf, "There are no reminders for today\n")
		// There are no entries waiting for this day's render.
		return buf.String()
	}
	if len(todayLog.PastReferences) > 0 {
		fmt.Fprintf(buf, "Reminders:\n")

		for originDate, messages := range todayLog.PastReferences {
			for _, m := range messages {
				fmt.Fprintf(buf, "From %s: %s\n", originDate.ToYmd(), m)
			}
		}
	}

	fmt.Fprintf(buf, "\n")

	fmt.Fprintf(os.Stderr, "Wrote file\n")
	return buf.String()
}
