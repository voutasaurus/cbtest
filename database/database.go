package database

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/couchbase/gocb"
)

type DB struct {
	db *gocb.Bucket

	// bucket is the escaped bucket name connected to by db.
	bucket string
}

type Config struct {
	ConnectString string
	Username      string
	Password      string
	Bucket        string
}

func NewDB(ctx context.Context, c *Config) (*DB, error) {
	db, err := newDB(c)
	if err == nil {
		return db, nil
	}

	// If no timeout is set, exit early
	_, ok := ctx.Deadline()
	if !ok {
		return nil, err
	}

	// Wait for DB to become available (or timeout)
	for {
		time.Sleep(100 * time.Millisecond)
		db, err := newDB(c)
		if err == nil {
			return db, nil
		}
		select {
		case <-ctx.Done():
			return nil, err
		default:
		}
	}
}

var errBucketName = errors.New("invalid bucket name, cannot contain `")

func newDB(c *Config) (*DB, error) {
	if strings.ContainsRune(c.Bucket, '`') {
		return nil, errBucketName
	}
	cluster, err := gocb.Connect(c.ConnectString)
	if err != nil {
		return nil, fmt.Errorf("db cluster connect error: %v", err)
	}
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: c.Username,
		Password: c.Password,
	})
	bucket, err := cluster.OpenBucket(c.Bucket, "")
	if err != nil {
		return nil, fmt.Errorf("db bucket open error: %v", err)
	}
	// TODO: bucket.Ping to verify?
	return &DB{db: bucket, bucket: "`" + c.Bucket + "`"}, nil
}

// Close is not graceful. If in flight operations are important then ensure
// they are complete before calling Close.
//
// TODO: gracefully shut down database connection
func (db *DB) Close() error {
	return db.db.Close()
}

type Record struct {
	ID  string
	Doc struct {
		Type string `json:"type"`
		X    string `json:"x"`
	}
}

func (db *DB) New(ctx context.Context, r *Record) (*Record, error) {
	// TODO: use InsertDura to validate durability
	// TODO: retain the CAS?
	// TODO: add timeout from context ctx.Deadline()
	_, err := db.db.Insert(r.ID, r.Doc, 0)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (db *DB) Get(ctx context.Context, r *Record) (*Record, error) {
	_, err := db.db.Get(r.ID, &r.Doc)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (db *DB) query(stmt string, args map[string]interface{}) (gocb.QueryResults, error) {
	q := gocb.NewN1qlQuery(stmt).ReadOnly(true)
	return db.db.ExecuteN1qlQuery(q, args)
}

func (db *DB) exec(stmt string, args map[string]interface{}) (gocb.QueryResults, error) {
	q := gocb.NewN1qlQuery(stmt)
	return db.db.ExecuteN1qlQuery(q, args)
}

func (db *DB) Search(ctx context.Context, r *Record) ([]*Record, error) {
	// TODO: add timeout from context
	rows, err := db.query(
		`SELECT {type,x} AS doc,
		  META(`+db.bucket+`).id AS id
		  FROM `+db.bucket+`
		  WHERE type == $type AND x LIKE $prefix`,
		map[string]interface{}{
			"type":   "x",
			"prefix": r.Doc.X + "%",
		},
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rr []*Record
	for next := true; next; {
		var rec Record
		if next = rows.Next(&rec); next {
			rr = append(rr, &rec)
		}
	}
	return rr, nil
}

func (db *DB) Update(ctx context.Context, r *Record) (*Record, error) {
	return nil, nil
}

func (db *DB) Delete(ctx context.Context, r *Record) (*Record, error) {
	return nil, nil
}
