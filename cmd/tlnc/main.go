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
		fmt.Println("tlnc", buildinfo.String())
		return
	}

	fmt.Println("tlnc: phase 0 - not implemented yet")
}
