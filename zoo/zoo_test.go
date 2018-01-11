package zoo

import (
	"io/ioutil"
	"path/filepath"
	//"os"
	"fmt"
	"testing"
	"time"

	"github.com/achew22/logbook/config"
	"github.com/achew22/logbook/parser"
	"github.com/achew22/logbook/templater"
)

const dateFormat = "2006-01-02"

func TestData(t *testing.T) {
	c := &config.Config{
		Name:    "Andrew Allen",
		LogPath: "./testdata",
	}
	p := parser.New(c)

	todayTime := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	today := parser.TimeToDate(todayTime)
	todayPath := filepath.Join(c.LogPath, fmt.Sprintf("%s.out", todayTime.Format(dateFormat)))
	t.Run("", func(t *testing.T) {
		want := templater.Print(c, p.Parse(), today)
		got, err := ioutil.ReadFile(todayPath)
		if err != nil {
			t.Error(err)
		}

		if want != string(got) {
			t.Errorf("Want: %s\nGot:  %s", want, got)
		}
	})
}
