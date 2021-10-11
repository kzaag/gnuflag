package gnuflag

import (
	"fmt"
	"strings"
	"testing"
)

type entry struct {
	name, val string
}

func testParse(t *testing.T, args string, e []entry, _fmt ...string) {
	fmt.Println("------------")
	fmt.Printf("PARSE: \"%s\"\nFMT=\"%s\"\n", args, strings.Join(_fmt, " "))
	found := make(map[entry]struct{}, len(e))
	Getopt(
		strings.Split(args, " "),
		func(name, val string) bool {
			found[entry{
				name: name,
				val:  val,
			}] = struct{}{}
			if name == "" {
				fmt.Printf("ARG: \"%s\" = \"%s\"\n", name, val)
			} else {
				fmt.Printf("OPT: \"%s\" = \"%s\"\n", name, val)
			}
			return true
		},
		_fmt...,
	)
	if e == nil {
		return
	}

	for i := range e {
		if _, f := found[e[i]]; !f {
			t.Fatalf("flag %v not found", e[i])
		}
	}
	for fi := range found {
		f := false
		for i := range e {
			if e[i] == fi {
				f = true
				break
			}
		}
		if !f {
			t.Fatalf("flag %v found but not defined", fi)
		}
	}
}

func TestFlag(t *testing.T) {
	argv := "foo -cba a_val"
	testParse(
		t,
		argv,
		// expected opt-optarg pairs:
		[]entry{
			{"", "foo"},
			{"c", ""},
			{"b", ""},
			{"a", "a_val"},
		},
		// _fmt
		"c", "b", "a:")

	// same argv but a is treated as option without value
	testParse(t, argv, []entry{
		{"c", ""},
		{"b", ""},
		{"a", ""},
		{"", "a_val"},
		{"", "foo"},
	}, "c", "b", "a")

	argv = "foo -cba"
	testParse(t, argv, []entry{
		{"c", "ba"},
		{"", "foo"},
	}, "c:", "b")

	argv = "foo -cba -v"
	testParse(t, argv, []entry{
		{"c", ""},
		{"b", "a"},
		{"v", ""},
		{"", "foo"},
	}, "c", "b:", "v")

	argv = "foo -a10 -b20 --foo bar --h"
	testParse(t, argv, []entry{
		{"a", "10"},
		{"b", "20"},
		{"foo", "bar"},
		{"h", ""},
		{"", "foo"},
	}, "a:", "b:", "foo:", "h")

	argv = "foo -a10 -b20 -- foo bar --h"
	testParse(t, argv, []entry{
		{"a", "10"},
		{"b", "20"},
		{"", "foo"},
		{"", "bar"},
		{"", "--h"},
	}, "a:", "b:", "foo:", "h")

	argv = "foo -a10 -ś20"
	testParse(t, argv, []entry{
		{"a", "10"},
		{"", "foo"},
		// not found options: unicode ś
		{"?", "\xC5"},
		{"?", "\x9B"},
		// not found options: 20
		{"?", "2"},
		{"?", "0"},
	}, "a:", "ś:")

	argv = "foo -a10 - -f bar"
	testParse(t, argv, []entry{
		{"a", "10"},
		{"", "foo"},
		{"", "-"},
		{"f", "bar"},
	}, "a:", "f:")

	argv = "foo -a统一码 --"
	testParse(t, argv, []entry{
		{"a", "统一码"},
		{"", "foo"},
	}, "a:", "f:")

	argv = "x --foo=bar --h=sad --f="
	testParse(t, argv, []entry{
		{"", "x"},
		{"foo", "bar"},
		{"h", "sad"},
		{"?", "f"},
	}, "foo:", "h:", "f")

	argv = "x --foo bar --h=sad --f="
	testParse(t, argv, []entry{
		{"", "x"},
		{"foo", "bar"},
		{"h", "sad"},
		{"?", "f"},
	}, "foo:", "h:", "f")

	argv = "program --fo x somevalue --bar y"
	testParse(t, argv, []entry{
		{"", "program"},
		{"?", "fo"}, // note x is ignored since it belongs to option which is not found
		{"", "somevalue"},
		{"bar", "y"},
	}, "foo:", "bar:")

	// missing arguments will cause following opts to be interpreted as arguments

	argv = "program -n -t -a"
	testParse(t, argv, []entry{
		{"", "program"},
		{"t", "-a"},
		{"n", ""},
	}, "n", "a", "t:")
	argv = "program -n -ta"
	testParse(t, argv, []entry{
		{"", "program"},
		{"t", "a"},
		{"n", ""},
	}, "n", "a", "t:")
}
