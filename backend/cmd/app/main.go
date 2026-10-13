package main

import (
	"fmt"
	"log"
	"github.com/jmoiron/sqlx"
	"github.com/lms-project/backend/cmd/app/config"
)

func main() {

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.Config.Database.DB_USER,
		config.Config.Database.DB_PASSWORD,
		config.Config.Database.DB_HOST,
		config.Config.Database.DB_PORT,
		config.Config.Database.DB_NAME,
		config.Config.Database.DB_SSLMODE,
	)

	db, err := sqlx.Connect("pgx", dsn)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()




	log.Println("user created successfully")

}
