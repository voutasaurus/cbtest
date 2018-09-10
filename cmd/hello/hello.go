package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/voutasaurus/cbtest/env"
)

func main() {
	logger := log.New(os.Stderr, "hello: ", log.Llongfile|log.Lmicroseconds|log.LstdFlags)
	logger.Println("starting...")

	var (
		addr = env.Get("HELLO_ADDR").WithDefault(":8080")
		_    = env.Get("HELLO_DATABASE_ADDR").WithDefault(":12345")
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Println("hit")
		fmt.Fprintln(w, "hello")
	})

	logger.Println("serving on", addr)
	logger.Fatal(http.ListenAndServe(addr, mux))
}
