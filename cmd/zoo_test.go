package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/achew22/logbook/config"
	"github.com/achew22/logbook/parser"
)

var (
	updateGoldens = flag.Bool("update_goldens", false, "Set to true to update the goldens on disk")
)

const dateFormat = "2006-01-02"

func TestZoo(t *testing.T) {
	files, err := ioutil.ReadDir("./testdata")
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			testCase(t, filepath.Join("./testdata", file.Name()))
		}
	}
}

func testCase(t *testing.T, dir string) {
	t.Run(dir, func(t *testing.T) {
		homeDir, cleanup := makeFakeHome(t)
		_ = cleanup
		//defer cleanup()

		makeLogbookDirectoryInHome(t, homeDir)

		logPath := filepath.Join(homeDir, "logbook")
		t.Logf("Logbook path: %s", logPath)
		c := &config.Config{
			Name:    "Andrew Allen",
			LogPath: logPath,
		}

		// Find all the .out files and generate the matching .md file. If the
		// .md file doesn't match fail.

		files, err := ioutil.ReadDir(dir)
		if err != nil {
			t.Fatalf("Unable to ReadDir(%s): %v", dir, err)
		}

		// Iterate over the directory and copy all .md and .out files to the
		// temporary home directory. That way modification isn't destructive.
		var toGenerate []parser.Date
		for _, file := range files {
			if file.IsDir() {
				continue
			}

			name := file.Name()
			switch filepath.Ext(name) {
			case ".out":
				toGenFilename := name[0 : len(name)-len(".out")]
				d, err := parser.YmdToDate(toGenFilename)
				if err != nil {
					t.Error(err)
					continue
				}
				toGenerate = append(toGenerate, d)
				fallthrough
			case ".md":
				// Copy the file over into the temp logbook.
				s := filepath.Join(dir, name)
				d := filepath.Join(logPath, name)
				cp(t, s, d)
			default:
				t.Logf("%s was not copied", name)
			}
		}

		for _, d := range toGenerate {
			generateAndCompare(t, c.LogPath, d)

			if *updateGoldens {
				src := filepath.Join(homeDir, "logbook", d.ToYmd()+".md")
				dst := filepath.Join(dir, d.ToYmd()+".out")
				t.Logf("CPing: %s to %s", src, dst)
				cp(t, src, dst)
			}
		}
	})
}

func generateAndCompare(t *testing.T, logPath string, d parser.Date) {
	t.Run(d.ToYmd(), func(t *testing.T) {
		// Invoke the main function on the date
		homeDir := filepath.Dir(logPath)
		out, err := helperCommand(
			t, homeDir,
			"--date_override="+d.ToYmd(),
			"--name_override=Andrew Allen",
		).CombinedOutput()
		if err != nil {
			t.Errorf("Unable to run generator: %v\nOutput:%s", err, out)
		}

		wantFileName := fmt.Sprintf("%s.out", d.ToYmd())
		wantLongPath := filepath.Join(logPath, wantFileName)
		want, err := ioutil.ReadFile(wantLongPath)
		if err != nil {
			t.Errorf("Unable to read file: %v", err)
		}

		assertLogEntry(t, homeDir, d.ToYmd(), string(want))
	})
}
