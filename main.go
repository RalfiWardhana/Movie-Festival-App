package main

import (
	"log"
	"movies/config"
	"movies/handler"
	"movies/repository"
	"movies/router"
	"movies/usecase"
)

func main() {
	// Initialize MySQL database connection
	db, err := config.InitMysql()

	if err != nil {
		log.Println("Failed connect to DB") // Log error if DB connection fails
		return
	}

	// Set up repository, use case, and handler for movie-related functionality
	movieRepo := repository.NewMoviesRepo(db)
	movieUseCase := usecase.NewMoviesUseCase(movieRepo)
	movieHandler := handler.NewMoviesHandler(movieUseCase)

	// Set up repository, use case, and handler for stats-related functionality
	statsRepo := repository.NewStatsRepo(db)
	statsUseCase := usecase.NewStatsUseCase(statsRepo)
	statsHandler := handler.NewStatsHandler(statsUseCase)

	// Set up repository, use case, and handler for user-related functionality
	userRepo := repository.NewUsersRepo(db)
	userUseCase := usecase.NewUsersUseCase(userRepo)
	userHandler := handler.NewUsersHandler(userUseCase)

	// Initialize router with handlers
	r := router.Router(movieHandler, statsHandler, userHandler)

	// Start the server on port 9191
	err = r.Run(":9191")

	if err != nil {
		log.Println("Failed to starting server") // Log error if server fails to start
		return
	}

	log.Println("Server starting at port 9191") // Log success message
}
