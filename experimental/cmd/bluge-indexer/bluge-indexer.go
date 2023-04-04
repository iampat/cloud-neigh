package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/blugelabs/bluge"
)

const maxLineSize int = 1000 * 1000 // Reserve 1MB
const vebosity int = 500 * 1000
const maxNumberOfItems = 10 * 1000 * 1000

func main() {
	defer func(tStart time.Time) {
		fmt.Println("Elapsed Time:", time.Since(tStart))
	}(time.Now())
	// Read data file
	datasetName := "20220301_en_1k"
	fname := "/Users/ali/workspace/data/wikipedia/" + datasetName + ".json"
	indexDir := "/Users/ali/workspace/data/bluge/index/"
	const keepIndex = true
	var err error
	if keepIndex {
		indexDir = path.Join(indexDir, "keep", datasetName)
		err = os.MkdirAll(indexDir, 0750)
	} else {
		indexDir, err = os.MkdirTemp(indexDir, "wiki_"+datasetName+"_*")
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Println("input json file:", fname)
	readFile, err := os.Open(fname)
	if err != nil {
		fmt.Println(err)
	}
	defer readFile.Close()
	fileScanner := bufio.NewScanner(readFile)

	cfg := bluge.DefaultConfig(indexDir)
	indexWriter, err := bluge.OpenWriter(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		cerr := indexWriter.Close()
		if cerr != nil {
			log.Fatal(cerr)
		}
	}()

	type Document struct {
		Id    string `json:"id"`
		Title string `json:"title"`
		Url   string `json:"url"`
		Text  string `json:"text"`
	}

	counter := 0
	batch := bluge.NewBatch()
	fileScanner.Buffer(make([]byte, maxLineSize), maxLineSize)
	for fileScanner.Scan() {

		data := fileScanner.Bytes()
		inputDoc := &Document{}
		json.Unmarshal(data, inputDoc)
		bdoc := bluge.NewDocument(inputDoc.Id).
			AddField(bluge.NewTextField("title", inputDoc.Title).StoreValue()).
			AddField(bluge.NewTextField("url", inputDoc.Url).StoreValue()).
			AddField(bluge.NewTextField("text", inputDoc.Text))

		batch.Update(bdoc.ID(), bdoc)

		counter++
		if counter%vebosity == 0 {
			err = indexWriter.Batch(batch)
			if err != nil {
				log.Fatalf("error executing batch: %v", err)
			}
			batch.Reset()
			log.Printf("%d items have been added", counter)
		}
		if counter == maxNumberOfItems {
			break
		}
	}
	if fileScanner.Err() != nil {
		log.Panicln(fileScanner.Err())
	}

	log.Printf("document indexed: %s", indexDir)
}
