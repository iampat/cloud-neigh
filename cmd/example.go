package main

import (
	"fmt"
	"log"
	"net/http"
	"gonum.org/v1/hdf5"
)

func main() {
	fmt.Println("Running cloudy neigh single local machine example ☁🏇☁🏇☁🏇")
	fmt.Println("Loads Dataset & Queries ...")
	fmt.Println("=== go-hdf5 ===")
	version, err := hdf5.LibVersion()
	if err != nil {
		fmt.Printf("** error ** %s\n", err)
		return
	}
	fmt.Printf("=== version: %s", version)
	fmt.Println("=== bye.")
}
	// gloveAngular100Url := "http://ann-benchmarks.com/glove-100-angular.hdf5"
	// resp, err := http.Get(gloveAngular100Url)
	// if err != nil {
	// 	log.Fatalln("failed to load the dataset", err)
	// }
	// resp.Body.Read()
	// fmt.Println("Creates Searcher ...")
	// fmt.Println("Runs Partitioning ...")
	// fmt.Println("Runs Scoring ...")
	// fmt.Println("Runs Rescoring ...")
}