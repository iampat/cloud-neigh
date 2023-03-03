package main

import (
	"fmt"
	"log"

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
	fname := "data/glove-100-angular.hdf5"
	f, err := hdf5.OpenFile(fname, hdf5.F_ACC_RDONLY)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	num, err := f.NumObjects()
	if err != nil {
		log.Panicln("num objects:", err)
	}
	fmt.Println("number of object", num)
	for idx:= uint(0); idx < num; idx++ {
		name, _ := f.ObjectNameByIndex(idx)
		fmt.Println(idx, ":", name)
	}
	dsname := "test"
	dset, err := f.OpenDataset(dsname)
	if err != nil {
		panic(err)
	}
	defer dset.Close()
	dtype, err := dset.Datatype()
	if err != nil {
		panic(err)
	}
	dtype.Close()
	fmt.Println("datatype:", dtype.GoType()	)
	// dset, err := f.OpenDataset(dsname)
	// defer dset.Close()
	// if err != nil {
	// 	panic(err)
	// }

	// // read it back into a new slice
	// s2 := make([]s1Type, length)
	// err = dset.Read(&s2)
	// if err != nil {
	// 	panic(err)
	// }

	// // display the fields
	// fmt.Printf(":: data: %v\n", s2)

	// release resources
//s	hdf5.

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