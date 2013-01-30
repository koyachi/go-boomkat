package boomkat

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
)

var reFilename = regexp.MustCompile(`^.*?(\d+.mp3)`)
var reFileId = regexp.MustCompile(`^.*?(\d+).mp3`)

type Track struct {
	ProductId int64  `json:"product_id"`
	Title     string `json:"title"`
	Url       string `json:"url"`
	Image     string `json:"image"`
	SampleMp3 string `json:"sample_mp3"`
	SampleOgg string `json:"sample_ogg"`
	record    Record
}

func (t *Track) Record() Record {
	return t.record
}

func (t *Track) SetRecord(r Record) {
	t.record = r
}

func (t *Track) filename() string {
	matchs := reFilename.FindStringSubmatch(t.SampleMp3)
	if len(matchs) > 1 {
		return matchs[1]
	}

	// TODO: 適当にユニーク文字列生成 prefix + time
	return ""
}

func (t *Track) Id() string {
	matchs := reFileId.FindStringSubmatch(t.SampleMp3)
	if len(matchs) > 1 {
		return matchs[1]
	}

	// TODO:
	return ""
}

func (t *Track) Download() error {
	res, err := http.Get(t.SampleMp3)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return err
	}

	dir, err := t.record.WorkDir()
	if err != nil {
		return err
	}

	trackPath := filepath.Join(dir, t.filename())
	err = ioutil.WriteFile(trackPath, body, 0666)
	if err != nil {
		return err
	}
	return nil
}
