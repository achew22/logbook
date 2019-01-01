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
package parser

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/achew22/logbook/config"
)

var (
	updateGoldens = flag.Bool("update_goldens", false, "Set this to true to update the golden file that you normally compare against")
)

func TestParsing(t *testing.T) {
	files, err := ioutil.ReadDir("./testdata")
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		t.Run(f.Name(), func(t *testing.T) {
			if !f.IsDir() {
				return
			}

			p := New(&config.Config{
				Name:    "Demo person",
				LogPath: filepath.Join("testdata", f.Name()),
			})

			parsedOut := map[string]*LogEntry{}
			for k, v := range p.Parse() {
				parsedOut[k.ToYmd()] = v
			}

			var out bytes.Buffer
			b, err := json.Marshal(parsedOut)
			if err != nil {
				t.Error(err)
			}

			json.Indent(&out, b, "", "  ")

			goldenPath := filepath.Join("testdata", f.Name(), "golden.json")
			goldenData, err := ioutil.ReadFile(goldenPath)
			if err != nil {
				t.Error(err)
			}

			got := strings.Split(trim(out.String()), "\n")
			want := strings.Split(trim(string(goldenData)), "\n")

			if diff := cmp.Diff(got, want); diff != "" {
				t.Errorf("Differences:\n%s", diff)
			}

			if *updateGoldens {
				goldenWriter, err := os.OpenFile(goldenPath, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0644)

				if err != nil {
					t.Error(err)
				}

				_, err = out.WriteTo(goldenWriter)
				if err != nil {
					t.Error(err)
				}
			}
		})
	}
}
