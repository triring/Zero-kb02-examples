/*
> tinygo flash -target=waveshare-rp2040-zero -size short -monitor .
> tinygo build -o heartbeat_sound.uf2 -target=waveshare-rp2040-zero -size short .
*/
/*
このプログラムについて
Raspberry pi のアクセスランプ（緑）は、ACT LEDは設定により外部LEDとしての機能を割り当てることができます。
デフォルトでは、ACT LEDはSDへのアクセスで点滅するように設定されていますが、これをheartbeatモードに設定すると一定間隔で「短い点滅が2回、その後に少し長い休止」の点滅し続けるようになります。
この点滅周期は、システムの負荷（ロードアベレージ）と連動しており、以下のように変動します。

* 低負荷: 点滅の間隔はゆっくりになる。
* 高負荷: 点滅の間隔が速くなり、最大で毎分180回程度の速い点滅になる。

今回は、これを模倣し、Beep音で、heartbeat(心音)を表現してみました。
英語の医学的・科学的な文脈では、heartbeat(心音)をLub-dubと表現されています。
これは、心臓の弁が閉じる音を正確に表したもので、「lub」が心室の収縮（第一心音）、「dub」が心室の弛緩（第二心音）に対応しています。
今回は、Lub-dubを以下の音階で表現しています。
* 5 octaves ラ->lub:心室の収縮(第一心音)
* 4 octaves ラ->dub:心室の弛緩(第二心音)
*/

/*
```bash
> go mod tidy
> go get tinygo.org/x/drivers/tone
```
*/
package main

import (
	"machine"
	"time"

	"tinygo.org/x/drivers/tone"
)

// zero-kb02用出力ポート等の定義
// Please connect a piezo buzzer to the 3V3 and EX01 pins on the back terminal.
// | EX01 | EX03 | 3V3 | SDA0 | 3V3 | 3V3 |     |        GROVE            |
// | EX02 | EX04 | GND | SCL0 | GND | GND | - - | GND | 3V3 | SDA0 | SCL0 |
// 今回は、圧電ブザーをGPIO15とGNDに接続した。

var pinToPWM = map[machine.Pin]tone.PWM{
	machine.GPIO14: machine.PWM7, // for EX01
	machine.GPIO15: machine.PWM7, // for EX02
	machine.GPIO26: machine.PWM5, // for EX03
	machine.GPIO27: machine.PWM5, // for EX04
}

/*
   beep音を出力する圧電ブザーを接続しているGPIOの設定について
   RP2040では8つのスライス、最大16チャンネルのPWMが扱える。
   今回は、GPIO15を圧電ブザーの出力に設定し、このGPIOに対応するPWMのチャンネルとしてPWM7を設定した。
   他のGPIOを使用する場合は、以下の表から、使用するGPIOに対応するPWM チャンネルを設定すること。

   GPIO		0	1	2	3	4	5	6	7	8	9	10	11	12	13	14	15
   PWM	Ch	0A	0B	1A	1B	2A	2B	3A	3B	4A	4B	5A	5B	6A	6B	7A	7B

   GPIO		16	17	18	19	20	21	22	23	24	25	26	27	28	29
   PWM	Ch	0A	0B	1A	1B	2A	2B	3A	3B	4A	4B	5A	5B	6A	6B
*/

// 無音の定義
var mute tone.Note = 0

// 心音のパターンデータ
var heartbeat_pattern = [8]tone.Note{
	tone.A5, // 5 octaves ラ->lub:心室の収縮(第一心音)
	mute,
	tone.A4, // 4 octaves ラ->dub:心室の弛緩(第二心音)
	mute,
	mute,
	mute,
	mute,
	mute,
}

func main() {
	// ブザーが接続されたPinの設定と初期化
	bzrPin := machine.GPIO15
	pwm := pinToPWM[bzrPin]
	speaker, err := tone.New(pwm, bzrPin)
	if err != nil {
		println("failed to configure PWM")
		return
	}

	// 一拍の長さ(us)
	// var OneBeat time.Duration = 100000
	var OneBeat time.Duration = 125000
	// var OneBeat time.Duration = 150000

	// var OneBeat int = 125
	// heartbeat
	for loop := 0; loop < 20; loop++ {
		for i := 0; i < len(heartbeat_pattern); i++ {
			speaker.SetNote(heartbeat_pattern[i])
			// time.Sleep(time.Millisecond * OneBeat)
			time.Sleep(time.Microsecond * OneBeat)
			speaker.SetNote(mute)
		}
	}
	speaker.SetNote(mute)
}
