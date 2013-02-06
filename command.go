package main

import (
	"fmt"
	"github.com/koyachi/go-boomkat/boomkat"
	"log"
)

func search(word string) {
	tag := fmt.Sprintf("CMD:SEARCH:%s", word)
	fmt.Printf("%s:START\n", tag)
	var err error

	records, err := boomkat.Search(word)
	if err != nil {
		log.Fatal(err)
	}
	dumpRecords(records, tag, "RES", false)
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
	tag := fmt.Sprintf("TSK:DLTRACK:%s", track.Id())
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

func dumpRecord(index int, commandTag, typeTag string, record *boomkat.Record) {
	fmt.Printf("%s:%s:[%d] = {id = %s, artist = %s, title = %s, label = %s, genre = %s, url = %s}\n",
		commandTag, typeTag, index, record.Id, record.Artist, record.Title, record.Label, record.Genre, record.Url())
}

func dumpRecords(records []*boomkat.Record, commandTag, typeTag string, doAsync bool) {
	for i, record := range records {
		if doAsync {
			go dumpRecord(i, commandTag, typeTag, record)
		} else {
			dumpRecord(i, commandTag, typeTag, record)
		}
	}
}

func recordInfo(recordId string) {
	tag := fmt.Sprintf("CMD:RecordInfo:%s", recordId)
	fmt.Printf("%s:START\n", tag)
	record, err := boomkat.NewRecordFromId(recordId)
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan bool)
	go func(record *boomkat.Record) {
		recordsAlsoBought, err := record.RecordsAlsoBought()
		if err != nil {
			log.Fatal(err)
		}
		dumpRecords(recordsAlsoBought, tag, "ALSO_BOUGHT", true)
		done <- true
	}(record)
	go func(record *boomkat.Record) {
		recordsByTheSameArtist, err := record.RecordsByTheSameArtist()
		if err != nil {
			log.Fatal(err)
		}
		dumpRecords(recordsByTheSameArtist, tag, "BY_THE_SAME_ARTIST", true)
		done <- true
	}(record)
	go func(record *boomkat.Record) {
		recordsByTheSameLabel, err := record.RecordsByTheSameLabel()
		if err != nil {
			log.Fatal(err)
		}
		dumpRecords(recordsByTheSameLabel, tag, "BY_THE_SAME_LABEL", true)
		done <- true
	}(record)
	go func(record *boomkat.Record) {
		recordsYouMightLike, err := record.RecordsYouMightLike()
		if err != nil {
			log.Fatal(err)
		}
		dumpRecords(recordsYouMightLike, tag, "YOU_MIGHT_LIKE", true)
		done <- true
	}(record)
	for i := 0; i < 4; i++ {
		<-done
	}
	fmt.Printf("%s:END\n", tag)
}
