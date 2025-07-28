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
	log.Println("🔧 Inicializando servidor...")

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	wsServer := wsAdapters.NewWebSocketServer()
	go func() {
		log.Println("🧠 Servidor WebSocket iniciado")
		wsServer.Run()
	}()

	wsController := wsControllers.NewWebSocketController(wsServer)

	wsRoutes.RegisterWSEndpoints(router, wsController)

	port := ":3000"
	log.Println("🚀 Servidor HTTP iniciado en http://localhost" + port)
	log.Fatal(router.Run(port))
}
