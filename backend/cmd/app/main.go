package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/lms-project/backend/cmd/app/config"
	courseServices "github.com/lms-project/backend/internal/course/domain/services"
	coursePostgres "github.com/lms-project/backend/internal/course/infrastructure/postgres"
	coursehttp "github.com/lms-project/backend/internal/course/interfaces/http"
	enrollmentServices "github.com/lms-project/backend/internal/enrollment/domain/services"
	enrollmentPostgres "github.com/lms-project/backend/internal/enrollment/infrastructure/postgres"
	enrollmenthttp "github.com/lms-project/backend/internal/enrollment/interfaces/http"
	"github.com/lms-project/backend/internal/user/domain/entity"
	"github.com/lms-project/backend/internal/user/domain/services"
	"github.com/lms-project/backend/internal/user/infrastructure/postgres"
	userhttp "github.com/lms-project/backend/internal/user/interfaces/http"
	"github.com/lms-project/backend/pkg/middleware"
)

func main() {

	if config.Config.Auth.JWTAccessSecret == "" || config.Config.Auth.JWTRefreshSecret == "" {
		log.Fatal("JWT_ACCESS_SECRET and JWT_REFRESH_SECRET must be set")
	}

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
	refreshTokenRepository := postgres.NewRefreshTokenRepository(db)
	userCommandServices := services.NewUserCommandDomainServices(userRepository)
	userQuesryServices := services.NewUserQueryServices(userRepository, refreshTokenRepository)
	userHandler := userhttp.NewUserHandler(userCommandServices, userQuesryServices)
	courseRepository := coursePostgres.NewCourseRepository(db)
	courseCommandServices := courseServices.NewCourseCommandServices(courseRepository, userRepository)
	courseQueryServices := courseServices.NewCourseQueryServices(courseRepository, userRepository)
	courseHandler := coursehttp.NewCourseHandler(courseCommandServices, courseQueryServices)
	enrollmentRepository := enrollmentPostgres.NewEnrollmentRepository(db)
	enrollmentCommandServices := enrollmentServices.NewEnrollmentCommandServices(enrollmentRepository, courseRepository)
	enrollmentQueryServices := enrollmentServices.NewEnrollmentQueryServices(enrollmentRepository, courseRepository)
	enrollmentHandler := enrollmenthttp.NewEnrollmentHandler(enrollmentCommandServices, enrollmentQueryServices)

	router := gin.Default()

	//routes

	// user
	router.POST("/signup", userHandler.SignUp())
	router.POST("/login", userHandler.Login())
	router.GET("/whoami", middleware.AuthMiddleware(userQuesryServices), func(c *gin.Context) {
		c.JSON(200, gin.H{"userID": c.MustGet("userID")})
	})

	router.GET("teacher-only", middleware.AuthMiddleware(userQuesryServices), middleware.RequireRole(entity.RoleTeacher), func(c *gin.Context) {
		c.JSON(200, gin.H{"userID": c.MustGet("userID")})

	})

	router.POST("/refresh", userHandler.Refresh())
	router.POST("/logout", userHandler.Logout())

	// course
	router.GET("/courses", courseHandler.ListPublishedCourses())
	router.GET("/courses/:id", courseHandler.GetCourse())

	router.POST("/courses", middleware.AuthMiddleware(userQuesryServices), middleware.RequireRole(entity.RoleTeacher), courseHandler.CreateCourse())
	router.GET("/courses/mine", middleware.AuthMiddleware(userQuesryServices), middleware.RequireRole(entity.RoleTeacher), courseHandler.ListMyCourses())

	router.GET("/courses/pending", middleware.AuthMiddleware(userQuesryServices), middleware.RequireRole(entity.RoleAdmin), courseHandler.ListPendingCourses())
	router.PATCH("/courses/:id/approve", middleware.AuthMiddleware(userQuesryServices), middleware.RequireRole(entity.RoleAdmin), courseHandler.ApproveCourse())

	router.DELETE("/courses/:id", middleware.AuthMiddleware(userQuesryServices), middleware.RequireRole(entity.RoleTeacher, entity.RoleAdmin), courseHandler.DeleteCourse())

	log.Println("server listening on :", config.Config.Application.PORT)

	// enrollment
	router.POST("/courses/:id/enroll", middleware.AuthMiddleware(userQuesryServices), middleware.RequireRole(entity.RoleStudent), enrollmentHandler.Enroll())
	router.GET("/enrollments/mine", middleware.AuthMiddleware(userQuesryServices), enrollmentHandler.ListMyEnrollments())
	router.GET("/courses/:id/enrollments", middleware.AuthMiddleware(userQuesryServices), middleware.RequireRole(entity.RoleTeacher, entity.RoleAdmin), enrollmentHandler.ListCourseEnrollments())
	router.DELETE("/courses/:id/enroll", middleware.AuthMiddleware(userQuesryServices), middleware.RequireRole(entity.RoleStudent), enrollmentHandler.Unenroll())

	if err := router.Run(":" + config.Config.Application.PORT); err != nil {
		log.Fatal(err)
	}

}
