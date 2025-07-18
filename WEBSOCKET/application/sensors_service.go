package app

import (
	"ws-vitalvest/WEBSOCKET/domain"
	"ws-vitalvest/WEBSOCKET/infraestructure/controllers"
)

type SensorService struct {
	WebSocketController *controllers.WebSocketController
}

func NewSensorService(wsController *controllers.WebSocketController) *SensorService {
	return &SensorService{WebSocketController: wsController}
}

func (s *SensorService) EnviarDatos(data domain.Sensors) error {
	return s.WebSocketController.SendData(data)
}
