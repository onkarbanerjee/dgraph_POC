package main

import (
	"context"
	"fmt"

	"math/rand"
	"os"
	"time"

	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc"
)

const (
	// 2GB
	MAX_MSG_SIZE = 2 * 1024 * 1024 * 1024
	// 1K
	BATCH_SIZE = 1000
)

var (
	n,
	p int
)

type CancelFunc func()

func main() {
	fmt.Println("Start")

	file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	log.SetOutput(file)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)

	log.Info("Application start at", time.Now())

	parseFlags()
	ctx := context.TODO()

	dgClient, cancel := getDgraphClient(cfg.endPoint)
	defer cancel()

	// defer txn.Commit(ctx)

	schema := `<child>: [uid] @reverse .
	<name>: string @index(exact) .

	type FRE {
		name: string
		child: [uid]
		<~child>: [uid]
	}

	type TPE {
		name: string
		child: [uid]
		<~child>: [uid]
	}

	type Alarm {
		name: string
		child: [uid]
		<~child>: [uid]
	}
   `

	if err := createSchema(ctx, dgClient, schema); err != nil {
		log.Fatal("failed to create schema")
	}

	// 10k
	n = 10000
	p = 10
	// 10 batches
	maxNumberOfBatches := n / BATCH_SIZE
	if n%BATCH_SIZE > 0 {
		maxNumberOfBatches++
	}

	for batchNo := 0; batchNo < maxNumberOfBatches; batchNo++ {

		txn := dgClient.NewTxn()
		start := batchNo*BATCH_SIZE + 1
		end := (batchNo + 1) * BATCH_SIZE
		log.Infof("Batch %d start here, start idx %d, end %d, time now %v\n", batchNo, start, end, time.Now())
		if _, err := createNode(ctx, txn, "TPE", start, end); err != nil {
			log.Fatal("failed to mutate ", err)
		}
		if err := txn.Commit(ctx); err != nil {
			log.Fatal("failed to commit transaction", err)
		}
		<-time.After(3 * time.Second)

	}

	log.Info("FREs now!!")
	m := 10
	selfName := "FRE"
	childName := "TPE"
	// m = 3
	for level := 2; level <= m; level++ {

		if level == 3 {
			childName = "FRE"
		}
		// 10 batches , each of 1000, at each level
		for batchNo := 1; batchNo <= maxNumberOfBatches; batchNo++ {
			createFREs(ctx, dgClient, batchNo, level, selfName, childName)

		}
		log.Info("Finishing level ", level)
		<-time.After(7 * time.Second)

	}
	log.Info("Program exiting now")
	fmt.Println("Bye now!!")

}

// create the FREs at the given level and for the given batchno.
func createFREs(ctx context.Context, dgClient *dgo.Dgraph, batchNo int, level int, selfName, childName string) {
	txn := dgClient.NewTxn()

	batchStart, batchEnd := getBatchStartAndEnd(level, batchNo)

	log.Infof("Level %d Batch %d start here, start idx %d, end %d, time now %v\n", level, batchNo, batchStart, batchEnd, time.Now())
	if _, err := createNode1(ctx, txn, selfName, childName, level, batchStart, batchEnd); err != nil {
		log.Fatal("failed to mutate ", err)
	}

	if err := txn.Commit(ctx); err != nil {
		log.Fatal("failed to commit transaction")
	}
	<-time.After(3 * time.Second)
}

func createNode(ctx context.Context, txn *dgo.Txn, node_name string, start, end int) (*api.Response, error) {
	// query := "query {\n"
	selfQuad := "_:self%d <name> \"%s%d\" .\n_:self%d <dgraph.type> \"TPE\" .\n"
	nQuads := []byte{}
	for i := start; i <= end; i++ {
		m := []byte(fmt.Sprintf(selfQuad, i, node_name, i, i))
		nQuads = append(nQuads, m...)
	}

	mu := &api.Mutation{
		SetNquads: nQuads,
	}

	return txn.Mutate(ctx, mu)
	// req := &api.Request{
	// 	Mutations: []*api.Mutation{mu},
	// }
	// if len(node.Children) > 0 {
	// 	req.Query = query
	// }
	// return txn.Do(ctx, req)

}

func createNode1(ctx context.Context, txn *dgo.Txn, selfName, childName string, level, batchStart, batchEnd int) (*api.Response, error) {

	query := "query {\n"
	queryQuad := "child%d%d as var(func: eq(name,\"%s%d\"))\n"
	selfQuad := "_:self%d <name> \"%s%d\" .\n_:self%d <dgraph.type> \"FRE\" .\n"
	childQuad := "_:self%d <child> uid(child%d%d) .\n"
	nQuads := []byte{}
	var q []byte
	q = []byte(query)

	for nodeNumber := batchStart; nodeNumber <= batchEnd; nodeNumber++ {
		m := []byte(fmt.Sprintf(selfQuad, nodeNumber, selfName, nodeNumber, nodeNumber))
		maxChildNumber := (level - 1) * n
		minChildNumber := maxChildNumber - n + 1
		// endRange := n + 10000
		childrenNumbers := getChildrenNumbers(int64(nodeNumber*2), minChildNumber, maxChildNumber, p)
		for idx, childrenNumber := range childrenNumbers {
			q1 := []byte(fmt.Sprintf(queryQuad, nodeNumber, idx, childName, childrenNumber))
			m1 := []byte(fmt.Sprintf(childQuad, nodeNumber, nodeNumber, idx))
			m = append(m, m1...)
			q = append(q, q1...)
		}
		nQuads = append(nQuads, m...)

	}
	q = append(q, []byte("}")...)
	mu := &api.Mutation{
		SetNquads: nQuads,
	}

	log.Info("q is ", string(q))
	log.Info("mu is ", string(nQuads))

	req := &api.Request{
		Mutations: []*api.Mutation{mu},
	}
	if len(q) > 0 {
		req.Query = string(q)
	}

	return txn.Do(ctx, req)

}

func getDgraphClient(endpoint string) (*dgo.Dgraph, CancelFunc) {
	conn, err := grpc.Dial(endpoint,
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MAX_MSG_SIZE)))

	if err != nil {
		log.Fatal("While trying to dial gRPC")
	}

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)
	// ctx := context.Background()

	// TODO: Skip now, but will need to be implemented with ACL and enterprise features enabled
	// Perform login call. If the Dgraph cluster does not have ACL and
	// enterprise features enabled, this call should be skipped.
	// for {
	// 	// Keep retrying until we succeed or receive a non-retriable error.
	// 	err = dg.Login(ctx, "groot", "password")
	// 	if err == nil || !strings.Contains(err.Error(), "Please retry") {
	// 		break
	// 	}
	// 	time.Sleep(time.Second)
	// }
	// if err != nil {
	// 	log.Fatalf("While trying to login %v", err.Error())
	// }

	return dg, func() {
		if err := conn.Close(); err != nil {
			log.Errorf("Error while closing connection:%v", err)
		}
	}
}

func createTPEs(ctx context.Context, txn *dgo.Txn, n int) {
	// for i := 1; i <= n; i++ {
	// 	node := New(fmt.Sprintf("tpe%d", i), []string{})
	// 	res, err := createNode(ctx, txn, node)
	// 	if err != nil {
	// 		log.Fatal("failed to mutate ", err)
	// 	}
	// 	fmt.Println("res is", res.String())
	// }
}

func getChildrenNumbers(seed int64, min, max, numberOfRands int) []int {
	childrenNumbers := make([]int, numberOfRands)
	rand.Seed(seed)
	for i := 0; i < numberOfRands; i++ {
		r := rand.Intn(max-min+1) + min
		childrenNumbers[i] = r
	}
	return childrenNumbers
}

func getBatchStartAndEnd(level, batchNo int) (int, int) {
	start := (level-1)*n + ((batchNo - 1) * BATCH_SIZE)
	end := start + BATCH_SIZE - 1
	return start, end
}
