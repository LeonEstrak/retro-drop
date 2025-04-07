package routes

import (
	"maps"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/LeonEstrak/retro-drop/backend/constants"
	"github.com/LeonEstrak/retro-drop/backend/database"
	"github.com/LeonEstrak/retro-drop/backend/internalUtils"
	"github.com/LeonEstrak/retro-drop/backend/scrape"
	"github.com/gin-gonic/gin"
)

var (
	route_map = []gin.RouteInfo{
		{
			Method:      "GET",
			Path:        "/games",
			HandlerFunc: getGamesHandler(),
		},
		{
			Method:      "GET",
			Path:        "/systems",
			HandlerFunc: getSystemsHandler(),
		},
		{
			Method:      "POST",
			Path:        "/init",
			HandlerFunc: initDBHandler(),
		},
	}
)

var logger = internalUtils.GetLogger()

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	for _, route := range route_map {
		router.Handle(route.Method, route.Path, route.HandlerFunc)
	}

	return router
}

func getGamesHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		systemQuery := ctx.Query("system")
		limitQuery, err := strconv.Atoi(ctx.Query("limit"))

		if err != nil {
			limitQuery = 0
		}

		db := database.GetDB()

		responseBody, err := db.GetGamesFromDB(systemQuery, limitQuery)
		if err != nil {
			if strings.Contains(err.Error(), "no such table") {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Database is not populated"})
				return
			}
			logger.Error("Error executing query: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		ctx.JSON(http.StatusOK, responseBody)
	}
}

func getSystemsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		listOfSystems := slices.Collect(maps.Keys(constants.SYSTEMS_TO_ERISTA_MAPPING))
		ctx.JSON(http.StatusOK, listOfSystems)
	}
}

func initDBHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := database.GetDB()

		dropDbErr := db.DropTables()
		if dropDbErr != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": dropDbErr.Error(), "message": "Failed to drop tables"})
			return
		}
		createDbErr := db.CreateTables()
		if createDbErr != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": createDbErr.Error(), "message": "Failed to create tables"})
			return
		}

		listOfSystems := slices.Collect(maps.Keys(constants.SYSTEMS_TO_ERISTA_MAPPING))
		for _, system := range listOfSystems {
			logger.Debug("Scraping %s", system)
			listOfGames := scrape.ScrapeAllDownloadLinks(constants.MYRIENT_ERISTA_ME_BASE_URL + constants.SYSTEMS_TO_ERISTA_MAPPING[system])

			err := db.InsertListOfGamesToDB(listOfGames)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed to insert data to db"})
				return
			}
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Database populated successfully"})
	}
}
