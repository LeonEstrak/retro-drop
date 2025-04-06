package main

import (
	"github.com/LeonEstrak/retro-drop/backend/routes"
	"github.com/LeonEstrak/retro-drop/backend/utils"

	"fmt"
)

func main() {
	PORT := 9090

	logger := utils.GetLogger()

	router := routes.SetupRoutes()

	logger.Debug("Starting server on port %d", PORT)
	router.Run(":" + fmt.Sprint(PORT))
}
