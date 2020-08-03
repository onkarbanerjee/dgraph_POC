package main

import "flag"

var cfg config

type config struct {
	m,
	n,
	p,
	q int
	endPoint string
}

func init() {
	cfg = config{}
}

func parseFlags() {
	flag.StringVar(&cfg.endPoint, "endpoint", "localhost:9080", "Dgraph endpint URL")
	flag.IntVar(&cfg.m, "m", 10, "Depth of the dataset")
	flag.IntVar(&cfg.n, "n", 10000, "Breadth of the dataset")
	flag.IntVar(&cfg.p, "p", 10, "Number of children of each FRE")
	flag.IntVar(&cfg.q, "q", 5, "Number of alarms on each TPE")
	flag.Parse()
	return
}
