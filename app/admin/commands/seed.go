package commands

import (
	"context"
	"log"
	"time"

	"github.com/ardanlabs/dgraph/business/feeds/twitter"
)

// Seed will seed the database for a given user.
func Seed(log *log.Logger, token string, screenName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t := twitter.New(log, token)
	u, err := t.RetrieveUser(ctx, screenName)
	if err != nil {
		return err
	}

	_, err = t.RetrieveFriends(ctx, u.ID)
	if err != nil {
		return err
	}

	return nil
}
