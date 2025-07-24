package routes

import (
	"ws-vitalvest/WEBSOCKET/infraestructure/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RegisterWSEndpoints registra las rutas WebSocket usando Gin
func RegisterWSEndpoints(router *gin.Engine, wsController *controllers.WebSocketController) {
	// Ruta WebSocket
	router.GET("/ws", wsController.HandleWebSocket)

	// Ruta para recibir datos y enviarlos a WebSocket
	router.POST("/sendData", wsController.HandleSendData)

	// Ruta de status del WebSocket
	router.GET("/ws-status", wsController.HandleStatus)

	// Ruta alertas
	router.POST("/sendAlerta", wsController.HandleSendAlerta)

	router.POST("/login", func(c *gin.Context) {
		var loginRequest struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
			return
		}

		// Validación simple (reemplaza con tu lógica real)
		validUsers := map[string]string{
			"admin":  "admin123",
			"juan":   "juan123",
			"maria":  "maria123",
			"carlos": "carlos123",
		}

		if password, exists := validUsers[loginRequest.Username]; exists && password == loginRequest.Password {
			// Usuario válido - devolver datos del usuario
			userData := []map[string]interface{}{
				{
					"id":       1,
					"username": loginRequest.Username,
					"name":     loginRequest.Username,
					"role":     "user",
				},
			}
			c.JSON(http.StatusOK, userData)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales incorrectas"})
		}
	})
}
