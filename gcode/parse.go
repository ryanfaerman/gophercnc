package gcode

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var (
	codeRegex = regexp.MustCompile(`([A-Z])(\d+(\.\d+)?)`)
)

// Parse works through the reader line-wise, converting each line into the
// correct code values. It attempts to use the correct addresses and vaues that
// will match the original input. Addresses (G,M, T, etc.) are expected to be
// upper case, otherwise they will be considered as string arguments.
//
// Empty lines are ignored and will not be present in the resulting output of
// Fragment.String().
func Parse(r io.Reader) (Fragment, error) {
	var f Fragment

	s := bufio.NewScanner(r)
	s.Split(bufio.ScanLines)

	for s.Scan() {
		line := strings.TrimSpace(s.Text())

		if len(line) == 0 {
			continue
		}

		fields := strings.Fields(line) // this gets us the "tokens" of the line

		c := code{}
		for i, field := range fields {
			current := code{}

			// Once we have the prefix for a comment, the rest of the line is a
			// comment. This block attempts to consume the rest of the line as
			// soon as that prefix is detected.
			if strings.HasPrefix(field, string(Comment.Address())) {
				str := strings.TrimLeft(strings.Join(fields[i:], " "), "; ")
				current = Comment(str)

				if i > 0 {
					c.codes = append(c.codes, current)
				} else {
					c = current
				}

				break
			}

			possibleCodes := codeRegex.FindStringSubmatch(field)

			// We have something that resembles G92.3, where the address (G)
			// and the value (92.3) are the matches from above.
			if len(possibleCodes) >= 1 {
				current.address = []byte(possibleCodes[1])[0]

				val, err := strconv.ParseFloat(possibleCodes[2], 64)
				if err != nil {
					return f, err
				}
				current.value = val
			} else {
				// If we can't be sure we have a proper code, we'll drop it into
				// a string. This may not be correct, but at least it won't be
				// wrong. The output will be the same when converted to a string
				// form.
				current.address = String.Address()
				current.comment = field
			}

			// If this isn't the first, that means it is a child of a code --
			// G0 X73.3 Y34.1 -- G0 is the c above and X73.3 would be the
			// current code.
			if i > 0 {
				c.codes = append(c.codes, current)
			} else {
				c = current
			}
		}

		// Our line is fully read and converted into codes, add it to the
		// output fragment
		f = append(f, c)

	}

	// spw := spew.ConfigState{ContinueOnMethod: true}
	// spw.Dump(f)

	return f, nil
}

// ParseString is a helper function. Rather than a reader, it takes a string
// and turns it into a reader for you.
func ParseString(input string) (Fragment, error) {
	return Parse(strings.NewReader(input))
}
