package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var rootPath = flag.String("path", ".", "Root path")

func main() {
	flag.Parse()

	var files list

	filepath.Walk(*rootPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		files = append(files, newNatStr(path))
		return nil
	})

	sort.Sort(files)

	for i, file := range files {
		path := file.s
		fmt.Println(path)
		dir, file := filepath.Split(path)
		if err := os.Rename(path, filepath.Join(dir, fmt.Sprintf("s01e%d.%s", i+1, file))); err != nil {
			log.Fatal(err)
		}
	}
}

// natStr associates a string with a preprocessed form
type natStr struct {
	s string // original
	t []tok  // preprocessed "sub-fields"
}

func newNatStr(s string) (t natStr) {
	t.s = s
	s = strings.ToLower(strings.Join(strings.Fields(s), " "))
	x := dx.FindAllString(s, -1)
	t.t = make([]tok, len(x))
	for i, s := range x {
		if n, err := strconv.Atoi(s); err == nil {
			t.t[i].n = n
		} else {
			t.t[i].s = s
		}
	}
	return t
}

var dx = regexp.MustCompile(`\d+|\D+`)

// rule is to use s unless it is empty, then use n
type tok struct {
	s string
	n int
}

// rule 2 of "numeric sub-fields" from talk page
func (f1 tok) Cmp(f2 tok) int {
	switch {
	case f1.s == "":
		switch {
		case f2.s > "" || f1.n < f2.n:
			return -1
		case f1.n > f2.n:
			return 1
		}
	case f2.s == "" || f1.s > f2.s:
		return 1
	case f1.s < f2.s:
		return -1
	}
	return 0
}

type list []natStr

func (l list) Len() int      { return len(l) }
func (l list) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l list) Less(i, j int) bool {
	ti := l[i].t
	for k, t := range l[j].t {
		if k == len(ti) {
			return true
		}
		switch ti[k].Cmp(t) {
		case -1:
			return true
		case 1:
			return false
		}
	}
	return false
}
