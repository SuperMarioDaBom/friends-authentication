package main

type Config struct {
	Mongo struct {
	}
	Cassandra struct{}
}

func init() {
	connectMongo()
}
