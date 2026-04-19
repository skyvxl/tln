package main

import (
	"flag"
	"fmt"

	"github.com/skyvxl/tln/internal/buildinfo"
)

func main() {
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println("tlnd", buildinfo.String())
		return
	}

	fmt.Println("tlnd: phase 0 - not implemented yet")
}
