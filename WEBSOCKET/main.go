package main

import (
	wsControllers "ws-vitalvest/WEBSOCKET/infraestructure/controllers"
	wsAdapters "ws-vitalvest/WEBSOCKET/infraestructure/adapters"
	"log"
	"net/http"

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


	//WEBSOCKET
	wsServer := wsAdapters.NewWebSocketServer()
	go wsServer.Run()
	wsController := wsControllers.NewWebSocketController(wsServer)

	// Registrar rutas WebSocket directamente en main
	router.GET("/ws", wsController.HandleWebSocket)
	router.POST("/sendData", wsController.HandleSendData)
	router.GET("/ws-status", wsController.HandleStatus)

	// LOGIN SIMPLE
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

	port := ":8080"
	log.Println("🚀 Servidor corriendo en el puerto", port)
	log.Println("📡 WebSocket disponible en ws://localhost:8080/ws")
	log.Println("📤 Endpoint sendData disponible en http://localhost:8080/sendData")
	log.Println("📊 Status WebSocket en http://localhost:8080/ws-status")
	log.Println("🔐 Login en http://localhost:8080/login")
	log.Fatal(router.Run(port))
}
