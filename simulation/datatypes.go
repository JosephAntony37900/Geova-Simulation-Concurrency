package simulation

// JSON para IMX477 (http://localhost:8000/imx477/sensor)
type IMXData struct {
	IDProject      int     `json:"id_project"`
	Resolution     string  `json:"resolution"`
	Luminosidad    float64 `json:"luminosidad_promedio"`
	Nitidez        float64 `json:"nitidez_score"`
	LaserDetectado bool    `json:"laser_detectado"`
	CalidadFrame   float64 `json:"calidad_frame"`
	Confiabilidad  float64 `json:"probabilidad_confiabilidad"`
	Event          bool    `json:"event"`
	Timestamp      string  `json:"timestamp"`
}

// JSON para MPU (http://localhost:8000/mpu/sensor)
type MPUData struct {
	IDProject int     `json:"id_project"`
	Ax        float64 `json:"ax"`
	Ay        float64 `json:"ay"`
	Az        float64 `json:"az"`
	Gx        float64 `json:"gx"`
	Gy        float64 `json:"gy"`
	Gz        float64 `json:"gz"`
	Roll      float64 `json:"roll"`
	Pitch     float64 `json:"pitch"`
	Apertura  float64 `json:"apertura"`
	Event     bool    `json:"event"`
	Timestamp string  `json:"timestamp"`
}

// JSON para TFLuna (http://localhost:8000/tfluna/sensor)
type TFLunaData struct {
	IDProject   int     `json:"id_project"`
	DistanciaCm int     `json:"distancia_cm"`
	DistanciaM  float64 `json:"distancia_m"`
	FuerzaSenal int     `json:"fuerza_senal"`
	Temperatura float64 `json:"temperatura"`
	Event       bool    `json:"event"`
	Timestamp   string  `json:"timestamp"`
}