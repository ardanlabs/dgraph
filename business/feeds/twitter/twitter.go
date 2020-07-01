// Package twitter provides support for extracting a rescurion of followers
// from a single account.
package twitter

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

// levels specifies how many layers of friends we will
// retrieve from twitter.
const levels = 1

// User represents information about a twitter user.
type User struct {
	ID           int    `json:"id"`
	ScreenName   string `json:"screen_name"`
	Name         string `json:"name"`
	Location     string `json:"location"`
	FriendsCount int    `json:"friends_count"`
	Friends      []User
}

// Twitter represents the set of API's to access twitter data.
type Twitter struct {
	client http.Client
	token  string
}

// New constructs a Twitter value for use.
func New(token string) *Twitter {
	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	return &Twitter{
		client: client,
		token:  fmt.Sprintf("bearer %s", token),
	}
}

// RetrieveUser returns information for the specifed screen name
// includes their friends.
func (t *Twitter) RetrieveUser(ctx context.Context, screenName string) (User, error) {
	const twitterURL = "https://api.twitter.com/1.1/users/show.json?screen_name=%s"

	url := fmt.Sprintf(twitterURL, screenName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return User{}, fmt.Errorf("twitter create request error: %w", err)
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", t.token)

	resp, err := t.client.Do(req)
	if err != nil {
		return User{}, fmt.Errorf("twitter request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return User{}, fmt.Errorf("twitter op error: status code: %s", resp.Status)
	}

	var u User
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return User{}, fmt.Errorf("twitter decoding error: %w", err)
	}

	return u, nil
}

// RetrieveUserByID returns information for the specifed screen name
// includes their friends.
func (t *Twitter) RetrieveUserByID(ctx context.Context, id int) (User, error) {
	const twitterURL = "https://api.twitter.com/1.1/users/show.json?user_id=%d"

	url := fmt.Sprintf(twitterURL, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return User{}, fmt.Errorf("twitter create request error: %w", err)
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", t.token)

	resp, err := t.client.Do(req)
	if err != nil {
		return User{}, fmt.Errorf("twitter request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return User{}, fmt.Errorf("twitter op error: status code: %s", resp.Status)
	}

	var u User
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return User{}, fmt.Errorf("twitter decoding error: %w", err)
	}

	return u, nil
}

// RetrieveFriends returns information for the specifed screen name
// includes their friends.
func (t *Twitter) RetrieveFriends(ctx context.Context, id int) ([]User, error) {
	const twitterURL = "https://api.twitter.com/1.1/friends/ids.json?user_id=%d"

	url := fmt.Sprintf(twitterURL, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("twitter create request error: %w", err)
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", t.token)

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("twitter request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("twitter op error: status code: %s", resp.Status)
	}

	var friends struct {
		IDS []int `json:"ids"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&friends); err != nil {
		return nil, fmt.Errorf("twitter decoding error: %w", err)
	}

	users := make([]User, len(friends.IDS))
	for i, id := range friends.IDS {
		u, err := t.RetrieveUserByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("twitter retrieve user: %w", err)
		}
		users[i] = u
		fmt.Println(u)
	}

	return users, nil
}
