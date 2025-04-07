package main

import (
	"github.com/LeonEstrak/retro-drop/backend/database"
	"github.com/LeonEstrak/retro-drop/backend/internalUtils"
	"github.com/LeonEstrak/retro-drop/backend/routes"

	"fmt"
)

func main() {
	PORT := 9090

	logger := internalUtils.GetLogger()

	// Initialize the DB connection
	db := database.GetDB()
	defer db.Close()

	// Initialize the routes
	router := routes.SetupRoutes()

	logger.Debug("Starting server on port %d", PORT)
	router.Run(":" + fmt.Sprint(PORT))
}
