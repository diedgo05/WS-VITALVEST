package adapters

import (
	"encoding/json"
	"log"
	"ws-vitalvest/WEBSOCKET/domain"

	"github.com/gorilla/websocket"
)

// WebSocketServer gestiona las conexiones WebSocket.
type WebSocketServer struct {
	Clients    map[*websocket.Conn]bool
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
	Broadcast  chan domain.Sensors
}

// NewWebSocketServer crea una nueva instancia de WebSocketServer.
func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		Clients:    make(map[*websocket.Conn]bool),
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),
		Broadcast:  make(chan domain.Sensors),
	}
}

// Run inicia el servidor WebSocket y maneja la conexión de los clientes.
func (s *WebSocketServer) Run() {
	log.Println("🚀 WebSocket Server iniciado y esperando conexiones...")

	for {
		select {
		case conn := <-s.Register:
			s.Clients[conn] = true
			log.Printf("✅ Nueva conexión WebSocket registrada. Total clientes: %d", len(s.Clients))

		case conn := <-s.Unregister:
			if _, ok := s.Clients[conn]; ok {
				delete(s.Clients, conn)
				conn.Close()
				log.Printf("❌ Conexión WebSocket cerrada. Total clientes: %d", len(s.Clients))
			}

		case data := <-s.Broadcast:
			log.Printf("📡 Enviando datos a %d clientes: %+v", len(s.Clients), data)

			// Convertir datos a JSON para logging
			jsonData, _ := json.Marshal(data)
			log.Printf("📦 JSON enviado: %s", string(jsonData))

			// Enviar a todos los clientes conectados
			for client := range s.Clients {
				err := client.WriteJSON(data)
				if err != nil {
					log.Printf("❌ Error al enviar mensaje al cliente: %v", err)
					client.Close()
					delete(s.Clients, client)
				}
			}
		}
	}
}

// SendData envía datos a todos los clientes conectados
func (s *WebSocketServer) SendData(data domain.Sensors) error {
	select {
	case s.Broadcast <- data:
		return nil
	default:
		log.Println("⚠️ Canal de broadcast lleno, descartando mensaje")
		return nil
	}
}

// GetClientsCount devuelve el número de clientes conectados
func (s *WebSocketServer) GetClientsCount() int {
	return len(s.Clients)
}

// GetClients devuelve el mapa de clientes (solo para testing/debugging)
func (s *WebSocketServer) GetClients() map[*websocket.Conn]bool {
	return s.Clients
}
