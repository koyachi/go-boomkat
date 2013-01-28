package boomkat

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Record struct {
	Id        string
	Artist    string
	Title     string
	Label     string
	Genre     []string
	Thumbnail string
	Review    string
	PageUrl   string
	tracks    []*Track
}

func (r *Record) Url() string {
	return fmt.Sprintf("http://boomkat.com%s", r.PageUrl)
}

type SampleTrackResponse struct {
	Id     int64    `json:"id"`
	Tracks []*Track `json:"boomboxx_sample_tracks"`
}

func (r *Record) SampleTracks() ([]*Track, error) {
	if len(r.tracks) > 0 {
		return r.tracks, nil
	}

	sampleTracksUrl := fmt.Sprintf("http://boomkat.com/boomboxx_sample_tracks_by_album?id=%s&product_id=%s", r.Id, r.Id)
	res, err := http.Get(sampleTracksUrl)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	var sampleTrackResponse SampleTrackResponse
	err = json.Unmarshal([]byte(body), &sampleTrackResponse)
	if err != nil {
		return nil, err
	}

	r.tracks = sampleTrackResponse.Tracks
	for i, track := range r.tracks {
		track.SetRecord(*r)
		log.Printf("  [%d] track.SetRecord()", i)
	}
	return r.tracks, nil
}

func (r *Record) DownloadSampleTracks() {
	for _, track := range r.tracks {
		track.Download()
	}
}

func (r *Record) WorkDir() (string, error) {
	var dirName = filepath.Join(BoomkatDir(), r.Id)
	var f *os.File
	var err error
	f, err = os.OpenFile(dirName, os.O_RDONLY, 0600)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dirName, 0766)
			if err != nil && !os.IsNotExist(err) {
				return "", err
			}
		} else {
			return "", err
		}
	} else {
		defer f.Close()
	}
	return dirName, nil
}