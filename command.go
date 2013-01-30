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

func downloadRecord(recordId string) {
	fmt.Printf("downloat record[record_id:%s]\n", recordId)
	tracks, err := tracksFromRecordId(recordId)
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan bool)
	sem := make(chan int, 2)
	for _, track := range tracks {
		go func(track *boomkat.Track) {
			sem <- 1
			downloadTrackTask(track)
			<-sem
			done <- true
		}(track)
	}
	for i := 0; i < len(tracks); i++ {
		<-done
		fmt.Printf("[%d:done]\n", i)
	}
}

func downloadTrack(recordId, trackId string) {
	fmt.Printf("downloat track[record_id:%s, track_id:%s]\n", recordId, trackId)
	tracks, err := tracksFromRecordId(recordId)
	if err != nil {
		log.Fatal(err)
	}
	for _, track := range tracks {
		if track.Id() == trackId {
			downloadTrackTask(track)
			break
		}
	}
}

func downloadTrackTask(track *boomkat.Track) {
	fmt.Printf("track_id:%s start...", track.Id())
	err := track.Download()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("done.\n")
}

func tracksFromRecordId(recordId string) ([]*boomkat.Track, error) {
	record, err := boomkat.NewRecordFromId(recordId)
	if err != nil {
		return nil, err
	}
	tracks, err := record.SampleTracks()
	if err != nil {
		return nil, err
	}
	return tracks, nil
}
