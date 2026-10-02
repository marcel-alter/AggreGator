/*
package main

import (

	"context"
	"fmt"
	"time"

	"github.com/marcel-alter/AggreGator/internal/database"
	"github.com/marcel-alter/AggreGator/internal/rss"

)

	func ScrapeFeeds(s *state) error {
		ctx := context.Background()
		next, err := s.db.GetNextFeedToFetch(ctx)
		if err != nil {
			return fmt.Errorf("something went wrong in aggregation fetching Next feed! Error: %v", err)
		}
		args := database.MarkFeedFetchedParams{
			UpdatedAt: time.Now(),
			ID:        next.ID,
		}
		if err := s.db.MarkFeedFetched(ctx, args); err != nil {
			return fmt.Errorf("something went wrong in aggregation Marking feed as fetched! Error: %v", err)
		}
		feed, err := rss.FetchFeed(ctx, next.Url)
		if err != nil {
			return fmt.Errorf("something went wrong in aggregation fetching feed again! Error: %v", err)
		}
		fmt.Printf("Success fetching <%v>! Now Printing Items titles...\n", feed.Channel.Title)
		for _, item := range feed.Channel.Item {
			fmt.Printf("Title:%v | URL: %v | Channel: %v\n", item.Title, item.Link, feed.Channel.Title)
		}
		return nil
	}
*/
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"

	"github.com/google/uuid"
	"github.com/marcel-alter/AggreGator/internal/database"
	"github.com/marcel-alter/AggreGator/internal/rss"
)

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("Couldn't get next feeds to fetch", err)
		return
	}
	log.Println("Found a feed to fetch!")
	scrapeFeed(s.db, feed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	args := database.MarkFeedFetchedParams{
		UpdatedAt: time.Now(),
		ID:        feed.ID,
	}
	err := db.MarkFeedFetched(context.Background(), args)
	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
		return
	}

	feedData, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}
	for _, item := range feedData.Channel.Item {
		var pub sql.NullTime
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			pub = sql.NullTime{Time: t, Valid: true}
		}
		args := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
			PublishedAt: pub,
			FeedID:      feed.ID,
		}
		post, err := db.CreatePost(context.Background(), args)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				// 23505 = unique_violation in Postgres
				continue // skip duplicate, no need to log
			}
			log.Printf("error creating post: %v", err)
			continue
		}
		fmt.Printf("Found post: %s\n", post.Title)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
}
