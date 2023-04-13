package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/blugelabs/bluge"
)

func main() {
	indexDir := flag.String("index", "/Users/ali/workspace/data/bluge/index/99p/", "where to load the data from")
	flag.Parse()
	cfg := bluge.DefaultConfig(*indexDir)
	indexReader, err := bluge.OpenReader(cfg)
	if err != nil {
		log.Fatalf("unable to open reader: %v", err)
	}
	defer func() {
		err = indexReader.Close()
		if err != nil {
			log.Fatalf("error closing reader: %v", err)
		}
	}()
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Press Ctrl-D to exit.")
	fmt.Printf("query:")
	for scanner.Scan() {
		query := scanner.Text()
		fmt.Println("--------------------------------")
		fmt.Println(query)
		fmt.Println("--------------------------------")
		fmt.Printf("query:")
	}
	if err := scanner.Err(); err != nil {

		log.Fatalln("oops!", err)
	}
	// defer func(tStart time.Time) {
	// 	fmt.Println("Elapsed Time:", time.Since(tStart))
	// }(time.Now())

	// q := bluge.NewMatchQuery("pto").SetField("title").SetFuzziness(1)

	// req := bluge.NewTopNSearch(10, q)

	// dmi, err := indexReader.Search(context.Background(), req)
	// if err != nil {
	// 	log.Fatalf("error executing search: %v", err)
	// }

	// next, err := dmi.Next()
	// for err == nil && next != nil {
	// 	err = next.VisitStoredFields(func(field string, value []byte) bool {
	// 		if field == "title" {
	// 			fmt.Println(string(value))
	// 		}
	// 		fmt.Println(field)
	// 		return true
	// 	})
	// 	if err != nil {
	// 		log.Fatalf("error accessing stored fields: %v", err)
	// 	}
	// 	fmt.Println("----------------------------------")
	// 	next, err = dmi.Next()
	// }
	// if err != nil {
	// 	log.Fatalf("error iterating results: %v", err)
	// }

}
