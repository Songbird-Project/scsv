package main

import (
	"fmt"
	"os"

	"github.com/Songbird-Project/scsv"
)

func main() {
	path := "./pkglist.scsv"
	scsvMap, err := scsv.ParseFile(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for repo, pkgs := range scsvMap {
		fmt.Printf("\n%s:\n", repo)

		for _, pkg := range pkgs {
			fmt.Printf(" - %s\n", pkg)
		}
	}
}
