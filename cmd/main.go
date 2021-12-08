package main

import (
	"flag"
	"log"

	"github.com/goplus/build/cmd/internal/builder"
)

var (
	verbose = flag.Bool("v", false, "verbose output")
	repo    = flag.String("repo", "https://github.com/goplus/gop.git", "repo git")
	branch  = flag.String("b", "", "branch")
	o       = flag.String("o", ".", "out put path")
)

func main() {
	flag.Parse()

	var b builder.Build
	// TODO: add os and archs
	b.OS = "darwin"
	b.Arch = ""
	b.Repo = *repo
	b.Branch = *branch
	b.Verbose = *verbose
	b.Out = *o

	if err := b.Do(); err != nil {
		log.Printf("err: %v", err)
	}

}
