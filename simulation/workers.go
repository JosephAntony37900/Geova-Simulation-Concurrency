package simulation

import (
	"bytes"
	"encoding/json"
	"geova-simulation/state"
	"image/color"
	"math/rand"
	"net/http"
	"time"
)

// --- Generadores de Datos Aleatorios ---

func GenerateRandomIMXData() IMXData {
	// ... (código igual)
	return IMXData{
		IDProject:      1,
		Resolution:     "640x480",
		Luminosidad:    5.0 + rand.Float64()*10.0,
		Nitidez:        4.0 + rand.Float64()*2.0,
		LaserDetectado: rand.Intn(2) == 1,
		CalidadFrame:   20.0,
		Confiabilidad:  0.0,
		Event:          true,
		Timestamp:      time.Now().Format(time.RFC3339),
	}
}

// ¡FUNCIÓN CORREGIDA!
func GenerateRandomMPUData(tilt float64) MPUData { // <-- 1. Acepta 'tilt'
	return MPUData{
		IDProject: 1,
		Ax:        0.1 + rand.Float64()*0.1,
		Ay:        -0.05 + rand.Float64()*0.1,
		Az:        9.8 + rand.Float64()*0.1,
		Gx:        0.01 + rand.Float64()*0.02,
		Gy:        0.02 + rand.Float64()*0.02,
		Gz:        0.03 + rand.Float64()*0.02,
		Roll:      tilt, // <-- 2. Usa el 'tilt' del juego
		Pitch:     0.5 + rand.Float64()*1.0,
		Apertura:  tilt * 1.5,
		Event:     true,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// (GenerateRandomTFLunaData se queda igual)
func GenerateRandomTFLunaData() TFLunaData {
	// ... (código igual)
	distCm := 200 + rand.Intn(100)
	return TFLunaData{
		IDProject:   1,
		DistanciaCm: distCm,
		DistanciaM:  float64(distCm) / 100.0,
		FuerzaSenal: 5000 + rand.Intn(1000),
		Temperatura: 50.0 + rand.Float64()*5.0,
		Event:       true,
		Timestamp:   time.Now().Format(time.RFC3339),
	}
}


// --- El Worker Concurrente ---
// (Asegúrate de que la importación esté corregida aquí también)
func SendPOSTRequest(url string, payload interface{}, packetID string, visState *state.VisualState, startY float64, c color.Color) {
	
	// 1. Crear el paquete y actualizar el estado a "enviando"
	visState.Mutex.Lock()
	packet := &state.PacketState{
		ID:       packetID,
		Active:   true,
		X:        150.0, 
		Y:        startY,
		TargetX:  300.0, 
		TargetY:  200.0,
		Color:    c,
		Status:   state.SendingToAPI,
		Payload:  payload,
	}
	visState.Packets[packetID] = packet
	visState.Mutex.Unlock()

	// ... (El resto de la función: Marshal, Sleep, http.Post, etc... se queda igual) ...
	jsonData, err := json.Marshal(payload)
	if err != nil {
		visState.Mutex.Lock()
		visState.Packets[packetID].Status = state.Error
		visState.Mutex.Unlock()
		return
	}

	time.Sleep(time.Duration(500+rand.Intn(500)) * time.Millisecond)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))

	visState.Mutex.Lock()
	if err != nil {
		visState.Packets[packetID].Status = state.Error
	} else if resp.StatusCode >= 400 {
		visState.Packets[packetID].Status = state.Error
	} else {
		visState.Packets[packetID].Status = state.ArrivedAtAPI
	}
	visState.Mutex.Unlock()
}