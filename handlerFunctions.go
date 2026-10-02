package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/marcel-alter/AggreGator/internal/database"
	"github.com/marcel-alter/AggreGator/internal/rss"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("something went wrong in middlewareLoggedIn: %v", err)
		}
		return handler(s, cmd, user)
	}
}

// handler Functions ↓↓↓↓↓↓↓↓↓↓↓
func handlerLogin(s *state, cmd command) error {
	//fmt.Printf("LOGIN TRIGGERED, len of args = %d\n", len(cmd.args))
	if len(cmd.args) == 0 {
		//fmt.Println("Login args 0 triggered")
		return fmt.Errorf("Error: Login requieres a username argument!")
	}
	ctx := context.Background()
	if _, err := s.db.GetUser(ctx, cmd.args[0]); err != nil {
		fmt.Printf("User %v doesn't exist! Error: %v \nShutting down now!\n", cmd.args[0], err)
		os.Exit(1)
	}
	/*for _, arg := range cmd.args {
		fmt.Println(arg)
	}*/
	if err := s.cfg.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("Error at handlerLogin: %w", err)
	}
	fmt.Printf("User %s was set.\n", cmd.args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		//fmt.Println("Login args 0 triggered")
		return fmt.Errorf("Error: Register requieres a name argument!")
	}
	ctx := context.Background()
	//question: what is a context.Context? what do we use it for?

	if getU, err := s.db.GetUser(ctx, cmd.args[0]); getU.Name == cmd.args[0] {
		fmt.Printf("User %v already exist's! Shutting down now!\n", getU.Name)
		os.Exit(1)
	} else if err != nil {
		fmt.Printf("No user %v known yet! Status code: %v\n", cmd.args[0], err)
	}
	user := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}

	returnUser, err := s.db.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("Error: Couldn't Create User: %v", err)
	}

	fmt.Printf("User %v was created! Starting Login for new User now!\n%v\n", returnUser.Name, returnUser)
	if err := handlerLogin(s, cmd); err != nil {
		fmt.Printf("something went wrong trying to Login\n")
		return err
	}
	return nil
}

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()
	if err := s.db.RemoveAllUsers(ctx); err != nil {
		return fmt.Errorf("something went wrong when Resetting! Error: %v\n", err)
	}
	fmt.Println("Reset users table was Successful!")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	ctx := context.Background()
	allUsers, err := s.db.GetAllUser(ctx)
	if err != nil {
		return fmt.Errorf("something went wrong wiht GetAllUsers! Error: %v", err)
	}
	if len(allUsers) == 0 {
		fmt.Println("No users added yet!")
		return nil
	}
	for _, oneUser := range allUsers {
		name := oneUser.Name
		if name == s.cfg.CurrentUserName {
			name += " (current)"
		}
		fmt.Println(name)
	}
	return nil
}
func handlerAggUrl(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Error: agg requieres a URL argument!")
	}
	aggregate, err := rss.FetchFeed(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("AggUrl failed with FetchFeed! Error: %v", err)
	}
	fmt.Printf("Name: %v | URL: %v\n", aggregate.Channel.Title, aggregate.Channel.Link)
	return nil
}

func handlerAgg(s *state, cmd command) error {

	if len(cmd.args) == 0 {
		return fmt.Errorf("Error: agg requieres a time argument! e.g 1h20m5s")
	}
	duration, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("something went wrong in agg handler parsing time string '%v' to time. Error: %v", cmd.args[0], err)
	}

	ticker := time.NewTicker(duration)
	for ; ; <-ticker.C {
		scrapeFeeds(s)

		fmt.Println("Scraped Feeds at ...", time.Now())
	}
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) <= 1 {
		//fmt.Println("Login args 0 triggered")
		return fmt.Errorf("Error: addfeed requieres a name and url argument!")
	}
	ctx := context.Background()
	//question: what is a context.Context? what do we use it for?

	if checkFeed, err := s.db.GetFeed(ctx, cmd.args[1]); checkFeed.Name == cmd.args[0] {
		fmt.Printf("Feed '%v' already exist's! Shutting down now!\n", checkFeed.Name)
		os.Exit(1)
	} else if err != nil {
		fmt.Printf("No Feed '%v' known yet! Status code: %v\n", cmd.args[1], err)
	}
	/*agg, err := rss.FetchFeed(ctx, cmd.args[0])
	if err != nil {
		return fmt.Errorf("something went wrong in handlerAddFeed using FetchFeed! Error: %v", err)
	}*/
	feed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}

	returnFeed, err := s.db.CreateFeed(ctx, feed)
	if err != nil {
		return fmt.Errorf("Error: Couldn't Create Feed: %v", err)
	}
	follow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    returnFeed.ID,
	}
	if _, err := s.db.CreateFeedFollow(ctx, follow); err != nil {
		return fmt.Errorf("Something went wrong following user's own feed: %v", err)
	}
	fmt.Printf("Feed '%v' at '%v' was created by '%v'!\n", returnFeed.Name, returnFeed.Url, s.cfg.CurrentUserName)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	ctx := context.Background()
	allFeeds, err := s.db.GetAllFeeds(ctx)
	if err != nil {
		return fmt.Errorf("something went wrong wiht GetAllFeeds! Error: %v", err)
	}
	if len(allFeeds) == 0 {
		fmt.Println("No Feeds added yet!")
		return nil
	}
	for _, oneFeed := range allFeeds {
		name := oneFeed.Name
		url := oneFeed.Url
		/*userName, err := s.db.GetUserFromId(ctx, oneFeed.UserID)
		if err != nil {
			return fmt.Errorf("Error trying to get name from Feed with name '%v'. error: %v", name, err)
		}*/
		//question: why does oneFeed.Owner print {users.name true}
		fmt.Printf("Feed: '%v' | URL: '%v' | from User: '%v'\n", name, url, oneFeed.Owner.String)
	}
	return nil
}

func handlerHelp(s *state, cmd command) error {
	fmt.Println(`Available Commands:
	login <user>    - let's you change to a different user
	register <user> | let's you register a new user
	users		- display's all registered users
	reset		| delet's all registered users
	addfeed <URL>- creates a feed for the current user
	feeds		| display's all feeds from all users
	agg <time>	- starts an aggregation loop that iterates every <time> units
	aggurl <URL>| aggregates the contents of a website
	help		| display's all command handlers
	follow <URL>- let's the current user follow other user's feed's
	following	| display's all follows current user is folling
	unfollow <URL> - let's the current user unfollow the named feed
	browse		| display's all posts of all feeds the current user follows`)
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		//fmt.Println("Login args 0 triggered")
		return fmt.Errorf("Error: follow requieres a url argument!")
	}
	ctx := context.Background()
	//getFF, err := s.db.GetFeedFollow(ctx, cmd.args[0])
	if getFF, err := s.db.GetFeedFollow(ctx, cmd.args[0]); getFF.Url != cmd.args[0] {
		fmt.Printf("Feed '%v' doesn't exist! Shutting down now!\n", getFF.Url)
		os.Exit(1)
	} else if getFF.UserName == s.cfg.CurrentUserName {
		fmt.Printf("User already follows Feed '%v' owned by user '%v'! Status Code: %v\n", getFF.FeedName, getFF.CreatorName, err)
		os.Exit(1)
	}
	feed, err := s.db.GetFeed(ctx, cmd.args[0])
	if err != nil {
		return fmt.Errorf("Error getting current user's feed data! %v", err)
	}
	follow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	returnFeedFollow, err := s.db.CreateFeedFollow(ctx, follow)
	if err != nil {
		return fmt.Errorf("Error: Couldn't Create FeedFollow: %v", err)
	}
	for _, item := range returnFeedFollow {
		fmt.Printf("Feed '%v' by '%v' is now followed by current user: '%v' !\n", item.FeedName, item.CreatorName, s.cfg.CurrentUserName)
	}
	return nil
}

func handlerFollowing(s *state, cmd command) error {
	ctx := context.Background()
	allFollows, err := s.db.GetFeedFollowForUser(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("something went wrong wiht GetFeedFollowForUser! Error: %v", err)
	}
	if len(allFollows) == 0 {
		fmt.Println("No Feeds followed yet!")
		return nil
	}
	for _, oneFollow := range allFollows {
		name := oneFollow.FeedName
		url := oneFollow.Url
		/*feed, err := s.db.GetFeed(ctx, url)
		if err != nil {
			return fmt.Errorf("something went wrong in handler following trying to fetch feed data: %v", err)
		}*/

		fmt.Printf("'%v follows Feed: %v | URL: %v | from User: %v\n", oneFollow.UserName, name, url, oneFollow.CreatorName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Unfollow requires <URL> argument!")
	}
	ctx := context.Background()
	feed, err := s.db.GetFeed(ctx, cmd.args[0])
	if err != nil {
		return fmt.Errorf("Error trying to GetFeed for unfollow: %v", err)
	}
	deletion := database.DeleteFollowForUserParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	if err := s.db.DeleteFollowForUser(ctx, deletion); err != nil {
		return fmt.Errorf("Error trying to delete follow by user '%v' to feed '%v' owned by '%v'! Error: %v", user.Name, feed.Name, feed.Owner, err)
	}
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	var num int
	if len(cmd.args) != 0 {
		val, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("'%v' is not a number\n", cmd.args[0])
		}
		num = val
	} else {
		num = 2
	}

	args := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(num),
	}
	posts, err := s.db.GetPostsForUser(context.Background(), args)
	if err != nil {
		return fmt.Errorf("something went wrong with GetPostsForUser! Error: %v", err)
	}
	if len(posts) == 0 {
		fmt.Println("no posts fetched yet!")
		return nil
	}
	for _, post := range posts {
		fmt.Printf("Post: %v | URL: %v | published at: %v\n", post.Title, post.Url, post.PublishedAt)
	}
	return nil
}
