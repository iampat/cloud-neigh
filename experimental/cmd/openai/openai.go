package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"github/iampat/cloudy-neigh/document"
	"github/iampat/cloudy-neigh/embeddings"
	"github/iampat/cloudy-neigh/lsh"
	"log"
	"os"
)

const maxLineSize int = 1000 * 1000 // Reserve 1MB
const batchSize = 200
const embeddingDim = 1536
const costDollarPerTocken = float64(0.0004 / 1000.0)
const maxOpenAiTextLength = 10 * 1000
const lshSize = 10

func writeDocsToJson(batch []*document.Document, encoder *json.Encoder) {
	for _, d := range batch {
		if err := encoder.Encode(d); err != nil {
			log.Fatalln(err)
		}
	}
}
func hydraiteBatch(batch []*document.Document, client *embeddings.OpenAIClient, lsh *lsh.Lsh) int {
	titles := []string{}
	contents := []string{}
	for _, doc := range batch {
		t := doc.Title
		if len(t) > maxOpenAiTextLength {
			t = t[:maxOpenAiTextLength]
		}
		titles = append(titles, t)
		c := doc.Text
		if len(c) > maxOpenAiTextLength {
			c = c[:maxOpenAiTextLength]
		}
		contents = append(contents, c)
	}

	contentEmbeddings, cost1, err := client.GetWithUsage(contents)
	if err != nil {
		log.Fatalln("content embedding", err)
	}
	titleEmbeddings, cost2, err := client.GetWithUsage(titles)
	if err != nil {
		log.Fatalln("title embedding", err)
	}
	for idx, doc := range batch {
		doc.TextEmbedding = contentEmbeddings[idx]
		doc.TextLshHash = lsh.Hash(doc.TextEmbedding)
		doc.TitleEmbedding = titleEmbeddings[idx]
		doc.TitleLshHash = lsh.Hash(doc.TitleEmbedding)
	}
	return cost1 + cost2

}

func main() {
	var inputJson = flag.String("input_json", "/Users/ali/src/misc/dash_answering/notebooks/data/99p/dataset.json", "where to load the data")
	var outputJson = flag.String("output_json", "/Users/ali/src/misc/dash_answering/notebooks/data/99p/dataset_embedding.json", "where to write the data")
	flag.Parse()
	client := embeddings.NewOpenAIClient(os.Getenv("OPENAI_API_KEY"))
	lsh := lsh.NewLSH42(lshSize, embeddingDim)

	log.Println("input json file:", *inputJson)
	readFile, err := os.Open(*inputJson)
	if err != nil {
		log.Fatalln(err)
	}
	defer readFile.Close()
	writeFile, err := os.Create(*outputJson)
	if err != nil {
		log.Fatalln(err)
	}
	defer writeFile.Close()
	encoder := json.NewEncoder(writeFile)

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Buffer(make([]byte, maxLineSize), maxLineSize)

	batch := []*document.Document{}
	cost := 0.0
	for fileScanner.Scan() {
		data := fileScanner.Bytes()
		doc := &document.Document{}
		json.Unmarshal(data, doc)
		batch = append(batch, doc)
		if len(batch) == batchSize {
			cost += float64(hydraiteBatch(batch, client, lsh)) * costDollarPerTocken
			log.Printf("cost: %0.02f$", cost)
			writeDocsToJson(batch, encoder)
			batch = []*document.Document{}
		}
	}
	if fileScanner.Err() != nil {
		log.Fatalln(fileScanner.Err())
	}
	cost += float64(hydraiteBatch(batch, client, lsh)) * costDollarPerTocken
	log.Printf("cost: %0.02f$", cost)
	writeDocsToJson(batch, encoder)
}
