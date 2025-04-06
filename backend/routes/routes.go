package routes

import (
	"maps"
	"net/http"
	"slices"
	"strings"

	"github.com/LeonEstrak/retro-drop/backend/constants"
	"github.com/LeonEstrak/retro-drop/backend/database"
	"github.com/LeonEstrak/retro-drop/backend/scrape"
	"github.com/LeonEstrak/retro-drop/backend/utils"
	"github.com/gin-gonic/gin"
)

var (
	route_map = []gin.RouteInfo{
		{
			Method:      "GET",
			Path:        "/games",
			HandlerFunc: getGames(),
		},
		{
			Method:      "GET",
			Path:        "/systems",
			HandlerFunc: getSystems(),
		},
		{
			Method:      "POST",
			Path:        "/init",
			HandlerFunc: initializeDB(),
		},
	}
)

var logger = utils.GetLogger()

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	for _, route := range route_map {
		router.Handle(route.Method, route.Path, route.HandlerFunc)
	}

	return router
}

func getGames() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		system := ctx.Query("system")

		db := database.OpenDB(constants.DB_PATH)
		defer db.Close()

		responseBody := []database.Games{}
		if len(system) == 0 {
			row, err := db.Query("SELECT id, game_title, system, download_url FROM games")
			if err != nil {
				logger.Error("Error executing query: %v", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer row.Close()

			for row.Next() {
				var id int
				var gameTitle string
				var downloadURL string
				var system string
				if err := row.Scan(&id, &gameTitle, &system, &downloadURL); err != nil {
					logger.Error("Error scanning row: %v", err)
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				responseBody = append(responseBody, database.Games{
					ID:          id,
					GameTitle:   gameTitle,
					System:      system,
					DownloadURL: downloadURL,
				})
			}
			ctx.JSON(http.StatusOK, gin.H{"games": responseBody})

			ctx.JSON(http.StatusOK, row)
		}

		rows, err := db.Query("SELECT game_title, download_url FROM games WHERE system = ?", system)
		if err != nil {
			if strings.Contains(err.Error(), "no such table") {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Database is not populated"})
				return
			}
			logger.Error("Error executing query: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()
	}
}

func getSystems() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		listOfSystems := slices.Collect(maps.Keys(constants.SYSTEMS_TO_ERISTA_MAPPING))
		ctx.JSON(http.StatusOK, listOfSystems)
	}
}

func initializeDB() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := database.OpenDB(constants.DB_PATH)
		defer db.Close()

		dropDbErr := database.DropTables(db)
		if dropDbErr != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": dropDbErr.Error(), "message": "Failed to drop tables"})
			return
		}
		createDbErr := database.CreateTables(db)
		if createDbErr != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": createDbErr.Error(), "message": "Failed to create tables"})
			return
		}

		listOfSystems := slices.Collect(maps.Keys(constants.SYSTEMS_TO_ERISTA_MAPPING))
		for _, system := range listOfSystems {
			logger.Debug("Scraping %s", system)
			listOfGames := scrape.ScrapeAllDownloadLinks(constants.MYRIENT_ERISTA_ME_BASE_URL + constants.SYSTEMS_TO_ERISTA_MAPPING[system])
			for gameTitle, downloadUrl := range listOfGames {
				logger.Debug("Scraped %s", gameTitle)
				_, err := db.Exec("INSERT INTO games (game_title, system, download_url) VALUES (?, ?, ?)", gameTitle, system, downloadUrl)
				if err != nil {
					logger.Error("Error executing query: %v", err)
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Database populated successfully"})
	}
}
