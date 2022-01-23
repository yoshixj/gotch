package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"net/url"
	"os"
	"strings"
	"time"
)

const METADATA_DIR = ".metadata"

type MetaData struct {
	Site      string     `json:"site"`
	NumLinks  int        `json:"num_links"`
	Images    int        `json:"images"`
	LastFetch *time.Time `json:"last_fetch"`
}

func NewMetaData(site string) (*MetaData, error) {
	return &MetaData{Site: site, NumLinks: 0, Images: 0, LastFetch: nil}, nil
}

func (m *MetaData) SetMetaData(htmlString string) {
	m.countNode(htmlString)
}

func (m *MetaData) countNode(htmlString string) {
	r := strings.NewReader(htmlString)
	tokenizer := html.NewTokenizer(r)
	end := false

	for {
		tt := tokenizer.Next()
		switch {
		case tt == html.StartTagToken:
			token := tokenizer.Token()
			fmt.Println(token.Data)
			switch token.Data {
			case "a":

				for _, attr := range token.Attr {

					if attr.Key == "href" {
						link := attr.Val

						_, parseLinkErr := url.Parse(link)
						if parseLinkErr == nil {
							m.NumLinks++
						}

						if parseLinkErr != nil {
							fmt.Println("Can't parse: " + token.Data)
						}
					}
				}
				break
			case "img":
				fmt.Printf("%v", m)
				m.Images++
			}
		case tt == html.ErrorToken:
			end = true
			break
		}
		if end {
			break
		}
	}
}

func (m *MetaData) Store() {
	data, _ := json.Marshal(m)
	os.Mkdir(METADATA_DIR, 0777)
	filePath := fmt.Sprintf("%s/%s", METADATA_DIR, m.Site)
	writeFile(string(data), filePath)
}

func (m *MetaData) Read() {
	jsonFromFile, err := os.ReadFile(fmt.Sprintf("%s/%s", METADATA_DIR, m.Site))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	json.Unmarshal(jsonFromFile, m)

	fmt.Printf("site: %s \n", m.Site)
	fmt.Printf("num_links: %d \n", m.NumLinks)
	fmt.Printf("images: %d \n", m.Images)
	fmt.Printf("last_fetch: %s \n", m.LastFetch)
}
