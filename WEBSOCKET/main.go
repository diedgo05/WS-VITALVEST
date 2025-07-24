package main

import (
	wsRoutes "ws-vitalvest/WEBSOCKET/infraestructure/routes"
	wsAdapters "ws-vitalvest/WEBSOCKET/infraestructure/adapters"
	wsControllers "ws-vitalvest/WEBSOCKET/infraestructure/controllers"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	wsServer := wsAdapters.NewWebSocketServer()
	go wsServer.Run()

	wsController := wsControllers.NewWebSocketController(wsServer)
	
	wsRoutes.RegisterWSEndpoints(router, wsController)

	port := ":3000"
	log.Fatal(router.Run(port))
}
