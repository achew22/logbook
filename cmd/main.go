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
package main

import (
	"fmt"
	"os"

	"github.com/achew22/logbook/config"
	"github.com/achew22/logbook/parser"
	"github.com/achew22/logbook/templater"
)

const dateFormat = "2006-01-02"

func main() {
	c := &config.Config{
		Name:    "Andrew Allen",
		LogPath: "/usr/local/google/home/achew/logbook",
	}
	p := parser.New(c)

	todayPath := filepath.Join(c.LogPath, fmt.Sprintf("%s.md", time.Now().Format(dateFormat)))

	if _, err := os.Stat(todayPath); err == nil {
		fmt.Fprintf(os.Stderr, "A file already exists by the name %s\n", todayPath)
		os.Exit(1)
	}

	out, err := os.Create(todayPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create a log entry named %s.\nErr: %v", todayPath, err)
		os.Exit(1)
	}

	today := parser.TimeToDate(time.Now())
	text := templater.Print(c, p.Parse(), today)

	out.Write([]byte(text))
	if err := out.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}
