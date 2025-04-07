package scrape

import (
	"strings"

	"github.com/LeonEstrak/retro-drop/backend/internalTypes"
	"github.com/LeonEstrak/retro-drop/backend/internalUtils"
	"github.com/gocolly/colly/v2"
)

var logger = internalUtils.GetLogger()

func ScrapeListOfGames(url string) []string {
	listOfGames := []string{}

	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {
		logger.Debug("Visited %s", r.Request.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		logger.Debug("Error: %v", err)
	})

	c.OnHTML("a", func(e *colly.HTMLElement) {
		gameTitle := e.Attr("title")

		if gameTitle == "" || !strings.HasSuffix(gameTitle, ".zip") {
			return
		}
		gameTitle = strings.TrimSuffix(gameTitle, ".zip")

		logger.Debug("Found Title: %s", gameTitle)

		listOfGames = append(listOfGames, gameTitle)
	})

	c.Visit(url)
	return listOfGames
}

func ScrapeAllDownloadLinks(url string) []internalTypes.Games {
	gameObjects := []internalTypes.Games{}

	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {
		logger.Debug("Visited %s", r.Request.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		logger.Debug("Error: %v", err)
	})

	c.OnHTML("a", func(e *colly.HTMLElement) {
		gameTitle := e.Attr("title")
		gameDownloadLink := e.Attr("href")

		if gameTitle == "" || !strings.HasSuffix(gameTitle, ".zip") {
			return
		}
		gameTitle = strings.TrimSuffix(gameTitle, ".zip")

		logger.Debug("Found Title: %s", gameTitle)

		gameObjects = append(gameObjects, internalTypes.Games{
			GameTitle:   gameTitle,
			DownloadURL: gameDownloadLink,
		})
	})

	c.Visit(url)
	return gameObjects
}
