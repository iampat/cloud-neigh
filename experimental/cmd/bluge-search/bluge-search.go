package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/blugelabs/bluge"
)

func main() {
	// indexDir := "/Users/ali/workspace/data/bluge/index/" + "keep/wiki_en_full"
	indexDir := "/Users/ali/workspace/data/bluge/index/keep/20220301_en_1k"
	cfg := bluge.DefaultConfig(indexDir)
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

	defer func(tStart time.Time) {
		fmt.Println("Elapsed Time:", time.Since(tStart))
	}(time.Now())

	q := bluge.NewMatchQuery("alisp").SetField("title").SetFuzziness(1)

	req := bluge.NewTopNSearch(10, q)

	dmi, err := indexReader.Search(context.Background(), req)
	if err != nil {
		log.Fatalf("error executing search: %v", err)
	}

	next, err := dmi.Next()
	for err == nil && next != nil {
		err = next.VisitStoredFields(func(field string, value []byte) bool {
			if field == "title" {
				fmt.Println(string(value))
			}

			return true
		})
		if err != nil {
			log.Fatalf("error accessing stored fields: %v", err)
		}
		fmt.Println("----------------------------------")
		next, err = dmi.Next()
	}
	if err != nil {
		log.Fatalf("error iterating results: %v", err)
	}

}
