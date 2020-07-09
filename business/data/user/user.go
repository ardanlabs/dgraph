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
mutation {
	addUser(input: [{
		source_id: "123456"
    	source: "twitter"
		screen_name: "goinggodotnet"
		name: "Bill"
		location: "miami, fl"
		friends_count: 10
	}])
	{
    	user {
	  		id
		}
  	}
}

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
