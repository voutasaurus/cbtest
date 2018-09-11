package database

import (
	"context"

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
		return nil, err
	}
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: c.Username,
		Password: c.Password,
	})
	bucket, err := cluster.OpenBucket(c.Bucket, c.Password)
	if err != nil {
		return nil, err
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
