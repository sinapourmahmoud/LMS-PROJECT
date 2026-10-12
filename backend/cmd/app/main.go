package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lms-project/backend/cmd/app/config"
	"github.com/lms-project/backend/internal/user/domain/entity"
	"github.com/lms-project/backend/internal/user/infrastructure/postgres"
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

	repo := postgres.NewUserRepository(db)

	testUser := entity.NewUser(uuid.New(), "amoo sina", "sinapoormahmood@gmail.com", "123456789", "teacher", time.Now(), nil)

	if err := repo.Create(context.Background(), testUser); err != nil {
		log.Fatalf("creating failed, %v", err.Message)
	}

	log.Println("user created successfully")

}
