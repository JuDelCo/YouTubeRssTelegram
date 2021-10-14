package lib

import (
	"context"
	"net/http"

	"github.com/mmcdole/gofeed"
)

// Private global variables (for lib package)
var cachedHttpClient *http.Client
var cachedRssParser *gofeed.Parser

func httpClient() *http.Client {
	if cachedHttpClient != nil {
		return cachedHttpClient
	}
	cachedHttpClient = &http.Client{}
	return cachedHttpClient
}

func rssParser() *gofeed.Parser {
	if cachedRssParser != nil {
		return cachedRssParser
	}
	cachedRssParser = gofeed.NewParser()
	return cachedRssParser
}

func ParseURL(feedURL string) (feed *gofeed.Feed, response *http.Response, err error) {
	client := httpClient()
	parser := rssParser()

	req, err := http.NewRequest("GET", feedURL, nil)
	if err != nil {
		return nil, nil, err
	}
	req = req.WithContext(context.Background())
	req.Header.Set("User-Agent", parser.UserAgent)
	resp, err := client.Do(req)

	if err != nil {
		return nil, nil, err
	}

	if resp != nil {
		defer func() {
			ce := resp.Body.Close()
			if ce != nil {
				err = ce
			}
		}()
	}

	feed, err = parser.Parse(resp.Body)
	response = resp

	if err != nil {
		return nil, nil, err
	}

	// TODO: Implement ETags or Last-Modified and return 304
	// https://pythonhosted.org/feedparser/http-etag.html
	// dt := time.Unix(rssSource.lastUpdate, 0).UTC() // compare "Last-Modified" header with this

	return feed, response, err
}
