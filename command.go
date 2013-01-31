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
	tag := fmt.Sprintf("CMD:SEARCH:%s", word)
	fmt.Printf("%s:START\n", tag)
	var err error

	records, err := boomkat.Search(word)
	if err != nil {
		log.Fatal(err)
	}
	for i, record := range records {
		fmt.Printf("%s:RES:[%s] = {id = %s, artist = %s, title = %s, label = %s, genre = %s, url = %s}\n",
			tag, i, record.Id, record.Artist, record.Title, record.Label, record.Genre, record.Url())
	}
	fmt.Printf("%s:END\n", tag)
}

func downloadRecord(recordId string) {
	tag := fmt.Sprintf("CMD:DownloadRecord:%s", recordId)
	fmt.Printf("%s:START\n", tag)
	tracks, err := tracksFromRecordId(recordId)
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan string)
	sem := make(chan int, 2)
	for _, track := range tracks {
		go func(track *boomkat.Track) {
			sem <- 1
			downloadTrackTask(track)
			<-sem
			done <- track.Id()
		}(track)
	}
	for i := 0; i < len(tracks); i++ {
		trackId := <-done
		fmt.Printf("%s:%s:DONE\n", tag, trackId)
	}
	fmt.Printf("%s:END\n", tag)
}

func downloadTrack(recordId, trackId string) {
	tag := fmt.Sprintf("CMD:DownloadTrack:%s:%s", recordId, trackId)
	fmt.Printf("%s:START\n", tag)
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
	fmt.Printf("%s:END\n", tag)
}

func downloadTrackTask(track *boomkat.Track) {
	tag := fmt.Sprintf("TASK:DLTRACK:%s", track.Id())
	fmt.Printf("%s:START\n", tag)
	err := track.Download()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s:END\n", tag)
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
