package main

import (
	"log"
	"net/http"
	"os"

	"github.com/voutasaurus/cbtest/api"
	"github.com/voutasaurus/cbtest/database"
	"github.com/voutasaurus/cbtest/env"
)

func main() {
	logger := log.New(os.Stderr, "hello: ", log.Llongfile|log.Lmicroseconds|log.LstdFlags)
	logger.Println("starting...")

	var (
		addr       = env.Get("HELLO_ADDR").WithDefault(":8080")
		dbConnect  = env.Get("COUCHBASE_CONNECT").WithDefault("couchbase://localhost")
		dbUsername = env.Get("COUCHBASE_USER").WithDefault("admin")
		dbPassword = env.Get("COUCHBASE_PASS").WithDefault("password")
		dbBucket   = env.Get("COUCHBASE_BUCKET").WithDefault("bucket")
	)

	h, err := api.NewHandler(&api.Config{
		DB: &database.Config{
			ConnectString: dbConnect,
			Username:      dbUsername,
			Password:      dbPassword,
			Bucket:        dbBucket,
		},
	})
	if err != nil {
		logger.Fatalf("initializing API failed: %v", err)
	}

	logger.Println("serving on", addr)
	logger.Fatal(http.ListenAndServe(addr, h))
}
