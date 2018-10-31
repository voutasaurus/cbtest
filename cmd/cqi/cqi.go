package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	c := &client{api: "http://localhost:8080"}
	cmd := ""
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}
	os.Args = os.Args[1:] // skip this for the subcommands
	switch cmd {
	case "read":
		do(c.read)
	case "write":
		do(c.write)
	default:
		log.Printf("cannot run command: %q", cmd)
		log.Fatalf("available commands: read, write")
	}
}

func do(query func(string, map[string]string) ([]byte, error)) {
	stmt, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	args, err := flags(string(stmt))
	if err != nil {
		log.Fatal(err)
	}
	b, err := query(string(stmt), args)
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout.Write(b)
}

func flags(query string) (map[string]string, error) {
	m := make(map[string]*string)
	for i, prearg := range strings.Split(query, "$") {
		if i == 0 {
			// skip the first part of the query
			continue
		}
		// split on the first non alphabetic rune
		arg := strings.FieldsFunc(prearg, func(r rune) bool {
			return !(('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z'))
		})[0]
		if _, ok := m[arg]; !ok {
			m[arg] = flag.String(arg, "", "")
		}
	}
	flag.Parse()
	mm := make(map[string]string)
	for arg, p := range m {
		if p == nil {
			return nil, fmt.Errorf("flag %q must be set", arg)
		}
		mm[arg] = *p
	}
	return mm, nil
}

type client struct {
	api string
}

func (c *client) read(stmt string, args map[string]string) ([]byte, error) {
	return c.send("/read", stmt, args)
}

func (c *client) write(stmt string, args map[string]string) ([]byte, error) {
	return c.send("/write", stmt, args)
}

func (c *client) send(ep, stmt string, args map[string]string) ([]byte, error) {
	b, err := json.Marshal(struct {
		Stmt string
		Args map[string]string
	}{
		Stmt: stmt,
		Args: args,
	})
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(b)
	res, err := http.Post(c.api+ep, "application/json", body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		io.Copy(os.Stderr, res.Body)
		return nil, fmt.Errorf("Error: %v", res.StatusCode)
	}
	return ioutil.ReadAll(res.Body)
}
