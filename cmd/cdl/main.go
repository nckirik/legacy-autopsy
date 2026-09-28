// Command cdl exposes the Canonical Deconstruction Language toolchain.
//
// Usage:
//
//	cdl fmt [-w] [--check] [files...]
//
// With no files, input is read from stdin and written to stdout. -w rewrites
// files in place; --check reports non-canonical files and exits 1.
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/nckirik/legacy-autopsy/cdl"
)

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "cdl:", err)
		os.Exit(1)
	}
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: cdl fmt [-w] [--check] [files...]")
	}
	switch args[0] {
	case "fmt":
		return runFmt(args[1:], stdin, stdout, stderr)
	case "version":
		fmt.Fprintf(stdout, "cdl %s (eir-format %d)\n", cdl.LanguageVersion, cdl.EIRFormat)
		return nil
	default:
		return fmt.Errorf("unsupported command %q; supported: fmt, version", args[0])
	}
}

func diagError(d cdl.Diagnostic) error {
	if d.Line > 0 {
		return fmt.Errorf("%s:%d: %s: %s", d.File, d.Line, d.Code, d.Msg)
	}
	return fmt.Errorf("%s: %s: %s", d.File, d.Code, d.Msg)
}

func runFmt(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	write := false
	check := false
	var files []string
	for _, arg := range args {
		switch arg {
		case "-w", "--write":
			write = true
		case "--check":
			check = true
		default:
			files = append(files, arg)
		}
	}
	if write && check {
		return fmt.Errorf("fmt: -w and --check are mutually exclusive")
	}
	if len(files) == 0 {
		if write || check {
			return fmt.Errorf("fmt: -w/--check require file arguments")
		}
		src, err := io.ReadAll(stdin)
		if err != nil {
			return err
		}
		formatted, diags := cdl.FormatSource("<stdin>", src)
		if len(diags) > 0 {
			return diagError(diags[0])
		}
		_, err = stdout.Write(formatted)
		return err
	}
	nonCanonical := 0
	for _, path := range files {
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		formatted, diags := cdl.FormatSource(path, src)
		if len(diags) > 0 {
			return diagError(diags[0])
		}
		switch {
		case check:
			if string(formatted) != string(src) {
				fmt.Fprintln(stderr, path)
				nonCanonical++
			}
		case write:
			if string(formatted) != string(src) {
				if err := os.WriteFile(path, formatted, 0o644); err != nil {
					return err
				}
				fmt.Fprintln(stderr, "formatted", path)
			}
		default:
			if _, err := stdout.Write(formatted); err != nil {
				return err
			}
		}
	}
	if nonCanonical > 0 {
		return fmt.Errorf("%d file(s) are not canonically formatted (run: cdl fmt -w)", nonCanonical)
	}
	return nil
}
