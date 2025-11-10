package main

import (
	"geova-simulation/assets"     // ¡Ajusta si tu 'go.mod' es diferente!
	"geova-simulation/game"
	"geova-simulation/state"
	"image"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// --- Constantes Globales ---
const (
	windowWidth  = 800
	windowHeight = 600
)

func main() {
	// 1. Inicializa el generador de números aleatorios (¡Importante!)
	// (En Go 1.20+ esto ya no es necesario, pero no hace daño)
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// 2. Cargar todos los Assets
	// Llama a la función LoadAssets que definimos en el paquete 'assets'
	gameAssets := assets.LoadAssets()
	log.Println("✅ Todos los assets cargados.")

	// 3. Crear el Estado Compartido
	// Este es el objeto que las goroutines (workers) y la UI (game)
	// usarán para comunicarse.
	visualState := &state.VisualState{
		Packets:      make(map[string]*state.PacketState),
		CurrentTilt:  0.0, // Inclinación inicial
	}

	// 4. Crear la Instancia del Juego
	// Define la "zona de clic" para el botón de crear
	// (Ajusta estos números para mover tu botón)
	btnX0 := float64(windowWidth - 120) // Esquina derecha
	btnY0 := float64(windowHeight - 60) // Abajo
	btnRect := image.Rect(int(btnX0), int(btnY0), int(btnX0+100), int(btnY0+40)) // (100x40 de tamaño)

	juego := game.NewGame(gameAssets, visualState, btnRect)

	// 5. Configurar y Correr Ebitengine
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Simulación de Flujo Geova (Concurrente)")
	
	log.Println("🚀 Iniciando simulación...")
	
	// ebiten.RunGame toma control del hilo principal
	// y empezará a llamar a juego.Update() y juego.Draw()
	if err := ebiten.RunGame(juego); err != nil {
		log.Fatal(err)
	}
}