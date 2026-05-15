package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// PrintDiff writes a human-readable diff summary to stdout.
func PrintDiff(result DiffResult) {
	printDiffTo(os.Stdout, result)
}

func printDiffTo(w io.Writer, result DiffResult) {
	if len(result.OnlyLeft) == 0 && len(result.OnlyRight) == 0 {
		fmt.Fprintln(w, "# no differences found")
		return
	}
	for _, l := range result.OnlyLeft {
		fmt.Fprintf(w, "< %s\n", l.Raw)
	}
	for _, r := range result.OnlyRight {
		fmt.Fprintf(w, "> %s\n", r.Raw)
	}
	fmt.Fprintf(w, "# common: %d  only-left: %d  only-right: %d\n",
		len(result.Common), len(result.OnlyLeft), len(result.OnlyRight))
}

type diffJSON struct {
	OnlyLeft  []string `json:"only_left"`
	OnlyRight []string `json:"only_right"`
	Common    int      `json:"common"`
}

// DiffJSON serialises a DiffResult to a JSON string.
func DiffJSON(result DiffResult) string {
	d := diffJSON{
		Common: len(result.Common),
	}
	for _, l := range result.OnlyLeft {
		d.OnlyLeft = append(d.OnlyLeft, l.Raw)
	}
	for _, r := range result.OnlyRight {
		d.OnlyRight = append(d.OnlyRight, r.Raw)
	}
	if d.OnlyLeft == nil {
		d.OnlyLeft = []string{}
	}
	if d.OnlyRight == nil {
		d.OnlyRight = []string{}
	}
	b, err := json.Marshal(d)
	if err != nil {
		return "{}"
	}
	return string(b)
}
