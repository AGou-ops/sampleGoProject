// main.go

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	// Handle Index
	router.GET("/", showIndexPage)
	// Handle GET requests at /article/view/some_article_id
	router.GET("/article/view/:article_id", getArticle)

	router.Use(func(ctx *gin.Context) {
		if ctx.Request.Method != http.MethodGet {
			ctx.JSON(http.StatusMethodNotAllowed, gin.H{
				"msg": "Method Not Allowed",
			})
		}
	})

	router.Run()
}
