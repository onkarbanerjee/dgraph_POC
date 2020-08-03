package main

import (
	"context"

	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
)

func serviceToAlarmCorrelationNormalize(ctx context.Context, dgClient *dgo.Dgraph, fre string) (*api.Response, error) {
	q := `query all($a: string) {
		all(func: eq(name, $a))@recurse@normalize {
		  uid
		  nam : name
		  child
		}
	  }`

	txn := dgClient.NewReadOnlyTxn().BestEffort()
	defer txn.Commit(ctx)
	return txn.QueryWithVars(ctx, q, map[string]string{"$a": fre})
}

func serviceToAlarmCorrelation(ctx context.Context, dgClient *dgo.Dgraph, fre string) (*api.Response, error) {
	q := `query all($a: string) {
		all(func: eq(name, $a))@recurse(loop:false) {
		  uid
		  nam : name
		  child
		}
	  }`

	txn := dgClient.NewReadOnlyTxn().BestEffort()
	defer txn.Commit(ctx)
	return txn.QueryWithVars(ctx, q, map[string]string{"$a": fre})
}

func createSchema(ctx context.Context, dgClient *dgo.Dgraph, schema string) error {
	return dgClient.Alter(ctx, &api.Operation{
		Schema: schema,
	})
}
