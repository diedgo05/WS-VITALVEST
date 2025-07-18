package domain

type ISensors interface {
	SendData(data Sensors) error
}
