// Package user provides CRUD access to the database.
package user

import (
	"context"
	"fmt"

	"github.com/ardanlabs/graphql"
	"github.com/pkg/errors"
)

// Set of error variables for CRUD operations.
var (
	ErrNotExists = errors.New("user does not exist")
	ErrExists    = errors.New("user exists")
	ErrNotFound  = errors.New("user not found")
)

// Add adds a new user to the database. If the user already exists
// this function will fail but the found user is returned. If the user is
// being added, the user with the id from the database is returned.
func Add(ctx context.Context, gql *graphql.GraphQL, nu NewUser) (User, error) {
	u := User{
		SourceID:     nu.SourceID,
		Source:       nu.Source,
		ScreenName:   nu.ScreenName,
		Name:         nu.Name,
		Location:     nu.Location,
		FriendsCount: nu.FriendsCount,
		Friends:      nu.Friends,
	}

	u, err := add(ctx, gql, u)
	if err != nil {
		return User{}, errors.Wrap(err, "adding user to database")
	}

	return u, nil
}

// AddFriend adds a new user to the database if the user doesn't already exist.
// Then the user is added to the collection of friends for the specified user id.
func AddFriend(ctx context.Context, gql *graphql.GraphQL, userID string, nu NewUser) (User, error) {
	// Validate the user doesn't already exists by screen name.
	// Validate the user isn't already in the list of friends for userID.
	/*
			mutation {
			updateUser(input: {
				filter: {
				id: [%q]
				},
				set: {
					friends: [{
						id: %q
					}]
				}
			})
			%s
		}
	*/

	return User{}, nil
}

// One returns the specified user from the database by the city id.
func One(ctx context.Context, gql *graphql.GraphQL, userID string) (User, error) {
	query := fmt.Sprintf(`
query {
	getUser(id: %q) {
		id
		source_id
    	source
		screen_name
		name
		location
		friends_count
	}
}`, userID)

	var result struct {
		GetUser User `json:"getUser"`
	}
	if err := gql.Query(ctx, query, &result); err != nil {
		return User{}, errors.Wrap(err, "query failed")
	}

	if result.GetUser.ID == "" {
		return User{}, ErrNotFound
	}

	return result.GetUser, nil
}

// OneByScreenName returns the specified user from the database by screen name.
func OneByScreenName(ctx context.Context, gql *graphql.GraphQL, screenName string) (User, error) {
	query := fmt.Sprintf(`
query {
	queryUser(filter: { screen_name: { eq: %q } }) {
		id
		source_id
    	source
		screen_name
		name
		location
		friends_count
	}
}`, screenName)

	var result struct {
		QueryUser []User `json:"queryUser"`
	}
	if err := gql.Query(ctx, query, &result); err != nil {
		return User{}, errors.Wrap(err, "query failed")
	}

	if len(result.QueryUser) != 1 {
		return User{}, ErrNotFound
	}

	return result.QueryUser[0], nil
}

// =============================================================================

func add(ctx context.Context, gql *graphql.GraphQL, user User) (User, error) {
	mutation, result := prepareAdd(user)
	if err := gql.Query(ctx, mutation, &result); err != nil {
		return User{}, errors.Wrap(err, "failed to add user")
	}

	if len(result.AddUser.User) != 1 {
		return User{}, errors.New("user id not returned")
	}

	user.ID = result.AddUser.User[0].ID
	return user, nil
}

func prepareAdd(user User) (string, addResult) {
	var result addResult
	mutation := fmt.Sprintf(`
mutation {
	addUser(input: [{
		source_id: %q
    	source: %q
		screen_name: %q
		name: %q
		location: %q
		friends_count: %d
	}])
	%s
}`, user.SourceID, user.Source, user.ScreenName, user.Name,
		user.Location, user.FriendsCount, result.document())

	return mutation, result
}

/*
query {
	queryUser(filter: { screen_name: { eq: "goinggodotnet" } })
	{
			id
      source_id
      source
      screen_name
      name
      location
      friends_count
  }
}

query {
	getUser(id: "0x3")
	{
	  id
      source_id
      source
      screen_name
      name
      location
      friends_count
  }
}
*/
