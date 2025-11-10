package assets

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Assets almacena todos los sprites cargados en memoria.
type Assets struct {
	// Hardware
	GeovaTripod *ebiten.Image
	UITiltMeter *ebiten.Image // (El medidor de inclinación)

	// Iconos de Backend (Inactivos)
	IconPythonIdle    *ebiten.Image
	IconRabbitIdle    *ebiten.Image
	IconWebsocketIdle *ebiten.Image

	// Iconos de Backend (Animados)
	IconPythonActiveAnim    *ebiten.Image
	IconRabbitActiveAnim    *ebiten.Image
	IconWebsocketActiveAnim *ebiten.Image

	// Paquete de Datos (Animado)
	DataPacketAnim *ebiten.Image

	// Frontend
	IconMonitor    *ebiten.Image
	UIGaugeBG      *ebiten.Image
	UIGaugeNeedle  *ebiten.Image
	UIProgressBG   *ebiten.Image
	UIProgressFill *ebiten.Image

	// Botones
	ButtonCreateUp   *ebiten.Image
	ButtonCreateDown *ebiten.Image
}

// loadSprite es un helper interno para cargar una imagen o fallar.
func loadSprite(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatalf("Error: No se pudo cargar el asset '%s': %v", path, err)
	}
	return img
}

// LoadAssets carga todas las imágenes del juego desde la carpeta /images
// ¡Asegúrate de que los nombres de archivo coincidan!
func LoadAssets() *Assets {
	return &Assets{
		// Hardware
		GeovaTripod: loadSprite("images/geova_tripod.png"),
		UITiltMeter: loadSprite("images/geova_tilt_anim.png"),

		// Iconos Inactivos
		IconPythonIdle:    loadSprite("images/icon_api_python_idle.png"),
		IconRabbitIdle:    loadSprite("images/icon_rabbitmq_idle.png"),
		IconWebsocketIdle: loadSprite("images/icon_api_websocket_idle.png"),

		// Iconos Animados (Sprite Sheets)
		IconPythonActiveAnim:    loadSprite("images/icon_api_python_active_anim.png"),
		IconRabbitActiveAnim:    loadSprite("images/icon_rabbitmq_active_anim.png"),
		IconWebsocketActiveAnim: loadSprite("images/icon_api_websocket_active_anim.png"),

		// Paquete de Datos (Sprite Sheet)
		DataPacketAnim: loadSprite("images/data_packet_anim.png"),

		// Frontend
		IconMonitor:    loadSprite("images/frontend_monitor.png"),
		UIGaugeBG:      loadSprite("images/ui_gauge_background.png"),
		UIGaugeNeedle:  loadSprite("images/ui_gauge_needle.png"),
		UIProgressBG:   loadSprite("images/ui_progressbar_background.png"),
		UIProgressFill: loadSprite("images/ui_progressbar_fill.png"),

		// Botones
		ButtonCreateUp:   loadSprite("images/boton_crear_up.png"),
		ButtonCreateDown: loadSprite("images/boton_crear_down.png"),
	}
}