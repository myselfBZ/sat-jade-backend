package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	"github.com/myselfBZ/sat-jade/internal/db"
	"golang.org/x/crypto/bcrypt"
)

func seedAdmin(pool *pgxpool.Pool) error {
	q := `INSERT INTO 
			users(full_name, email, role, password_hash) 
		VALUES ('The Super User', 'theadmin@satjade.com', 'admin', $1)`
	hash, err := bcrypt.GenerateFromPassword([]byte("restfulapis1A"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = pool.Exec(context.TODO(), q, hash)
	if err != nil {
		return err
	}
	log.Println("Admin has been created successfully")
	return nil
}

func main() {
	pool, err := db.New(db.Config{
		Addr:        os.Getenv("DB"),
		MaxConns:    15,
		MinConns:    15,
		MaxIdleTime: "15m",
	})
	if err != nil {
		panic(err)
	}
	if err := seedAdmin(pool); err != nil {
		panic(err)
	}
}
