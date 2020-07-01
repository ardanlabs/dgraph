package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/ardanlabs/dgraph/business/feeds/twitter"
)

// Seed will seed the database for a given user.
func Seed(token string, screenName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t := twitter.New(token)
	u, err := t.RetrieveUser(ctx, screenName)
	if err != nil {
		return err
	}
	fmt.Println(u)

	users, err := t.RetrieveFriends(ctx, u.ID)
	if err != nil {
		return err
	}
	fmt.Println(users)

	return nil
}
