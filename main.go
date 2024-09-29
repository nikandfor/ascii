package main

import (
	"fmt"
	"os"

	"github.com/nikandfor/escape/color"
	"golang.org/x/term"
	stdflag "nikand.dev/go/cli/stdflag"
)

var (
	cols        = stdflag.Int("cols,c", 4, "number of columns (4, 8)")
	colorize    = stdflag.Bool("colorize,color,c", term.IsTerminal(int(os.Stdout.Fd())), "colorize")
	showDecimal = stdflag.Bool("decimal", true, "show decimal code")
	showHex     = stdflag.Bool("hex", true, "show hex code")
	showAbbr    = stdflag.Bool("abbr", true, "show control chars names abbreveations")
)

func main() {
	stdflag.CommandLine.EnvPrefix = "ASCII_"
	stdflag.Parse()

	var gray, reset []byte

	if *colorize {
		gray = color.New(90)
		reset = color.New(0)
	}

	var buf []byte

	rows := 128 / *cols
	if *cols*rows < 128 {
		rows++
	}

	for r := 0; r < rows; r++ {
		for c := r; c < 128; c += rows {
			if c != r {
				buf = append(buf, "      "...)
			}

			buf = append(buf, gray...)

			if *showDecimal {
				buf = fmt.Appendf(buf, "%3d", c)
			}

			if *showDecimal && *showHex {
				buf = append(buf, "  "...)
			}

			if *showHex {
				buf = fmt.Appendf(buf, "%2x", c)
			}

			buf = append(buf, reset...)

			if *showDecimal || *showHex {
				buf = append(buf, "  "...)
			}

			if c < 32 || c >= 127 {
				buf = fmt.Appendf(buf, "%-6q", c)
			} else {
				buf = fmt.Appendf(buf, "%c", c)
			}

			switch {
			case c < 32 && *showAbbr:
				buf = append(buf, gray...)
				buf = fmt.Appendf(buf, "  %-3s", names[c])
				buf = append(buf, reset...)
			case c == 0x7f && *showAbbr:
				buf = append(buf, gray...)
				buf = fmt.Appendf(buf, "  del")
				buf = append(buf, reset...)
			case c-c%rows < 32 && *showAbbr:
				buf = fmt.Appendf(buf, "          ")
			}
		}

		buf = append(buf, '\n')
	}

	_, err := os.Stdout.Write(buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

var names = []string{
	"nul", "soh", "stx", "etx",
	"eot", "enq", "ack", "bel",
	"bs", "ht", "lf", "vt",
	"ff", "cr", "so", "si",
	"dle", "dc1", "dc2", "dc3",
	"dc4", "nak", "syn", "etb",
	"can", "em", "sub", "esc",
	"fs", "gs", "rs", "us",
}
