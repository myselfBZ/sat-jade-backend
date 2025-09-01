package db

import (
	"testing"
)

func TestNewDB(t *testing.T) {
	config := Config{
		Addr:        "host=localhost port=32768 user=postgres password=new_password dbname=satjade sslmode=disable",
		MaxConns:    15,
		MinConns:    15,
		MaxIdleTime: "15m",
	}

	db, err := New(config)
	if err != nil {
		t.Fatalf("New() returned an error %v\n", err)
	}

	if db == nil {
		t.Fatal("db is nil")
	}
}
