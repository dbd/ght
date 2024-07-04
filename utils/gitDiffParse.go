package utils

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

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
		line.LeftNumber = lc
		line.RightNumber = rc
		line.Raw = fmt.Sprintf("%s,%s", ls, rs)
		line.Raw = l
		lines = append(lines, line)
	}
	h.Lines = lines
}

func ParseDiffText(diff string) []File {
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

func ParseHunkDiff(hunk string) Hunk {
	hunkRegex := regexp.MustCompile(`@@ -(\d+),(\d+) \+(\d+),(\d+) @@`)
	hunkLocs := hunkRegex.FindAllStringSubmatch(hunk, -1)
	hunks := hunkRegex.Split(hunk, -1)
	h := Hunk{}
	for i, hunkLoc := range hunkLocs {
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
		return h
	}
	return h
}
