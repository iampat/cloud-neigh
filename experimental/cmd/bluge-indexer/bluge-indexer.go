package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github/iampat/cloudy-neigh/document"
	"log"
	"os"
	"time"

	"github.com/blugelabs/bluge"
)

const maxLineSize int = 1000 * 1000 // Reserve 1MB
const vebosity int = 100
const maxNumberOfItems = 10 * 1000 * 1000

const lshNumberOfBuckets = 20 // 2^20 buckets
const maxLshNumberOfBuckets = 64
const hashPattern = "%064b"

func dummyHash(num uint64, numberOfBuckets int) string {
	if numberOfBuckets > maxLshNumberOfBuckets {
		log.Panicf("numberOfBuckets (%d) > maxLshNumberOfBuckets(%d)", numberOfBuckets, maxLshNumberOfBuckets)
	}
	num = num % (uint64(1) << numberOfBuckets)
	return fmt.Sprintf(hashPattern, num)
}
func main() {
	defer func(tStart time.Time) {
		fmt.Println("Elapsed Time:", time.Since(tStart))
	}(time.Now())
	var inputJson = flag.String("input_json", "/Users/ali/src/misc/dash_answering/notebooks/data/99p/dataset.json", "where to load the data")
	var outputIndex = flag.String("output_index", "/Users/ali/workspace/data/bluge/index/99p/", "where to write the data")
	flag.Parse()
	// Read data file
	// datasetName := "20220301_en_1k"
	// datasetName := "20220301_en_full"
	// fname := "/Users/ali/workspace/data/wikipedia/" + datasetName + ".json"
	// indexDir := "/Users/ali/workspace/data/bluge/index/"
	// const keepIndex = true
	var err error
	// if keepIndex {
	// 	indexDir = path.Join(indexDir, "keep", datasetName)
	// 	err = os.MkdirAll(indexDir, 0750)
	// } else {
	// 	indexDir, err = os.MkdirTemp(indexDir, "wiki_"+datasetName+"_*")
	// }
	// if err != nil {
	// 	log.Fatal(err)
	// }
	log.Println("input json file:", *inputJson)
	readFile, err := os.Open(*inputJson)
	if err != nil {
		log.Fatalln(err)
	}
	defer readFile.Close()
	fileScanner := bufio.NewScanner(readFile)

	cfg := bluge.DefaultConfig(*outputIndex)
	indexWriter, err := bluge.OpenWriter(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		cerr := indexWriter.Close()
		if cerr != nil {
			log.Fatalln(cerr)
		}
	}()

	counter := 0
	batch := bluge.NewBatch()
	fileScanner.Buffer(make([]byte, maxLineSize), maxLineSize)
	for fileScanner.Scan() {

		data := fileScanner.Bytes()
		inputDoc := &document.Document{}
		// inputDoc.LshHash = dummyHash(rand.Uint64(), lshNumberOfBuckets)
		json.Unmarshal(data, inputDoc)
		bdoc := bluge.NewDocument(inputDoc.Id).
			AddField(bluge.NewTextField("title", inputDoc.Title).StoreValue()).
			AddField(bluge.NewTextField("url", inputDoc.Url).StoreValue()).
			AddField(bluge.NewTextField("text", inputDoc.Text).StoreValue()).
			AddField(bluge.NewTextField("text_lsh_hash", inputDoc.TextLshHash).StoreValue()).
			AddField(bluge.NewTextField("title_lsh_hash", inputDoc.TextLshHash).StoreValue())
		batch.Update(bdoc.ID(), bdoc)

		counter++
		if counter < 5 {
			log.Println("sample doc:", inputDoc.Id, inputDoc.Title, inputDoc.Url)
		}
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
		log.Fatalln(fileScanner.Err())
	}
	err = indexWriter.Batch(batch)
	if err != nil {
		log.Fatalf("error executing batch: %v", err)
	}
	batch.Reset()
	log.Printf("%d items have been added", counter)

	log.Printf("document indexed: %s", *outputIndex)
}
