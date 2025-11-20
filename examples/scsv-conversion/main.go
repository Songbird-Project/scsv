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

	os.WriteFile("pkglist.csv", []byte(scsv.AsCSV(scsvMap)), 0644)
	os.WriteFile("pkglist.json", []byte(scsv.AsJSON(scsvMap)), 0644)
	os.WriteFile("pkglist.yaml", []byte(scsv.AsYAML(scsvMap)), 0644)
}
