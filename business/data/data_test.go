package data_test

import (
	"context"
	"testing"
	"time"

	"github.com/ardanlabs/dgraph/business/data"
	"github.com/ardanlabs/dgraph/business/data/ready"
	"github.com/ardanlabs/dgraph/business/data/schema"
	"github.com/ardanlabs/dgraph/foundation/tests"
	"github.com/ardanlabs/graphql"
)

// TestData validates all the mutation support in data.
func TestData(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	url, teardown := tests.NewUnit(t)
	t.Cleanup(teardown)

	t.Run("readiness", readiness(url))
}

// waitReady provides support for making sure the database is ready to be used.
func waitReady(t *testing.T, ctx context.Context, testID int, url string) (*schema.Schema, *graphql.GraphQL) {
	err := ready.Validate(ctx, url, time.Second)
	if err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to see Dgraph is ready: %v", tests.Failed, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to to see Dgraph is ready.", tests.Success, testID)

	gqlConfig := data.GraphQLConfig{
		URL: url,
	}
	gql := data.NewGraphQL(gqlConfig)

	schema := schema.New(gql)
	t.Logf("\t%s\tTest %d:\tShould be able to prepare the schema.", tests.Success, testID)

	return schema, gql
}

// readiness validates the health check is working.
func readiness(url string) func(t *testing.T) {
	tf := func(t *testing.T) {
		type tableTest struct {
			name       string
			retryDelay time.Duration
			timeout    time.Duration
			success    bool
		}

		tt := []tableTest{
			{"timeout", 500 * time.Millisecond, time.Second, false},
			{"ready", 500 * time.Millisecond, 20 * time.Second, true},
		}

		t.Log("Given the need to be able to validate the database is ready.")
		{
			for testID, test := range tt {
				tf := func(t *testing.T) {
					t.Logf("\tTest %d:\tWhen waiting up to %v for the database to be ready.", testID, test.timeout)
					{
						ctx, cancel := context.WithTimeout(context.Background(), test.timeout)
						defer cancel()

						err := ready.Validate(ctx, url, test.retryDelay)
						switch test.success {
						case true:
							if err != nil {
								t.Fatalf("\t%s\tTest %d:\tShould be able to see Dgraph is ready: %v", tests.Failed, testID, err)
							}
							t.Logf("\t%s\tTest %d:\tShould be able to see Dgraph is ready.", tests.Success, testID)

						case false:
							if err == nil {
								t.Fatalf("\t%s\tTest %d:\tShould be able to see Dgraph is Not ready.", tests.Failed, testID)
							}
							t.Logf("\t%s\tTest %d:\tShould be able to see Dgraph is Not ready.", tests.Success, testID)
						}
					}
				}
				t.Run(test.name, tf)
			}
		}
	}
	return tf
}
