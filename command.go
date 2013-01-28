package main

import (
	"fmt"
	"github.com/koyachi/go-boomkat/boomkat"
	"log"
)

func _searchTest(word string) {
	var err error

	//records, err := boomkat.Search("Tim Hecker")
	records, err := boomkat.Search(word)
	if err != nil {
		log.Fatal(err)
	}
	for _, record := range records {
		fmt.Printf("[%s] artist = %s, title = %s, label = %s, genre = %s, url = %s\n",
			record.Id, record.Artist, record.Title, record.Label, record.Genre, record.Url())

		if record.Id == "598941" {
			sampleTracks, err := record.SampleTracks()
			if err != nil {
				log.Fatal(err)
			}
			for i, track := range sampleTracks {
				if i == 2 {
					err = track.Download()
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	}
}

func search(word string) {
	var err error

	records, err := boomkat.Search(word)
	if err != nil {
		log.Fatal(err)
	}
	for _, record := range records {
		fmt.Printf("[%s] artist = %s, title = %s, label = %s, genre = %s, url = %s\n",
			record.Id, record.Artist, record.Title, record.Label, record.Genre, record.Url())
	}
}

func downloadRecord() {
}

func downloadTrack() {
}
