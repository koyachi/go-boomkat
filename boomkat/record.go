package boomkat

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/koyachi/go-boomkat/goquerywrapper"
	"io/ioutil"
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
	CoverUrl  string
	tracks    []*Track
}

func NewRecordFromId(id string) (*Record, error) {
	var doc *goquery.Document
	var e error

	recordUrl := fmt.Sprintf("http://boomkat.com/downloads/%s", id)
	if doc, e = goquerywrapper.NewDocument(recordUrl); e != nil {
		return nil, e
	}

	var coverUrl string
	if val, ok := doc.Find("div#main-image img").Attr("src"); ok {
		coverUrl = val
	}

	record := &Record{
		Id:     id,
		Artist: doc.Find("h1.product-header-artist-value").Text(),
		Title:  doc.Find("h1.product-header-title").Text(),
		Label:  doc.Find("div#product-header-label").Text(),
		Genre:  GenresFromString(doc.Find("div#product-header-genre a").Text()),
		//Thumbnail
		Review:   doc.Find("div#product-description-text").Text(),
		PageUrl:  recordUrl, // ???
		CoverUrl: coverUrl,
	}

	return record, nil
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
	for _, track := range r.tracks {
		track.SetRecord(*r)
	}
	return r.tracks, nil
}

func (r *Record) DownloadSampleTracks() {
	for _, track := range r.tracks {
		track.Download()
	}
}

func (r *Record) moreRecords(cssQuery string) ([]*Record, error) {
	var doc *goquery.Document
	var e error
	if doc, e = goquerywrapper.NewDocument(r.PageUrl); e != nil {
		return nil, e
	}

	elmRecords := doc.Find(fmt.Sprintf("%s div.data div.meta", cssQuery))
	records := make([]*Record, elmRecords.Length())
	elmRecords.Each(func(i int, s *goquery.Selection) {
		artist := s.Find("p.artist").Text()
		title := s.Find("p.title").Text()
		label := s.Find("p.label").Text()
		var recordUrl, recordId string
		if val, ok := s.Find("p.artist a").Attr("href"); ok {
			recordUrl = val
		}
		if reId.MatchString(recordUrl) {
			recordId = reId.FindStringSubmatch(recordUrl)[1]
		}
		records[i] = &Record{
			Id:      recordId,
			Artist:  artist,
			Title:   title,
			Label:   label,
			PageUrl: recordUrl,
		}
	})

	elmRecords = doc.Find(fmt.Sprintf("%s div.image.medium a img", cssQuery))
	elmRecords.Each(func(i int, s *goquery.Selection) {
		var coverUrl string
		if val, ok := s.Attr("src"); ok {
			coverUrl = val
		}
		records[i].CoverUrl = coverUrl
	})
	return records, nil
}

func (r *Record) RecordsAlsoBought() ([]*Record, error) {
	return r.moreRecords("div#slider-group-cross-sell")
}

func (r *Record) RecordsByTheSameArtist() ([]*Record, error) {
	return r.moreRecords("div#slider-group-same-artist")
}

func (r *Record) RecordsByTheSameLabel() ([]*Record, error) {
	return r.moreRecords("div#slider-group-same-label")
}

func (r *Record) RecordsYouMightLike() ([]*Record, error) {
	return r.moreRecords("div#slider-group-same-genre")
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
