package controllers

import (
	"log"
	"net/http"
	"ws-vitalvest/WEBSOCKET/domain"
	"ws-vitalvest/WEBSOCKET/infraestructure/adapters"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketController struct {
	WebSocketServer *adapters.WebSocketServer
}

func NewWebSocketController(wsServer *adapters.WebSocketServer) *WebSocketController {
	return &WebSocketController{WebSocketServer: wsServer}
}

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

	log.Println("Nueva conexión WebSocket establecida")
	c.WebSocketServer.Register <- conn

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Conexión WebSocket cerrada:", err)
			c.WebSocketServer.Unregister <- conn
			break
		}
	}
}

func (c *WebSocketController) HandleSendData(ctx *gin.Context) {
	var data domain.Sensors

	if err := ctx.ShouldBindJSON(&data); err != nil {
		log.Println("Error al leer el cuerpo de la solicitud:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al leer el cuerpo de la solicitud"})
		return
	}

	if err := c.SendData(data); err != nil {
		log.Println("Error al enviar el data:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al enviar el data"})
		return
	}

	log.Printf("Datos enviados via WebSocket: %+v", data)
	ctx.JSON(http.StatusOK, gin.H{"message": "datos enviados exitosamente"})
}

func (c *WebSocketController) SendData(data domain.Sensors) error {
	c.WebSocketServer.Broadcast <- data
	return nil
}

func (c *WebSocketController) HandleStatus(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message":           "WebSocket server está funcionando",
		"clients_connected": c.WebSocketServer.GetClientsCount(),
		"websocket_url":     "ws://localhost:8080/ws",
		"send_data_url":     "http://localhost:8080/sendData",
	})
}

func (c *WebSocketController) HandleSendAlerta(ctx *gin.Context) {
	var data domain.Alerta

	if err := ctx.ShouldBindJSON(&data); err != nil {
		log.Println("Error al leer el cuerpo de la solicitud:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al leer el cuerpo de la solicitud"})
		return
	}

	if err := c.SendAlerta(data); err != nil {
		log.Println("Error al enviar el data:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al enviar el data"})
		return
	}

	log.Printf("Datos enviados via WebSocket: %+v", data)
	ctx.JSON(http.StatusOK, gin.H{"message": "datos enviados exitosamente"})
}

func (c *WebSocketController) SendAlerta(data domain.Alerta) error {
	c.WebSocketServer.Broadcast <- data
	return nil
}
