package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("Error when creating request: %v\n", err)
	}
	req.Header.Set("User-Agent", "gator")
	resp, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("Error when making request: %v\n", err)
	}
	respData, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	var rssFeed RSSFeed
	if err := xml.Unmarshal(respData, &rssFeed); err != nil {
		return &RSSFeed{}, fmt.Errorf("Error unmarshalling XML response: %v\n", err)
	}
	//return &rssFeed, nil
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	for i, _ := range rssFeed.Channel.Item {
		rssFeed.Channel.Item[i].Description = html.UnescapeString(rssFeed.Channel.Item[i].Description)
		rssFeed.Channel.Item[i].Title = html.UnescapeString(rssFeed.Channel.Item[i].Title)
		//fmt.Println(it.Description)
		//fmt.Println(html.UnescapeString(it.Description))

		//it.Description = html.UnescapeString(it.Description)
		//it.Title = html.UnescapeString(it.Title)
		//fmt.Println(it.Description)

	}
	fmt.Println(&rssFeed)
	return &rssFeed, nil
}
