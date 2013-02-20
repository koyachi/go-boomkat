package boomkat

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/koyachi/go-boomkat/goquerywrapper"
	"net/url"
	"regexp"
)

var reId = regexp.MustCompile(`^.*?(\d+)`)
var boomkatDir string = ""

func BoomkatDir() string {
	if boomkatDir != "" {
		return boomkatDir
	}
	return "/tmp/boomkat"
}

func SetBoomkatDir(dir string) {
	boomkatDir = dir
}

func Search(word string) ([]*Record, error) {
	var doc *goquery.Document
	var e error

	// TODO: urlEncode "word"
	searchUrl := fmt.Sprintf("http://boomkat.com/search?q=%s", url.QueryEscape(word))
	if doc, e = goquerywrapper.NewDocument(searchUrl); e != nil {
		return nil, e
	}

	elmRecords := doc.Find(".line")
	records := make([]*Record, elmRecords.Length())
	elmRecords.Each(func(i int, s *goquery.Selection) {
		var artist, title, label, review string
		var recordUrl, coverUrl string

		if val, ok := s.Find("div.image a img").Attr("src"); ok {
			coverUrl = val
		}
		meta := s.Find(".meta")
		artist = meta.Find("h4").Text()
		if val, ok := meta.Find("h4 a").Attr("href"); ok {
			recordUrl = val
		}
		var recordId string
		if reId.MatchString(recordUrl) {
			recordId = reId.FindStringSubmatch(recordUrl)[1]
		}
		title = meta.Find("p:nth-of-type(1)").Text()
		label = meta.Find("p:nth-of-type(2)").Text()
		// TODO: format
		genres := GenresFromString(meta.Find("p:nth-of-type(4)").Text())
		review = s.Find("div.review").Text()

		records[i] = &Record{
			Id:     recordId,
			Artist: artist,
			Title:  title,
			Label:  label,
			Genre:  genres,
			// TODO: Thumbnail
			Review:   review,
			PageUrl:  recordUrl,
			CoverUrl: coverUrl,
		}
	})

	return records, nil
}
