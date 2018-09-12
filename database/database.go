package database

import (
	"context"
	"fmt"

	"github.com/couchbase/gocb"
)

type DB struct {
	db *gocb.Bucket
}

type Config struct {
	ConnectString string
	Username      string
	Password      string
	Bucket        string
}

func NewDB(c *Config) (*DB, error) {
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
	return &DB{db: bucket}, nil
}

type Record struct {
	Key   string
	Value string
}

func (db *DB) New(ctx context.Context, r *Record) (*Record, error) {
	// TODO: use InsertDura to validate durability
	// TODO: retain the CAS?
	// TODO: add timeout from context
	_, err := db.db.Insert(r.Key, r.Value, 0)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (db *DB) Get(ctx context.Context, r *Record) (*Record, error) {
	_, err := db.db.Get(r.Key, &r.Value)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (db *DB) Search(ctx context.Context, r *Record) ([]*Record, error) {
	// TODO: add timeout from context
	// gocb.NewN1qlQuery(``).ReadOnly(true)
	return nil, nil
}

func (db *DB) Update(ctx context.Context, r *Record) (*Record, error) {
	return nil, nil
}

func (db *DB) Delete(ctx context.Context, r *Record) (*Record, error) {
	return nil, nil
}
