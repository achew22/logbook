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
	"github.com/achew22/logbook/config"
	"github.com/achew22/logbook/parser"
	"github.com/achew22/logbook/templater"
	"os"
)

func main() {
	c := &config.Config{
		Name:    "Andrew Allen",
		LogPath: "/usr/local/google/home/achew/logbook",
	}
	p := parser.New(c)
	templater.Print(c, p.Parse())

	os.Exit(0)
}