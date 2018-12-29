package templater

import (
	"flag"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/achew22/logbook/config"
	"github.com/achew22/logbook/parser"
)

var (
	updateGoldens = flag.Bool("update_goldens", false, "Set this to true to update the golden file that you normally compare against")
)

func ymd(s string) parser.Date {
	d, err := parser.YmdToDate(s)
	if err != nil {
		panic(err)
	}
	return d
}

func trim(s string) string {
	return strings.Trim(s, " \n\t")
}

func TestSimpleTemplating(t *testing.T) {
	c := &config.Config{
		Name: "Andrew Allen",
	}
	today := ymd("2014-02-14")

	got := strings.Split(strings.Trim(Print(c, map[parser.Date]*parser.LogEntry{
		ymd("2014-02-13"): &parser.LogEntry{},
	}, today), " \t\n"), "\n")
	want := strings.Split(strings.Trim(`
# Andrew Allen - 2014-02-14

There are no reminders for today`, " \t\n"), "\n")
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Differences:\n%s\nGot:  %q\nWant: %q", diff, got, want)
	}
}
