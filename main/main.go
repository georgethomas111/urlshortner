package main

import "github.com/georgethomas111/urlshortner"

func main() {
	err := urlshortner.Start(":8080", "couchdb", 5984, "tinyurl", "objectviews")
	if err != nil {
		panic(err)
	}
}
