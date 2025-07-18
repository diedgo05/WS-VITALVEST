package routes

import (
	"ws-vitalvest/WEBSOCKET/infraestructure/controllers"
	"github.com/gin-gonic/gin"
)

// RegisterWSEndpoints registra las rutas WebSocket usando Gin
func RegisterWSEndpoints(router *gin.Engine, wsController *controllers.WebSocketController) {
	// Ruta WebSocket
	router.GET("/ws", wsController.HandleWebSocket)

	// Ruta para recibir datos y enviarlos a WebSocket
	router.POST("/sendData", wsController.HandleSendData)

	// Ruta de status del WebSocket
	router.GET("/ws-status", wsController.HandleStatus)
}
