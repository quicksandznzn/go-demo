// Created by quicksandzn@gmail.com on 2018/7/25
package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/quicksandznzn/go-demo/go-gin/web"
)


func main() {
	// Disable Console Color
	// gin.DisableConsoleColor()

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()
	router.GET("/welcome", func(context *gin.Context) {
		firstname := context.DefaultQuery("firstname", "Guest")
		lastname := context.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")
		context.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})
	//router.POST("/somePost", posting)
	//router.PUT("/somePut", putting)
	//router.DELETE("/someDelete", deleting)
	//router.PATCH("/somePatch", patching)
	//router.HEAD("/someHead", head)
	//router.OPTIONS("/someOptions", options)

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run(":8881")
	// router.Run(":3000") for a hard coded port

}


func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(err)
	}
}
