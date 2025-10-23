// tinygo build -o dispQRcode.uf2 --target waveshare-rp2040-zero -size=short .
// tinygo flash --target waveshare-rp2040-zero -size=short -monitor .

package main

import (
	"fmt"
	"image/color"
	"machine"
	"time"

	qrcode "github.com/skip2/go-qrcode"
	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinydraw"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

// errorを表示するためのマイクロQRコード
var qr_err = [19][19]int{
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 0},
	{0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0},
	{0, 0, 1, 0, 1, 1, 1, 0, 1, 0, 1, 1, 1, 1, 1, 0, 1, 0, 0},
	{0, 0, 1, 0, 1, 1, 1, 0, 1, 0, 1, 0, 1, 0, 0, 1, 1, 0, 0},
	{0, 0, 1, 0, 1, 1, 1, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
	{0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 1, 0, 0},
	{0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 0, 1, 0, 0, 0, 0},
	{0, 0, 0, 1, 1, 1, 1, 0, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0},
	{0, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 0, 0, 0},
	{0, 0, 0, 1, 0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0},
	{0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 1, 0, 0},
	{0, 0, 0, 0, 0, 1, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 1, 0, 0},
	{0, 0, 1, 0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
}

func main() {
	// I2Cの初期設定
	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: 2.8 * machine.MHz,
		SDA:       machine.GPIO12,
		SCL:       machine.GPIO13,
	})

	// OLEDディスプレイの初期設定
	display := ssd1306.NewI2C(machine.I2C0)
	display.Configure(ssd1306.Config{
		Address: 0x3C,
		Width:   128,
		Height:  64,
	})

	time.Sleep(500 * time.Millisecond)
	fmt.Println("Display the QR code.\n")

	black := color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
	white := color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	tinydraw.FilledRectangle(&display, 0, 0, 64, 128, white)
	display.Display()

	// str := "イ" // 3 byte, size: 29 x 29 OK version 1(21x21)
	// str := "Hello tinygo!" // 13 byte, size: 29 x 29 OK version 1(21x21)
	// str := "TinyGo Keeb Tour" // 16 byte, size: 29 x 29 OK version 1(21x21)
	// str := "12345678901234567890123456789012" // 32 byte, size: 29 x 29 OK version 1(21x21)
	// str := "1234567890123456789012345678901234567890" // 40 byte, size: 29 x 29 OK version 1(21x21)
	str := "https://tinygo.org/" // 19 byte, size: 33 x 33 OK version 2(25x25)
	// str := "Go言語はシンプル" // 23 byte, size: 33 x 33 OK version 2(25x25)
	// str := "https://go.dev/blog/gopher" // 26 byte, size: 33 x 33 OK version 2(25x25)
	// str := "The Go gopher was born in 2009." // 31 byte, size: 33 x 33 OK version 2(25x25)
	// str := "天上天下唯我独尊" // 24 byte, size: 33 x 33 OK version 2(25x25)
	// str := "臨兵闘者皆陣烈在前" // 27 byte, size: 33 x 33 OK version 2(25x25)
	// str := "甲乙丙丁戊己庚辛壬癸" // 30 byte, size: 33 x 33 OK version 2(25x25)
	// str := "摩訶般若波羅蜜多心経" // 30 byte, size: 33 x 33 OK version 2(25x25)
	// str := "銃砲刀剣類所持等取締法" // 33 byte, size: 37 x 37 NG version 3(33x33)
	// str := "日本餃子焼売食品工業会" // 33 byte, size: 37 x 37 NG version 3(33x33)
	// str := "子丑寅卯辰巳午未申酉戌亥" // 36 byte, size: 37 x 37 NG version 3(33x33)
	// str := "The mascot of the golang is Gopher" // 34 byte, size: 37 x 37 NG version 3(29x29)
	// str := "https://tinygo.org/docs/reference/" // 34 byte, size: 37 x 37 NG version 3(29x29)
	// str := "https://tinygo.org/getting-started/" // 35 byte, size: 37 x 37 NG version 3(29x29)
	// str := "Jackdaws love my big sphinx of quartz." // 38 byte, size: 37 x 37 NG version 3(29x29)
	// str := "The Go gopher was designed by Renee French." // 43 byte,size: 37 x 37 NG version 3(29x29)
	// str := "https://tinygo.org/getting-started/install/" // 43 byte, size: 37 x 37 NG version 3(29x29)
	// str := "GoのマスコットキャラクターGopher" // 44 byte, size: 37 x 37 NG version 3(29x29)
	// str := "The quick brown fox jumps over the lazy dog." // 48 byte,size: 37 x 37 NG version 3(29x29)
	// str := "〒100-8111 東京都千代田区千代田１−１" // 51 byte,size: 37 x 37 NG version 3(29x29)

	// QRコードを作成（内容: strに定義した文字列、エラー訂正レベル: Low）
	qr, err := qrcode.New(str, qrcode.Low)
	if err != nil {
		fmt.Println("Error creating QR code:", err)
		return
	}
	time.Sleep(500 * time.Millisecond)

	//	var bitmap [][]bool
	// 2D bool配列を取得
	bitmap := qr.Bitmap()
	// 配列のサイズを表示（例: 25x25など、内容とレベルによって変動）
	height := len(bitmap)
	width := len(bitmap[0]) // 全ての行が同じ幅
	tinydraw.FilledRectangle(&display, 0, 0, 128, 64, white)
	if height > 33 {
		// version 3以上は、大きくて表示されないので、エラーを表示する。
		tinyfont.WriteLine(&display, &freemono.Regular9pt7b, 0, 40, "data is\ntoo large.\n", black)
		for y := 0; y < len(qr_err); y++ {
			for x := 0; x < len(qr_err[0]); x++ {
				if qr_err[y][x] == 1 {
					tinydraw.FilledRectangle(&display, int16(x*2)+88, int16(y*2)+4, 2, 2, black)
					fmt.Print("1 ")
				} else {
					tinydraw.FilledRectangle(&display, int16(x*2)+88, int16(y*2)+4, 2, 2, white)
					fmt.Print("0 ")
				}
			}
			display.Display()
			fmt.Println()
		}
		fmt.Printf("The data size is too large.\n")
	} else {
		// 例: 全配列を表示（true: 1, false: 0として）
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				if bitmap[y][x] {
					tinydraw.FilledRectangle(&display, int16(x*2), int16(y*2), 2, 2, black)
					fmt.Print("1 ")
				} else {
					tinydraw.FilledRectangle(&display, int16(x*2), int16(y*2), 2, 2, white)
					fmt.Print("0 ")
				}
			}
			display.Display()
			fmt.Println()
		}
	}
	fmt.Printf("QR Code bitmap size: %d x %d\n", height, width)
	fmt.Println("len:", len(str), str)
	display.Display()
	for {
		time.Sleep(60 * time.Second)
	}
}
