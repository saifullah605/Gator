package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	config "github.com/saifullah605/Gator/internal/config"
	"github.com/saifullah605/Gator/internal/database"
)

type state struct {
	db     *database.Queries
	config *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	commandFunc, ok := c.cmds[cmd.name]
	if !ok {
		return fmt.Errorf("command does not exist")
	}

	return commandFunc(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("the commands needs a username")
	}

	_, err := s.db.GetUser(context.Background(), cmd.arguments[0])

	if err == sql.ErrNoRows {
		return fmt.Errorf("user does not exist")
	} else if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	err = s.config.SetUser(cmd.arguments[0])

	if err != nil {
		return fmt.Errorf("cannot set username: %v", err)
	}

	fmt.Println("User has been set:", cmd.arguments[0])

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("the command neeeds a username")
	}

	_, err := s.db.GetUser(context.Background(), cmd.arguments[0])
	if err == nil {
		return fmt.Errorf("cannot create user, name already used")
	} else if err != sql.ErrNoRows {
		return fmt.Errorf("error: %v", err)
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
	})

	if err != nil {
		return fmt.Errorf("cannot create user: %v", err)
	}

	fmt.Println("User", user.Name, "created successfully")

	return s.config.SetUser(user.Name)

}

func handlerReset(s *state, cmd command) error {
	if err := s.db.ResetUsers(context.Background()); err != nil {
		fmt.Println("Reset users was not successful")
		return err
	}

	fmt.Println("Reset users was successful")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())

	if err != nil {
		return fmt.Errorf("cannot get users, error: %v", err)
	}

	for _, user := range users {
		if user.Name == s.config.CurrUserName {
			fmt.Println("*", user.Name, "(current)")
		}

		fmt.Println("*", user.Name)
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"
	feed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return err
	}

	fmt.Println(feed)
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	switch len(cmd.arguments) {
	case 0:
		return fmt.Errorf("need a name and url for feed, the first argument is the name, the second is the url, use quotation marks to wrap the name and url")
	case 1:
		return fmt.Errorf("need url for feed")
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
		Url:       cmd.arguments[1],
		UserID:    user.ID,
	})

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return fmt.Errorf("feed already exist, follow feed using the follow command")
			}
		}
		return fmt.Errorf("cannot add feed, error: %v", err)
	}

	if _, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}); err != nil {
		return fmt.Errorf("feed added, but cannot follow feed, please use feed command, error: %v\n%v", err, feed)
	}
	fmt.Println("feed added and automatically following feed")
	fmt.Println(feed)

	return nil
}

func handlerFeeds(s *state, cmd command) error {

	feeds, err := s.db.GetFeeds(context.Background())

	if err != nil {
		return fmt.Errorf("cannot get feeds, error: %v", err)
	}

	if len(feeds) == 0 {
		fmt.Println("There are no active feeds")
		return nil
	}

	for i, feed := range feeds {
		fmt.Printf("%v: user: %v name: %v url: %v\n", i+1, feed.User, feed.Name, feed.Url)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("needs a URL link")
	}

	feedId, err := s.db.GetFeedId(context.Background(), cmd.arguments[0])

	if err == sql.ErrNoRows {
		return fmt.Errorf("feed does not exist")
	} else if err != nil {
		return fmt.Errorf("cannot link feed, error: %v", err)
	}

	followedData, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feedId,
	})

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return fmt.Errorf("feed already followed")
			}
		}
		return fmt.Errorf("cannot follow feed, error: %v", err)
	}

	fmt.Printf("following feed success for user: %v of feed %v\n", followedData.UserName, followedData.FeedName)

	return nil
}

func handlerFollowingList(s *state, cmd command, user database.User) error {

	followList, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)

	if err != nil {
		return fmt.Errorf("cannot get followers list, error: %v", err)
	}

	for _, feed := range followList {
		fmt.Printf("%v\n", feed.FeedName)
	}

	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.config.CurrUserName)
		if err != nil {
			return fmt.Errorf("cannot get user info")
		}

		return handler(s, cmd, user)

	}

}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("need a url argument")
	}

	if _, err := s.db.Unfollow(context.Background(), database.UnfollowParams{
		UserID: user.ID,
		Url:    cmd.arguments[0],
	}); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user does not follow that feed")
		}

		return fmt.Errorf("cannot unfollow that feed, error: %v", err)
	}

	return nil
}
