package zoo

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/achew22/logbook/config"
	"github.com/achew22/logbook/parser"
	"github.com/achew22/logbook/templater"
)

const dateFormat = "2006-01-02"

func trim(s string) string {
	return strings.Trim(s, " \n\t")
}

func TestGeneration(t *testing.T) {
	files, err := ioutil.ReadDir("./testdata")
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			testCase(t, file.Name())
		}
	}
}

func testCase(t *testing.T, dir string) {
	t.Run(dir, func(t *testing.T) {
		logPath := filepath.Join("./testdata", dir)
		c := &config.Config{
			Name:    "Andrew Allen",
			LogPath: logPath,
		}

		// Find all the .out files and generate the matching .md file. If the
		// .md file doesn't match fail.

		files, err := ioutil.ReadDir("./testdata")
		if err != nil {
			t.Fatal(err)
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			name := file.Name()

			if !strings.HasSuffix(name, ".out") {
				return
			}

			toGenFilename := name[0 : len(name)-len(".out")]
			d, err := parser.YmdToDate(toGenFilename)
			if err != nil {
				t.Error(err)
				continue
			}

			generateAndCompare(t, c, d)
		}
	})
}

func generateAndCompare(t *testing.T, c *config.Config, d parser.Date) {
	fileName := fmt.Sprintf("%s.out", d.ToYmd())
	longPath := filepath.Join(c.LogPath, fileName)
	t.Run(fileName, func(t *testing.T) {
		p := parser.New(c)
		want := templater.Print(c, p.Parse(), d)
		got, err := ioutil.ReadFile(longPath)
		if err != nil {
			t.Error(err)
		}

		if diff := cmp.Diff(trim(string(got)), trim(want)); diff != "" {
			t.Errorf("Difference - want + got:\n%s", diff)
		}
	})
}
