package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/ardanlabs/dgraph/business/data"
	"github.com/ardanlabs/dgraph/business/data/schema"
)

// Schema handles the updating of the schema.
func Schema(gqlConfig data.GraphQLConfig) error {
	schema := schema.New(gqlConfig)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := schema.Create(ctx); err != nil {
		return err
	}

	fmt.Println("schema updated")
	return nil
}
