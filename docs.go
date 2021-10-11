// package gnuflag parses command line arguments, using following POSIX conventions
//
// - options are specified in format array
// 	example: format := []string{"a", "b"} - 'a' and 'b' are recognized options
//
// - wthin argv options begin with single hyphen
// 	example: argv := []string{"-a"} is option 'a'
//
// - multiple options may follow hyphen within single token:
// 	-example: argv := []string{"-abc"} is equivalent to argv := []string{"-a", "-b", "-c"}
//
// - option names must be alphanumeric characters; other values are ignored
// 	example: format := []string{"ś", "b"} - ś will be promptly ignored
//
// - options may have arguments attached
// 	example: argv := []string{"-a", "foobar"} is option 'a' with value 'foobar'
//
// - to specify that option should have argument, append ':' to option
// 	example: format := []string{"a:"} specifies option 'a' which has an argument
//
// - Option and its argument may or may not appear as separate token
// 	example: argv := []string{"-a", "foobar"} is equivalent to argv := []string{"-afoobar"}
//
// - an argument '--' terminates all options
// 	example: argv := []string{"-a", "--", "-b"} b will be treated as an argument, even if its contained withing format
//
// - single hyphen is interpreted as single argument not an option
// 	example: argv := []string{"-a", "-"} -> a is an option, '-' is an non-option argument
//
// - long options are options which consist of > 1 characters, and begin with double hyphen ('--')
// 	example: argv := []string{"--a"} is option 'a'
//
// - long option and its argument may appear as single token if they are separated by equal sign ('=')
// 	example 1: argv := []string{"--a=b"} is option 'a' with value 'b'
// 	example 2: argv := []string{"--a", "b"} is equivalent to example 1
//
//
// more @ https://www.gnu.org/software/libc/manual/html_node/Argument-Syntax.html
//
package gnuflag
