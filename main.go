package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
)

var memberRegs map[string]string

type linegroup struct {
	sortkey string
	lines   []string
}

func init() {
	memberRegs = map[string]string{
		"C#": `^( +)public (int|decimal)\?? ([a-z0-9_]+) .*$`,
	}
}

func main() {
	infile := flag.String("in", "in.txt", "input file")
	outfile := flag.String("out", "out.txt", "output file")
	lang := flag.String("lang", "C#", "language")
	flag.Parse()

	if *infile == "" || *outfile == "" {
		flag.PrintDefaults()
		return
	}

	var reg *regexp.Regexp
	if rs, ok := memberRegs[*lang]; ok {
		reg = regexp.MustCompile(rs)
	} else {
		log.Println("unsupported language at this time %s", *lang)
		return
	}

	lines, err := getLines(*infile)
	if err != nil {
		log.Fatal(err)
	}

	groups := getGroups(reg, lines)
	sortGroups(groups)
	err = writeGroups(*outfile, groups)
	if err != nil {
		log.Fatal(err)
	}
}

func writeGroups(outfile string, groups []linegroup) error {
	fout, err := os.OpenFile(outfile, os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("couldn't open file for writing %s. %w", outfile, err)
	}
	defer fout.Close()

	for _, g := range groups {
		for _, line := range g.lines {
			fmt.Fprintln(fout, line)
		}
	}
	return nil
}

func sortGroups(groups []linegroup) {
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].sortkey < groups[j].sortkey
	})
}

func getGroups(memberReg *regexp.Regexp, lines []string) []linegroup {
	groups := []linegroup{}
	for i := len(lines) - 1; i >= 0; i-- {
		s := lines[i]
		mm := memberReg.FindAllStringSubmatch(s, -1)
		if len(mm) == 0 {
			if len(groups) > 0 && len(s) > 0 {
				groups[len(groups)-1].lines = append([]string{s}, groups[len(groups)-1].lines...)
			}
			continue
		}

		lg := linegroup{sortkey: mm[0][3], lines: []string{mm[0][0]}}
		groups = append(groups, lg)
	}
	return groups
}

func getLines(filename string) ([]string, error) {
	f, err := os.Open(filename)

	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	lines := []string{}
	for scanner.Scan() {
		var txt = scanner.Text()
		lines = append(lines, txt)
	}

	return lines, nil
}
