package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/lms-project/backend/cmd/app/config"
	"github.com/lms-project/backend/internal/user/domain/services"
	"github.com/lms-project/backend/internal/user/infrastructure/postgres"
	userhttp "github.com/lms-project/backend/internal/user/interfaces/http"
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

	log.Println("database connection established")

	userRepository := postgres.NewUserRepository(db)
	userCommandServices := services.NewUserCommandDomainServices(userRepository)
	userQuesryServices :=services.NewUserQueryServices(userRepository)
	userHandler := userhttp.NewUserHandler(userCommandServices,userQuesryServices)

	router := gin.Default()

	//routes
	router.POST("/signup", userHandler.SignUp())
	router.POST("/login", userHandler.Login())

	log.Println("server listening on :5000")

	if err := router.Run(":5000"); err != nil {
		log.Fatal(err)
	}

}
