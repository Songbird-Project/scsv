package main

import (
	"fmt"
	"os"

	"github.com/Songbird-Project/scsv"
)

func main() {
	path := "../pkglist.scsv"
	scsvMap, err := scsv.ParseFile(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for repo, pkgs := range scsvMap {
		if len(pkgs) <= 1 {
			continue
		}

		fmt.Printf("%s:\n", repo)

		for _, pkgColumns := range pkgs {
			for _, pkg := range pkgColumns {
				fmt.Printf(" - %s\n", pkg)
			}
		}
	}
}
