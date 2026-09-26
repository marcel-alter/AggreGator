package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/marcel-alter/AggreGator/internal/database"
	"github.com/marcel-alter/AggreGator/internal/rss"
)

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

func handlerAgg(s *state, cmd command) error {
	ctx := context.Background()
	rssFeed, err := rss.FetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("Error using FetchFeed: %v\n", err)
	}
	fmt.Println(rssFeed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) <= 1 {
		//fmt.Println("Login args 0 triggered")
		return fmt.Errorf("Error: addfeed requieres a name and url argument!")
	}
	ctx := context.Background()
	//question: what is a context.Context? what do we use it for?

	if getU, err := s.db.GetFeed(ctx, cmd.args[0]); getU.Name == cmd.args[0] {
		fmt.Printf("Feed %v already exist's! Shutting down now!\n", getU.Name)
		os.Exit(1)
	} else if err != nil {
		fmt.Printf("No Feed %v known yet! Status code: %v\n", cmd.args[0], err)
	}
	user, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
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

	fmt.Printf("Feed %v was created!\n", returnFeed.Name)
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
		userName, err := s.db.GetUserFromId(ctx, oneFeed.UserID)
		if err != nil {
			return fmt.Errorf("Error trying to get name from Feed with name '%v'. error: %v", name, err)
		}
		fmt.Printf("Feed: %v | URL: %v | from User: %v\n", name, url, userName.Name)
	}
	return nil
}
