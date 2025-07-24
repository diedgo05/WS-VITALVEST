package domain

import "time"

type Alerta struct {
	ID                     int       `json:"id"`
	NombreDelSensor        string    `json:"nombre_del_sensor"`
	Fecha                  time.Time `json:"fecha"`
	CantidadDeVecesEnviado int       `json:"cantidad_de_veces_enviado"`
}
