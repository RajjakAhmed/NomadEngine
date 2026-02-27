package tools

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func init() {
	Register(Tool{
		Name: "fetch_url",
		Run:  FetchURL,
	})
}

func FetchURL(input map[string]interface{}) (map[string]interface{}, error) {

	url, ok := input["url"].(string)
	if !ok || url == "" {
		return nil, fmt.Errorf("url is required")
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// remove scripts and styles
	doc.Find("script").Remove()
	doc.Find("style").Remove()
	doc.Find("noscript").Remove()

	text := doc.Find("body").Text()

	// clean spaces
	text = strings.Join(strings.Fields(text), " ")

	// limit size so LLM doesn't explode
	if len(text) > 4000 {
		text = text[:4000]
	}

	return map[string]interface{}{
		"text": text,
	}, nil
}