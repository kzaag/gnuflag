package gnuflag

import (
	"strings"
)

func prepFmt(fmt []string) map[string]string {

	// alpanumeric charset as per GNU specification
	const argCharset = "abcdefghijklmnopqrstuvwxyz0123456789"

	mfmt := make(map[string]string)
O:
	for i := range fmt {
		v := fmt[i]
		for j, x := range v {
			if !strings.Contains(argCharset, strings.ToLower(string(x))) {
				// ':' at the end is allowed
				if j == len(v)-1 && x == ':' {
					break
				}
				// ignore this param
				continue O
			}
		}
		if strings.HasSuffix(v, ":") {
			l := len(v)
			mfmt[v[:l-1]] = ":"
		} else {
			mfmt[v] = ""
		}
	}
	return mfmt
}

// Getopt parses options and arguments according to package specification.
//
// argv is argument array (for example os.Args).
//
// _fmt is format array.
//
// optcb will be called every time function encounters argument.
//
// If argument is option present within _fmt array then value opt will be set to option name, and value optarg to option value (empty if not exists).
//
// If argument is option not present within _fmt array then opt will be set to "?" and optarg will contain option name.
//
// If argument is not an option then opt will be set to empty string and optarg will contain value.
//
// optcb may return false to terminate Getopt function immediately.
//
// see flag_test.go for examples.
func Getopt(
	argv []string,
	optcb func(opt, optarg string) bool,
	_fmt ...string,
) {
	var currentFlag string
	var isTerminated bool
	var isUnrecognizedOpt bool

	mfmt := prepFmt(_fmt)

	for i := range argv {
		e := argv[i]

		if isTerminated {
			if !optcb("", e) {
				return
			}
			continue
		}

		// we are still reading previous flag
		if currentFlag != "" {
			if !optcb(currentFlag, e) {
				return
			}
			currentFlag = ""
			continue
		}

		// not an option
		if !strings.HasPrefix(e, "-") {
			// if previous option was not recognized then ignore this value
			if isUnrecognizedOpt {
				isUnrecognizedOpt = false
				continue
			}

			if !optcb("", e) {
				return
			}

			continue
		}

		// is an option

		isUnrecognizedOpt = false

		// single hyphen - treat as an argument
		if len(e) == 1 {
			if !optcb("", "-") {
				return
			}
			continue
		}

		// long option
		if strings.HasPrefix(e, "--") {
			name := e[2:]

			if name == "" {
				// this is terminator
				currentFlag = ""
				isTerminated = true
				continue
			}

			if strings.Contains(name, "=") { // may be --foo=bar format
				parts := strings.SplitN(name, "=", 2)
				if len(parts) != 2 {
					continue
				}
				v, f := mfmt[parts[0]]
				if !f || v != ":" {
					isUnrecognizedOpt = true
					if !optcb("?", parts[0]) {
						return
					}
					continue
				}
				if !optcb(parts[0], parts[1]) {
					return
				}

				continue
			}

			// either bool or value in the next argument

			v, f := mfmt[name]
			if !f {
				isUnrecognizedOpt = true
				if !optcb("?", name) {
					return
				}
				continue
			}

			if v == ":" {
				currentFlag = name
			} else {
				if !optcb(name, "") {
					return
				}
			}
			continue
		}

		// short option

		for j := 1; j < len(e); j++ {
			nameb := byte(e[j])
			name := string([]byte{nameb})
			v, f := mfmt[name]
			if !f {
				isUnrecognizedOpt = true
				if !optcb("?", name) {
					return
				}
				continue
			}
			if v == ":" { // has value
				if j < len(e)-1 { // has value within this argv element
					if !optcb(name, e[j+1:]) {
						return
					}
					break
				} else { // has value in next argv element
					currentFlag = name
				}
			} else { // is bool
				if !optcb(name, "") {
					return
				}
			}
		}
	}

	if currentFlag != "" {
		// we are still reading previous flag
		optcb(currentFlag, "")
	}
}
