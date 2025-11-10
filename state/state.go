package state

import (
	"image/color"
	"sync"
)

// PacketStatus define los pasos de nuestro flujo visual
type PacketStatus int

const (
	Idle PacketStatus = iota // El paquete no existe
	SendingToAPI             // Moviéndose de Geova -> API Python
	ArrivedAtAPI             // ¡Éxito! El worker HTTP terminó
	SendingToRabbit          // Moviéndose de API Python -> RabbitMQ
	ArrivedAtRabbit          // ...
	SendingToWebsocket       // Moviéndose de RabbitMQ -> API Websocket
	ArrivedAtWebsocket       // ...
	SendingToFrontend        // Moviéndose de API Websocket -> Frontend
	Done                     // Llegó al frontend
	Error                    // La petición HTTP falló
)

// PacketState representa un paquete de datos en la pantalla
type PacketState struct {
	ID        string
	Active    bool
	X, Y      float64       // Posición actual
	TargetX, TargetY float64 // Hacia dónde se dirige
	Color     color.Color
	Status    PacketStatus
	Payload   interface{} // Guarda los datos aleatorios generados
}

// VisualState es el puente seguro (thread-safe)
type VisualState struct {
	Mutex   sync.Mutex
	Packets map[string]*PacketState // Un mapa para los 3 paquetes

	// Estado de los íconos (para animaciones)
	PythonAPITimer    int
	RabbitMQTimer     int
	WebsocketAPITimer int

	// Datos finales para el dashboard
	DisplayDistancia     float64
	DisplayRoll          float64
	DisplayNitidez       float64
	CurrentTilt          float64 // Para el medidor de inclinación
	SimulacionIniciada bool    // Para deshabilitar el botón
}