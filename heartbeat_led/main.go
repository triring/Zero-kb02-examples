/*
> tinygo flash -target=waveshare-rp2040-zero -size short -monitor .
> tinygo build -o heartbeat_led.uf2 -target=waveshare-rp2040-zero -size short .
*/
/*
このプログラムについて
Raspberry pi のアクセスランプ（緑）は、ACT LEDは設定により外部LEDとしての機能を割り当てることができます。
デフォルトでは、ACT LEDはSDへのアクセスで点滅するように設定されているが、これをheartbeatモードに設定すると一定間隔で「短い点滅が2回、その後に少し長い休止」の点滅し続けるようになります。
この点滅周期は、システムの負荷（ロードアベレージ）と連動しており、以下のように変動する。

* 低負荷: 点滅の間隔はゆっくりになる。
* 高負荷: 点滅の間隔が速くなり、最大で毎分180回程度の速い点滅になる。

今回は、これを模倣し、フルカラーLEDで、心拍を表現してみた。
英語の医学的・科学的な文脈では、心音をLub-dubと表現する。
これは、心臓の弁が閉じる音を正確に表したもので、「lub」が心室の収縮（第一心音）、「dub」が心室の弛緩（第二心音）に対応している。
今回は、以下の配色でLub-dubと表現する。
* 赤色->lub:心室の収縮(第一心音)
* 青色->dub:心室の弛緩(第二心音)
*/

package main

import (
	"image/color"
	"machine"
	"time"

	pio "github.com/tinygo-org/pio/rp2-pio"
	"github.com/tinygo-org/pio/rp2-pio/piolib"
)

type WS2812B struct {
	Pin machine.Pin
	ws  *piolib.WS2812B
}

func NewWS2812B(pin machine.Pin) *WS2812B {
	s, _ := pio.PIO0.ClaimStateMachine()
	ws, _ := piolib.NewWS2812B(s, pin)
	return &WS2812B{
		ws: ws,
	}
}

func (ws *WS2812B) PutColor(c color.Color) {
	ws.ws.PutColor(c)
}

// 色の定義
var (
	red   = color.RGBA{R: 0xFF, G: 0x00, B: 0x0, A: 0xFF}  //  Red : 赤
	blue  = color.RGBA{R: 0x0, G: 0x5A, B: 0xFF, A: 0xFF}  //  Blue : 青
	white = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF} //  White  白
	black = color.RGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xFF}    //  Black  黒
)

// LEDの発光パターン
var lighting_pattern = [8]color.Color{
	red,  // 赤色->lub:心室の収縮(第一心音)
	blue, // 青色->dub:心室の弛緩(第二心音)
	black,
	black,
	black,
	black,
	black,
	black,
}

func main() {
	ws := NewWS2812B(machine.GPIO16)
	ws.PutColor(black)
	time.Sleep(time.Millisecond * 500)
	/*
		for {
			time.Sleep(time.Millisecond * 500)
			ws.PutColor(black)
			time.Sleep(time.Millisecond * 500)
			ws.PutColor(white)
		}
	*/
	for {
		for i := 0; i < len(lighting_pattern); i++ {
			ws.PutColor(lighting_pattern[i])
			time.Sleep(time.Millisecond * 25)
			ws.PutColor(black)
			time.Sleep(time.Millisecond * 100)
		}
	}
}
