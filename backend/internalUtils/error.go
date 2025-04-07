package internalUtils

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func HandleDBQueryError(ctx *gin.Context, err error) {
	if err != nil {
		if strings.Contains(err.Error(), "no such table") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Database is not populated"})
			return
		}
		logger.Error("Error executing query: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
