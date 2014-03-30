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
	user = flag.String("user", "", "username")
	pass = flag.String("pass", "", "password")
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "needs exactly one file")
		os.Exit(1)
	}

	fin, e := os.Open(args[0])
	noError(e)

	url := "https://api.github.com/markdown/raw"
	req, e := http.NewRequest("POST", url, fin)
	noError(e)

	req.Header.Set("Content-Type", "text/plain")
	if *user != "" {
		req.SetBasicAuth(*user, *pass)
	}

	client := new(http.Client)
	resp, e := client.Do(req)
	noError(e)

	if !(*quiet || resp.StatusCode == 200) {
		fmt.Fprintln(os.Stderr, resp.Status)
	}

	var fout = os.Stdout
	if *output != "" {
		fout, e = os.Create(*output)
		noError(e)
	}

	_, e = io.Copy(fout, resp.Body)
	noError(e)

	fmt.Fprintln(fout)
	if fout != os.Stdout {
		noError(fout.Close())
	}
}
