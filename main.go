package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

func noError(e error) {
	if e != nil {
		panic(e)
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}

var (
	output = flag.String("out", "", "output file")
	quiet = flag.Bool("q", false, "quite")
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "needs exactly one file")
		os.Exit(1)
	}

	fin, e := os.Open(args[0])
	url := "https://api.github.com/markdown/raw"
	resp, e := http.Post(url, "text/plain", fin)
	noError(e)
	if !*quiet {
		fmt.Fprintln(os.Stderr, resp.Status)
	}

	var fout = os.Stdout
	if *output != "" {
		fout, e = os.Create(*output)
		noError(e)
	}

	_, e = io.Copy(fout, resp.Body)
	noError(e)

	noError(fout.Close())
}
