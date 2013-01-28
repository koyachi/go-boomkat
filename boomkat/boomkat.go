package boomkat

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"regexp"
	"strings"
)

var reId = regexp.MustCompile(`^.*?(\d+)`)

func Search(word string) ([]*Record, error) {
	var doc *goquery.Document
	var e error

	// TODO: urlEncode "word"
	searchUrl := fmt.Sprintf("http://boomkat.com/search?q=%s", url.QueryEscape(word))
	fmt.Printf("searchUrl = %s\n", searchUrl)
	if doc, e = goquery.NewDocument(searchUrl); e != nil {
		return nil, e
	}

	elmRecords := doc.Find(".line")
	records := make([]*Record, elmRecords.Length())
	elmRecords.Each(func(i int, s *goquery.Selection) {
		var artist, title, label, genre, review string
		var recordUrl string

		meta := s.Find(".meta")
		artist = meta.Find("h4").Text()
		if val, ok := meta.Find("h4 a").Attr("href"); ok {
			recordUrl = val
		}
		var recordId string
		if reId.MatchString(recordUrl) {
			recordId = reId.FindStringSubmatch(recordUrl)[1]
		}
		fmt.Printf("recordId = %s\n", recordId)
		title = meta.Find("p:nth-of-type(1)").Text()
		label = meta.Find("p:nth-of-type(2)").Text()
		// TODO: format
		genre = meta.Find("p:nth-of-type(4)").Text()
		genres := strings.Split(genre, "/")
		for i, genre := range genres {
			genres[i] = strings.TrimSpace(genre)
		}
		review = s.Find("div.review").Text()

		records[i] = &Record{
			Id:     recordId,
			Artist: artist,
			Title:  title,
			Label:  label,
			Genre:  genres,
			// TODO: Thumbnail
			Review:  review,
			PageUrl: recordUrl,
		}
	})

	return records, nil
}

func BoomkatDir() string {
	return "/tmp/boomkat"
}