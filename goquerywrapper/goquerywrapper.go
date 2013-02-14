package goquerywrapper

import (
	"exp/html"
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

func newHttpClient() *http.Client {
	transport := &http.Transport{
		Dial:            dialFunc(),
		TLSClientConfig: tlsClientConfig(),
	}
	return &http.Client{Transport: transport}
}

func NewDocument(url string) (d *goquery.Document, e error) {
	client := newHttpClient()
	res, e := client.Get(url)
	if e != nil {
		return
	}
	defer res.Body.Close()

	// Parse the HTML into nodes
	root, e := html.Parse(res.Body)
	if e != nil {
		return
	}

	// Create and fill the document
	d = goquery.NewDocumentFromNode(root)
	return
}
