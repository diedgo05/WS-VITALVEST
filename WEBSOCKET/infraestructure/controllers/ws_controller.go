package controllers

import (
	"ws-vitalvest/WEBSOCKET/infraestructure/adapters"
	"ws-vitalvest/WEBSOCKET/domain"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocketController maneja las conexiones WebSocket.
type WebSocketController struct {
	WebSocketServer *adapters.WebSocketServer
}

// NewWebSocketController crea un nuevo controlador WebSocket.
func NewWebSocketController(wsServer *adapters.WebSocketServer) *WebSocketController {
	return &WebSocketController{WebSocketServer: wsServer}
}

// HandleWebSocket maneja la conexión WebSocket con Gin Context.
func (c *WebSocketController) HandleWebSocket(ctx *gin.Context) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("Error al establecer WebSocket:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No se pudo establecer conexión WebSocket"})
		return
	}
	defer conn.Close()

	log.Println("✅ Nueva conexión WebSocket establecida")
	c.WebSocketServer.Register <- conn

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("🔌 Conexión WebSocket cerrada:", err)
			c.WebSocketServer.Unregister <- conn
			break
		}
	}
}

// HandleSendData maneja el endpoint POST /sendData para recibir datos y enviarlos via WebSocket
func (c *WebSocketController) HandleSendData(ctx *gin.Context) {
	var data domain.Sensors

	if err := ctx.ShouldBindJSON(&data); err != nil {
		log.Println("Error al leer el cuerpo de la solicitud:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al leer el cuerpo de la solicitud"})
		return
	}

	// Llamar al servicio WebSocket para enviar el data
	if err := c.SendData(data); err != nil {
		log.Println("Error al enviar el data:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al enviar el data"})
		return
	}

	log.Printf("Datos enviados via WebSocket: %+v", data)
	ctx.JSON(http.StatusOK, gin.H{"message": "datos enviados exitosamente"})
}

// SendData envía datos a través del WebSocket.
func (c *WebSocketController) SendData(data domain.Sensors) error {
	c.WebSocketServer.Broadcast <- data
	return nil
}

// HandleStatus devuelve el estado del WebSocket server
func (c *WebSocketController) HandleStatus(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message":           "WebSocket server está funcionando",
		"clients_connected": c.WebSocketServer.GetClientsCount(),
		"websocket_url":     "ws://localhost:8080/ws",
		"send_data_url":     "http://localhost:8080/sendData",
	})
}
