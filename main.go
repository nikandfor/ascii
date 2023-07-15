package main

import (
	"fmt"
	"os"

	stdflag "github.com/nikandfor/cli/stdflag"
	"github.com/nikandfor/escape/color"
	"golang.org/x/term"
)

var (
	rows        = stdflag.Int("rows,r", 32, "number of rows (16, 32)")
	colorize    = stdflag.Bool("colorize,color,c", term.IsTerminal(int(os.Stdout.Fd())), "colorize")
	showDecimal = stdflag.Bool("decimal", true, "show decimal code")
	showHex     = stdflag.Bool("hex", true, "show hex code")
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

	for c := 0; c < *rows; c++ {
		for b := 0; b < 128; b += *rows {
			if b != 0 {
				buf = append(buf, "      "...)
			}

			r := b + c

			buf = append(buf, gray...)

			if *showDecimal {
				buf = fmt.Appendf(buf, "%3d", r)
			}

			if *showDecimal && *showHex {
				buf = append(buf, "  "...)
			}

			if *showHex {
				buf = fmt.Appendf(buf, "%2x", r)
			}

			buf = append(buf, reset...)

			if *showDecimal || *showHex {
				buf = append(buf, "  "...)
			}

			if r < 32 || r >= 127 {
				buf = fmt.Appendf(buf, "%-6q", r)
			} else {
				buf = fmt.Appendf(buf, "%c", r)
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
