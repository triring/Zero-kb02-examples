// tinygo build -o dispQRcode.uf2 --target waveshare-rp2040-zero -size=short .
// tinygo flash --target waveshare-rp2040-zero -size=short -monitor .

package main

import (
	"fmt"
	"image/color"
	"machine"
	"time"

	qrcode "github.com/skip2/go-qrcode"
	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinydraw"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

// カラーユニバーサルデザイン(CUD) カラーセット
var (
	// Accent Colors アクセントカラー
	red      = color.RGBA{R: 0xFF, G: 0x4B, B: 0x0, A: 0xFF}  //  Red : 赤
	yellow   = color.RGBA{R: 0xFF, G: 0xF1, B: 0x0, A: 0xFF}  //  Yellow : 黄色
	green    = color.RGBA{R: 0x3, G: 0xAF, B: 0x7A, A: 0xFF}  //  Green : 緑
	blue     = color.RGBA{R: 0x0, G: 0x5A, B: 0xFF, A: 0xFF}  //  Blue : 青
	sky_blue = color.RGBA{R: 0x4D, G: 0xC4, B: 0xFF, A: 0xFF} //  Sky blue : 空色
	pink     = color.RGBA{R: 0xFF, G: 0x80, B: 0x82, A: 0xFF} //  Pink : ピンク
	orange   = color.RGBA{R: 0xF6, G: 0xAA, B: 0x0, A: 0xFF}  //  Orange : オレンジ
	purple   = color.RGBA{R: 0x99, G: 0x0, B: 0x99, A: 0xFF}  //  Purple : 紫
	brown    = color.RGBA{R: 0x80, G: 0x40, B: 0x0, A: 0xFF}  //  Brown : 茶色

	// Base Colors  ベースカラー
	light_pink         = color.RGBA{R: 0xFF, G: 0xCA, B: 0xBF, A: 0xFF} //  Light pink : 明るいピンク
	cream              = color.RGBA{R: 0xFF, G: 0xFF, B: 0x80, A: 0xFF} //  Cream : クリーム
	light_yellow_green = color.RGBA{R: 0xD8, G: 0xF2, B: 0x55, A: 0xFF} //  Light yellow-green : 明るい黄緑
	light_sky_blue     = color.RGBA{R: 0xBF, G: 0xE4, B: 0xFF, A: 0xFF} //  Light sky blue : 明るい空色
	beige              = color.RGBA{R: 0xFF, G: 0xCA, B: 0x80, A: 0xFF} //  Beige : ベージュ
	light_green        = color.RGBA{R: 0x77, G: 0xD9, B: 0xA8, A: 0xFF} //  Light green : 明るい緑
	light_purple       = color.RGBA{R: 0xC9, G: 0xAC, B: 0xE6, A: 0xFF} //  Light purple : 明るい紫

	// Achromatic Colors 無彩色
	white      = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF} //  White  白
	light_gray = color.RGBA{R: 0xC8, G: 0xC8, B: 0xCB, A: 0xFF} //  Light gray  明るいグレー
	gray       = color.RGBA{R: 0x84, G: 0x91, B: 0x9E, A: 0xFF} //  Gray  グレー
	black      = color.RGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xFF}    //  Black  黒
)

/*
x,yを左上端として、与えられた文字列をQRコードにして描画する。
戻り値
成功	QRコードの描画サイズが、現在の液晶サイズで描画可能な場合は、QRコードを描画してかから、そのバージョン番号を返す。
失敗	QRコードの描画サイズが、現在の液晶サイズよりも大きい場合は、何も描画せず、0を返す。
*/
func dispQR(display ssd1306.Device, offset_x int, offset_y int, str string) int {
	// QRコードを作成（内容: strに定義した文字列、エラー訂正レベル: Low）
	// コントラストがはっきりした液晶にQRコードを表示するので、エラー訂正レベルはLowでも問題はない。
	qr, err := qrcode.New(str, qrcode.Low)
	if err != nil {
		fmt.Println("Error creating QR code:", err)
		return 0
	}
	//	var bitmap [][]bool
	// 2D bool配列を取得
	bitmap := qr.Bitmap()
	// 配列のサイズを表示
	// 表示する文字列の長さ、字種、エラー訂正レベルによってそのデータサイズは変動する。
	// 今回、使用したライブラリでは、与えられた文字列とエラー訂正レベルに合わせて、収納するVersionが選択され、QRコードの配列が生成される。
	// 生成されるは、QRコードの配列は、QRコード本体と上下左右に2ドットづつの余白が追加されたものとなる。
	// Verion 1 本体25x25 余白込み 29x29,表示可
	// Verion 2 本体29x29 余白込み 33x33,表示可
	// Verion 3 本体33x33 余白込み 37x37,表示不可 液晶内に収まらない。
	height := len(bitmap)
	width := len(bitmap[0]) // 全ての行が同じ幅
	fmt.Printf("QR Code bitmap size: %d x %d\n", height, width)
	fmt.Println("len:", len(str), str)
	time.Sleep(50 * time.Millisecond)
	if height > 33 {
		fmt.Println("The data size is too large.", height*2)
		return 0
	}
	// 例: 全配列を表示（true: 1, false: 0として）
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if bitmap[y][x] {
				tinydraw.FilledRectangle(&display, int16(x*2+offset_x), int16(y*2+offset_y), 2, 2, black)
				fmt.Print("1 ")
			} else {
				tinydraw.FilledRectangle(&display, int16(x*2+offset_x), int16(y*2+offset_y), 2, 2, white)
				fmt.Print("0 ")
			}
		}
		display.Display()
		fmt.Println()
	}
	// versionをチェックする。
	version := 0
	switch height {
	case 29: //	25x25
		version = 1
	case 33: //	29x29
		version = 2
	default:
		version = 0
	}
	return version
}

/*
x,yを左上端として、2次元配列に定義してあるマイクロQRコードのデータを描画する。
戻り値はなし。
*/
func dispMicroQR(display ssd1306.Device, offset_x int, offset_y int) {
	// var bitmap [][]int
	// 与えられた2D int配列のデータを使用する。
	// errorを表示するためのマイクロQRコード
	var mqr_err = [19][19]int{
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

	height := len(mqr_err)
	width := len(mqr_err[0]) // 全ての行が同じ幅
	fmt.Printf("micro QR Code bitmap size: %d x %d\n", height, width)
	time.Sleep(50 * time.Millisecond)
	// 例: 全配列を表示（intの１と0で表現されたデータ）
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if mqr_err[y][x] == 1 {
				tinydraw.FilledRectangle(&display, int16(x*2+offset_x), int16(y*2+offset_y), 2, 2, black)
				fmt.Print("1 ")
			} else {
				tinydraw.FilledRectangle(&display, int16(x*2+offset_x), int16(y*2+offset_y), 2, 2, white)
				fmt.Print("0 ")
			}
		}
		fmt.Println()
	}
	display.Display()
	return
}

func main() {
	// str := "イ" // 3 byte, size: 29 x 29 OK version 1(21x21)
	// str := "Hello tinygo!" // 13 byte, size: 29 x 29 OK version 1(21x21)
	// str := "TinyGo Keeb Tour" // 16 byte, size: 29 x 29 OK version 1(21x21)
	// str := "12345678901234567890123456789012" // 32 byte, size: 29 x 29 OK version 1(21x21)
	// str := "1234567890123456789012345678901234567890" // 40 byte, size: 29 x 29 OK version 1(21x21)
	// str := "https://tinygo.org/" // 19 byte, size: 33 x 33 OK version 2(25x25)
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
	str := "〒100-8111 東京都千代田区千代田１−１" // 51 byte,size: 37 x 37 NG version 3(29x29)

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
	display.SetRotation(drivers.Rotation180)
	display.ClearDisplay()
	fmt.Println("Display the QR code.")
	w, h := display.Size()
	fmt.Printf("display:%T, width:%d, height:%d\n", display, w, h)
	time.Sleep(50 * time.Millisecond)
	// tinydraw.FilledRectangle(&display, 0, 0, 64, 128, white)
	tinydraw.FilledRectangle(&display, 0, 0, 128, 64, white)
	ver := dispQR(display, 32, 0, str)
	display.Display()
	time.Sleep(100 * time.Millisecond)
	if ver == 0 {
		// version 3以上は、大きくて表示されないので、エラーを表示する。
		fmt.Println("Display a micro QR code indicating the error.")
		tinyfont.WriteLine(&display, &freemono.Regular9pt7b, 4, 36, "data is\ntoo large.\n", black)
		dispMicroQR(display, 88, 4)
	}
	for {
		time.Sleep(60 * time.Second)
	}
}
