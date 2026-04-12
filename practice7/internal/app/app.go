package app

import (
	"log"
	"practice-7/config"
	v1 "practice-7/internal/controller/http/v1"
	"practice-7/internal/entity"
	"practice-7/internal/usecase"
	"practice-7/internal/usecase/repo"
	"practice-7/pkg/logger"
	"practice-7/pkg/postgres"

	"github.com/gin-gonic/gin"
)

func Run() {
	cfg := config.NewConfig()
	l := logger.New()

	pg, err := postgres.New(&cfg.PG)
	if err != nil {
		log.Fatalf("postgres connect error: %v", err)
	}
	l.Info("Connected to PostgreSQL")

	// Auto-migrate User table
	if err := pg.Conn.AutoMigrate(&entity.User{}); err != nil {
		log.Fatalf("automigrate error: %v", err)
	}
	l.Info("Database migrated")

	userRepo := repo.NewUserRepo(pg)
	userUseCase := usecase.NewUserUseCase(userRepo)

	router := gin.Default()
	v1.NewRouter(router, userUseCase, l)

	addr := ":" + cfg.HTTP.Port
	l.Info("Server starting on http://localhost" + addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
