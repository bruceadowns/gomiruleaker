package main

import (
	"log"
	"os"

	"github.com/bruceadowns/gomiruleaker/lib"
)

func main() {
	in, err := lib.InitInput(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	// construct pipeline
	// producer/source -> pipeline -> consumer/sink
	// generator -> parse[] -> accumulator -> poster

	log.Print("Start generating sources for the pipeline")
	genChan := lib.Generate(in.Targets, in.DownloadDelayMs)

	log.Print("Start parser channel to accept sources")
	parseChan := lib.Parse(genChan, in.ParserCount)

	log.Print("Start accumulator channel for parsed emails")
	accumChan, err := lib.Accum(parseChan, in.AccumBatchSize, in.OutputDir)
	if err != nil {
		log.Fatalf("Error creating accumulator channel: %s", err)
	}

	log.Print("Start posting channel to send emails to miru-leaks")
	postChan := lib.Post(accumChan, in.MiruURL, in.PostErrorDelayMs, in.PostCompress)

	// pull from sink until exhausted
	count := 0
	for p := range postChan {
		count += p
	}
	log.Printf("Posted %d emails", count)

	log.Printf("Done")
}
