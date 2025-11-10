package game

import (
	"fmt"
	"geova-simulation/assets"
	"geova-simulation/simulation"
	"geova-simulation/state"
	"image"
	"image/color"
	"math"
	_"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// --- Constantes de Posición ---
// (Ajusta estos valores para que coincidan con tu diseño)
const (
	tripodeX = 150.0
	tripodeY = 200.0

	iconPythonX    = 300.0
	iconPythonY    = 200.0
	iconRabbitX    = 450.0
	iconRabbitY    = 200.0
	iconWebsocketX = 600.0
	iconWebsocketY = 200.0

	monitorX = 700.0
	monitorY = 350.0

	tiltMeterX = 100.0
	tiltMeterY = 300.0

	dashboardX = 50.0
	dashboardY = 400.0

	packetSpeed = 2.0 // Píxeles por frame
)

// Game implementa la interfaz ebiten.Game
type Game struct {
	Assets *assets.Assets
	State  *state.VisualState

	BotonRect      image.Rectangle
	isBotonPressed bool

	// Contadores para animaciones (sprite sheets)
	animPacketCounter int
	animIconCounter   int
}

// NewGame es el constructor de nuestro juego
func NewGame(assets *assets.Assets, state *state.VisualState, btnRect image.Rectangle) *Game {
	return &Game{
		Assets:    assets,
		State:     state,
		BotonRect: btnRect,
	}
}

// --- Lógica Principal (Update) ---

func (g *Game) Update() error {
	// Incrementa los contadores de animación (para los sprite sheets)
	// (El % 6 asume 6 frames por animación, ajústalo)
	g.animPacketCounter = (g.animPacketCounter + 1) % 360 // Un bucle largo
	g.animIconCounter = (g.animIconCounter + 1) % 360     // Un bucle largo

	// --- 1. Manejar Input del Usuario ---
	g.handleInput()

	// --- 2. Actualizar la Máquina de Estados (FSM) ---
	g.updatePacketFSM()

	return nil
}

func (g *Game) handleInput() {
	x, y := ebiten.CursorPosition()
	clickPoint := image.Pt(x, y)

	// Lógica visual del botón (presionado o no)
	g.isBotonPressed = g.BotonRect.Bounds().Canon().Overlaps(image.Rectangle{Min: clickPoint, Max: clickPoint.Add(image.Pt(1, 1))}) &&
		ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	// Lógica de inclinación (simulada)
	if ebiten.IsKeyPressed(ebiten.KeyLeft) && g.State.CurrentTilt > -15.0 {
		g.State.CurrentTilt -= 0.5
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) && g.State.CurrentTilt < 15.0 {
		g.State.CurrentTilt += 0.5
	}

	// --- ¡EL FAN-OUT! (Al hacer clic en el botón) ---
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && !g.State.SimulacionIniciada {
		if g.BotonRect.Bounds().Canon().Overlaps(image.Rectangle{Min: clickPoint, Max: clickPoint.Add(image.Pt(1, 1))}) {
			g.startSimulation()
		}
	}
}

func (g *Game) startSimulation() {
	// Resetea el estado
	g.State.Mutex.Lock()
	g.State.Packets = make(map[string]*state.PacketState)
	g.State.DisplayDistancia = 0
	g.State.DisplayNitidez = 0
	g.State.DisplayRoll = 0
	g.State.SimulacionIniciada = true // Bloquea el botón
	g.State.Mutex.Unlock()

	// ¡Lanza las 3 goroutines (los workers)!
	// (Usamos los datos de inclinación simulados)
	tilt := g.State.CurrentTilt

	go simulation.SendPOSTRequest(
		"http://localhost:8000/imx477/sensor",
		simulation.GenerateRandomIMXData(),
		"imx", g.State, 200.0, color.RGBA{G: 255, A: 255}, // Verde
	)
	go simulation.SendPOSTRequest(
		"http://localhost:8000/mpu/sensor",
		simulation.GenerateRandomMPUData(tilt), // Pasa la inclinación
		"mpu", g.State, 230.0, color.RGBA{B: 255, A: 255}, // Azul
	)
	go simulation.SendPOSTRequest(
		"http://localhost:8000/tfluna/sensor",
		simulation.GenerateRandomTFLunaData(),
		"tfluna", g.State, 260.0, color.RGBA{R: 255, A: 255}, // Rojo
	)
}

// updatePacketFSM es la Máquina de Estados Finita que mueve los paquetes
// Se ejecuta 60 veces por segundo en el hilo principal (no bloquea)
func (g *Game) updatePacketFSM() {
	g.State.Mutex.Lock()
	defer g.State.Mutex.Unlock()

	// Decrementa los timers de animación de los iconos
	if g.State.PythonAPITimer > 0 { g.State.PythonAPITimer-- }
	if g.State.RabbitMQTimer > 0 { g.State.RabbitMQTimer-- }
	if g.State.WebsocketAPITimer > 0 { g.State.WebsocketAPITimer-- }

	allDone := true
	for _, packet := range g.State.Packets {
		if packet.Status == state.Error || packet.Status == state.Done {
			continue // Este paquete ya terminó
		}

		allDone = false // Si al menos uno no ha terminado, la simulación no ha acabado

		// Mover el paquete hacia su objetivo
		if math.Abs(packet.X - packet.TargetX) > packetSpeed {
			if packet.X < packet.TargetX { packet.X += packetSpeed }
		} else {
			packet.X = packet.TargetX
		}
		// (Aquí iría la lógica para mover Y si no están alineados)

		// Comprobar si llegó al objetivo
		if packet.X == packet.TargetX {
			
			// Lógica de la FSM basada en el estado actual
			switch packet.Status {
			
			case state.ArrivedAtAPI:
				g.State.PythonAPITimer = 60 // Activa anim (1 segundo)
				packet.Status = state.SendingToRabbit
				packet.TargetX = iconRabbitX
				packet.TargetY = iconRabbitY

			case state.SendingToRabbit:
				if packet.X >= iconRabbitX { // Comprueba si "llegó"
					packet.Status = state.ArrivedAtRabbit
				}

			case state.ArrivedAtRabbit:
				g.State.RabbitMQTimer = 60
				packet.Status = state.SendingToWebsocket
				packet.TargetX = iconWebsocketX
				packet.TargetY = iconWebsocketY

			case state.SendingToWebsocket:
				if packet.X >= iconWebsocketX {
					packet.Status = state.ArrivedAtWebsocket
				}
			
			case state.ArrivedAtWebsocket:
				g.State.WebsocketAPITimer = 60
				packet.Status = state.SendingToFrontend
				packet.TargetX = monitorX
				packet.TargetY = monitorY
			
			case state.SendingToFrontend:
				if packet.X >= monitorX {
					packet.Status = state.Done
					packet.Active = false // Deja de dibujarlo
					
					// ¡ACTUALIZA EL DASHBOARD!
					// Extrae los datos del payload guardado
					switch data := packet.Payload.(type) {
					case simulation.TFLunaData:
						g.State.DisplayDistancia = data.DistanciaM
					case simulation.MPUData:
						g.State.DisplayRoll = data.Roll
					case simulation.IMXData:
						g.State.DisplayNitidez = data.Nitidez
					}
				}
			}
		}
	}
	
	if allDone {
		g.State.SimulacionIniciada = false // Reactiva el botón
	}
}


// --- Lógica de Dibujado (Draw) ---

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 0x1a, G: 0x1a, B: 0x1a, A: 255}) // Fondo oscuro

	// 1. Dibuja el Trípode
	opTripode := &ebiten.DrawImageOptions{}
	opTripode.GeoM.Translate(tripodeX, tripodeY)
	screen.DrawImage(g.Assets.GeovaTripod, opTripode)

	// 2. Dibuja el Medidor de Inclinación (Opción B)
	opTiltBG := &ebiten.DrawImageOptions{}
	opTiltBG.GeoM.Translate(tiltMeterX, tiltMeterY)
	screen.DrawImage(g.Assets.UITiltMeter, opTiltBG)
	// (Aquí iría la lógica para rotar una "aguja" del medidor
	// basada en g.State.CurrentTilt, si tu diseñador la hizo)

	// 3. Dibuja el Botón "Crear"
	opBoton := &ebiten.DrawImageOptions{}
	opBoton.GeoM.Translate(float64(g.BotonRect.Min.X), float64(g.BotonRect.Min.Y))
	if g.isBotonPressed {
		screen.DrawImage(g.Assets.ButtonCreateDown, opBoton)
	} else {
		screen.DrawImage(g.Assets.ButtonCreateUp, opBoton)
	}

	// 4. Dibuja los 4 Iconos del Flujo
	g.drawIcon(screen, g.Assets.IconPythonIdle, g.Assets.IconPythonActiveAnim, g.State.PythonAPITimer, iconPythonX, iconPythonY)
	g.drawIcon(screen, g.Assets.IconRabbitIdle, g.Assets.IconRabbitActiveAnim, g.State.RabbitMQTimer, iconRabbitX, iconRabbitY)
	g.drawIcon(screen, g.Assets.IconWebsocketIdle, g.Assets.IconWebsocketActiveAnim, g.State.WebsocketAPITimer, iconWebsocketX, iconWebsocketY)
	opMonitor := &ebiten.DrawImageOptions{}
	opMonitor.GeoM.Translate(monitorX, monitorY)
	screen.DrawImage(g.Assets.IconMonitor, opMonitor)
	
	// 5. Dibuja los Paquetes de Datos
	g.drawPackets(screen)
	
	// 6. Dibuja el Dashboard
	g.drawDashboard(screen)
}

// drawIcon es un helper para dibujar un icono (animado o estático)
func (g *Game) drawIcon(screen *ebiten.Image, idle *ebiten.Image, anim *ebiten.Image, timer int, x, y float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	
	if timer > 0 {
		// Dibuja la animación
		// (Asume 6 frames de 64x64, ajusta '64' y '6' si es diferente)
		frameWidth := 64
		frameCount := 6
		frameIndex := (g.animIconCounter / 6) % frameCount
		
		sx := frameIndex * frameWidth
		sy := 0
		rect := image.Rect(sx, sy, sx+frameWidth, sy+64) // Asume 64 de alto
		screen.DrawImage(anim.SubImage(rect).(*ebiten.Image), op)
	} else {
		// Dibuja el estático
		screen.DrawImage(idle, op)
	}
}

// drawPackets dibuja todos los paquetes de datos activos
func (g *Game) drawPackets(screen *ebiten.Image) {
	g.State.Mutex.Lock()
	defer g.State.Mutex.Unlock()

	// (Asume 6 frames de 32x32, ajusta '32' y '6' si es diferente)
	frameWidth := 32
	frameCount := 6
	frameIndex := (g.animPacketCounter / 6) % frameCount
	sx := frameIndex * frameWidth
	sy := 0
	rect := image.Rect(sx, sy, sx+frameWidth, sy+32) // Asume 32 de alto
	packetFrame := g.Assets.DataPacketAnim.SubImage(rect).(*ebiten.Image)

	for _, packet := range g.State.Packets {
		if !packet.Active { continue }
		
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(packet.X, packet.Y)
		
		// Tinta el paquete de su color
		op.ColorScale.SetR(float32(packet.Color.(color.RGBA).R) / 255)
		op.ColorScale.SetG(float32(packet.Color.(color.RGBA).G) / 255)
		op.ColorScale.SetB(float32(packet.Color.(color.RGBA).B) / 255)
		
		screen.DrawImage(packetFrame, op)

		if packet.Status == state.Error {
			ebitenutil.DebugPrintAt(screen, "X", int(packet.X)+10, int(packet.Y)-10)
		}
	}
}

// drawDashboard dibuja los medidores y barras de progreso
func (g *Game) drawDashboard(screen *ebiten.Image) {
	// Dibuja un título
	ebitenutil.DebugPrintAt(screen, "--- Dashboard de Resultados ---", int(dashboardX), int(dashboardY))

	// 1. Dibuja el Medidor de Distancia (TF-Luna)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Distancia: %.2f m", g.State.DisplayDistancia), int(dashboardX), int(dashboardY)+20)
	// (Aquí iría el código para dibujar el UIGaugeBG y rotar el UIGaugeNeedle)

	// 2. Dibuja la Barra de Nitidez (IMX477)
	ebitenutil.DebugPrintAt(screen, "Nitidez:", int(dashboardX), int(dashboardY)+40)
	opBarBG := &ebiten.DrawImageOptions{}
	opBarBG.GeoM.Translate(dashboardX+60, dashboardY+40)
	screen.DrawImage(g.Assets.UIProgressBG, opBarBG)

	opBarFill := &ebiten.DrawImageOptions{}
	// Escala la barra de relleno (g.State.DisplayNitidez es 4.0-6.0, lo normalizamos a 0.0-1.0)
	normalizedNitidez := (g.State.DisplayNitidez - 4.0) / 2.0
	if normalizedNitidez < 0 { normalizedNitidez = 0 }
	if normalizedNitidez > 1 { normalizedNitidez = 1 }
	
	opBarFill.GeoM.Scale(normalizedNitidez, 1.0) // ¡Escala en X!
	opBarFill.GeoM.Translate(dashboardX+60, dashboardY+40)
	screen.DrawImage(g.Assets.UIProgressFill, opBarFill)
	
	// 3. Dibuja el Medidor de Inclinación (MPU)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Inclinacion (Roll): %.1f°", g.State.DisplayRoll), int(dashboardX), int(dashboardY)+60)
}


// --- Layout (Define el tamaño de la ventana) ---

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 600
}