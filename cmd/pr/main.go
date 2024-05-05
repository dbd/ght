package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/dbd/ght/components"
	"github.com/dbd/ght/utils"
)

var diff string = `diff --git a/go.mod b/go.mod
index c26c597..66acc82 100644
--- a/go.mod
+++ b/go.mod
@@ -1,6 +1,6 @@
 module github.com/dbd/strategy

-go 1.18
+go 1.20

 require (
        github.com/alecthomas/kong v0.2.17
@@ -32,7 +32,7 @@ require (
        github.com/pmezard/go-difflib v1.0.0 // indirect
        github.com/rivo/uniseg v0.2.0 // indirect
        github.com/spf13/pflag v1.0.5 // indirect
-       golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c // indirect
+       golang.org/x/sys v0.1.0 // indirect
        gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
        gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
 )
diff --git a/go.sum b/go.sum
index 30a364f..a8ec724 100644
--- a/go.sum
+++ b/go.sum
@@ -64,8 +64,9 @@ golang.org/x/sys v0.0.0-20200116001909-b77594299b42/go.mod h1:h1NjWce9XRLGQEsW7w
 golang.org/x/sys v0.0.0-20200916030750-2334cc1a136f/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20201119102817-f84b799fce68/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
 golang.org/x/sys v0.0.0-20210119212857-b64e53b001e4/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
-golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c h1:VwygUrnw9jn88c4u8GD3rZQbqrP/tgas88tPUbBxQrk=
 golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
+golang.org/x/sys v0.1.0 h1:kunALQeHf1/185U1i0GOB/fy1IPRDDpuoOOqRReG57U=
+golang.org/x/sys v0.1.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
 golang.org/x/term v0.0.0-20210422114643-f5beecf764ed h1:Ei4bQjjpYUsS4efOUz+5Nz++IVkHk87n2zBA0NxBWc0=
 golang.org/x/term v0.0.0-20210422114643-f5beecf764ed/go.mod h1:bj7SfCRtBDWHUb9snDiAeCFNEtKQo2Wmx5Cou7ajbmo=
 gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=`

type File struct {
	Path     string
	Hunks    []Hunk
	Preamble []string
}
type Hunk struct {
	LeftStart  int64
	LeftCount  int64
	RightStart int64
	RightCount int64
	lines      []string
	Lines      []Line
}
type Line struct {
	Raw         string
	LeftNumber  int64
	RightNumber int64
	Left        bool
	Right       bool
}

func main() {
	files := utils.ParseDiffText(diff)
	fmt.Println(renderUnifiedDiff(files))
}

func renderUnifiedDiff(files []utils.File) string {
	doc := strings.Builder{}
	for _, file := range files {
		body := strings.Builder{}
		header := strings.Builder{}
		header.WriteString(file.Path + "\n")
		for _, line := range file.Preamble {
			header.WriteString(line + "\n")
		}
		for _, hunk := range file.Hunks {
			for _, line := range hunk.Lines {
				var lf string
				if line.Left {
					lf = components.DeletionsStyle.Render(line.Raw)
				} else if line.Right {
					lf = components.AdditionsStyle.Render(line.Raw)
				} else {
					lf = line.Raw
				}
				lp := components.DiffLineNumberStyle.Render(fmt.Sprintf("%d,%d", line.LeftNumber, line.RightNumber))
				body.WriteString(fmt.Sprintf("%s %s\n", lp, lf))

			}
			body.WriteString("\n--------------------------------------\n")
		}
		doc.WriteString(components.RenderBoxWithTitle(header.String(), body.String(), 160))
		doc.WriteString("\n")
	}
	return doc.String()
}

func parseDiffText(diff string) []File {
	hunkRegex := regexp.MustCompile(`@@ -(\d+),(\d+) \+(\d+),(\d+) @@`)
	pathRegex := regexp.MustCompile("^ a/(.*) b/")

	files := []File{}

	for _, fs := range strings.Split(diff, "diff --git")[1:] {
		f := File{}
		paths := pathRegex.FindStringSubmatch(fs)
		path := paths[1]
		f.Path = path
		f.Preamble = append(f.Preamble, strings.Split(fs, "\n")[1:4]...)
		hunkLocs := hunkRegex.FindAllStringSubmatch(fs, -1)
		hunks := hunkRegex.Split(fs, -1)
		for i, hunkLoc := range hunkLocs {
			h := Hunk{}
			var m = []int64{}
			for i := 1; i <= 4; i++ {
				j, err := strconv.ParseInt(string(hunkLoc[i]), 10, 64)
				if err != nil {
					log.Fatal(err)
				}
				m = append(m, j)
			}
			h.LeftStart = m[0]
			h.LeftCount = m[1]
			h.RightStart = m[2]
			h.RightCount = m[3]
			h.lines = strings.Split(hunks[i+1], "\n")[1:]
			h.populateLines()
			f.Hunks = append(f.Hunks, h)
		}
		files = append(files, f)
	}
	return files
}

func (h *Hunk) populateLines() {
	var rc int64
	var lc int64
	var rs string
	var ls string
	var li int64
	var ri int64
	lines := []Line{}
	for _, l := range h.lines {
		line := Line{}
		if strings.HasPrefix(l, "-") {
			lc = h.LeftStart + li
			ls = fmt.Sprintf("%d", lc)
			rs = strings.Repeat(" ", len(ls))
			li += 1
			line.Left = true
		} else if strings.HasPrefix(l, "+") {
			rc = h.RightStart + ri
			ls = strings.Repeat(" ", len(rs))
			rs = fmt.Sprintf("%d", rc)
			ri += 1
			line.Right = true
		} else {
			lc = h.LeftStart + li
			rc = h.RightStart + ri
			ls = fmt.Sprintf("%d", lc)
			rs = fmt.Sprintf("%d", rc)
			li += 1
			ri += 1
		}
		line.LeftNumber = li
		line.RightNumber = ri
		line.Raw = fmt.Sprintf("%s,%s", ls, rs)
		lines = append(lines, line)
	}
	h.Lines = lines
}

func (h Hunk) renderUnifiedDiff(width int) string {
	doc := strings.Builder{}
	doc.WriteString("```diff\n")
	for _, line := range h.lines {
		doc.WriteString(line)
		doc.WriteString("\n")
	}
	doc.WriteString("\n```")
	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	body, err := r.Render(doc.String())
	if err != nil {
		body = "ERROR"
	}
	ll := strings.Split(body, "\n")
	_ = ll
	nl := []string{}
	var rc int64
	var lc int64
	var rs string
	var ls string
	var lf string
	var li int64
	var ri int64
	for _, l := range h.lines {
		if strings.HasPrefix(l, "-") {
			lc = h.LeftStart + li
			ls = fmt.Sprintf("%d", lc)
			rs = strings.Repeat(" ", len(ls))
			li += 1
			lf = components.DeletionsStyle.Render(l)
		} else if strings.HasPrefix(l, "+") {
			rc = h.RightStart + ri
			ls = strings.Repeat(" ", len(rs))
			rs = fmt.Sprintf("%d", rc)
			ri += 1
			lf = components.AdditionsStyle.Render(l)
		} else {
			lc = h.LeftStart + li
			rc = h.RightStart + ri
			ls = fmt.Sprintf("%d", lc)
			rs = fmt.Sprintf("%d", rc)
			li += 1
			ri += 1
			lf = l
		}
		lp := components.DiffLineNumberStyle.Render(fmt.Sprintf("%s,%s", ls, rs))
		nl = append(nl, fmt.Sprintf("%s %s", lp, lf))
	}
	return strings.Join(nl, "\n")
}
