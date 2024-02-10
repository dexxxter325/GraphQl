package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func ConnTODB() (*pgxpool.Pool, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("err in load env file:%s", err.Error())
	}
	var (
		host     = os.Getenv("HOST")
		port     = os.Getenv("PORT")
		user     = os.Getenv("USER")
		dbname   = os.Getenv("DB_NAME")
		password = os.Getenv("PASSWORD")
		sslmode  = os.Getenv("SSLMODE")
	)
	data := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", host, port, user, dbname, password, sslmode)
	conn, err := pgxpool.New(context.Background(), data)
	if err != nil {
		log.Fatalf("err in conn to db:%s", err)
	}
	return conn, nil
}
