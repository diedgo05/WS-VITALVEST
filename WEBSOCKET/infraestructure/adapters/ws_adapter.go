package adapters

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type WebSocketServer struct {
	Clients    map[*websocket.Conn]bool
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
	Broadcast  chan interface{}
}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		Clients:    make(map[*websocket.Conn]bool),
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),
		Broadcast:  make(chan interface{}, 100),
	}
}

func (s *WebSocketServer) Run() {
	log.Println("WebSocket Server iniciado")

	for {
		select {
		case conn := <-s.Register:
			s.Clients[conn] = true
			log.Printf("Nueva conexión registrada. Total clientes: %d", len(s.Clients))

		case conn := <-s.Unregister:
			if _, ok := s.Clients[conn]; ok {
				delete(s.Clients, conn)
				conn.Close()
				log.Printf("Conexión cerrada. Total clientes: %d", len(s.Clients))
			}

		case data := <-s.Broadcast:
			log.Printf("Enviando datos a %d clientes: %+v", len(s.Clients), data)

			if jsonData, err := json.Marshal(data); err == nil {
				log.Printf("JSON enviado: %s", string(jsonData))
			}

			for client := range s.Clients {
				if err := client.WriteJSON(data); err != nil {
					log.Printf("Error al enviar mensaje al cliente: %v", err)
					client.Close()
					delete(s.Clients, client)
				}
			}
		}
	}
}

func (s *WebSocketServer) SendData(data interface{}) error {
	select {
	case s.Broadcast <- data:
		return nil
	default:
		log.Println("Canal de broadcast lleno, mensaje descartado")
		return nil
	}
}

func (s *WebSocketServer) GetClientsCount() int {
	return len(s.Clients)
}

func (s *WebSocketServer) GetClients() map[*websocket.Conn]bool {
	return s.Clients
}
