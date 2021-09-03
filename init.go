package main

func init() {
	config, _ = ImportConfigFromFile("config.conf")
	connectMongo()
}
