// tinygo build -o dispQRcode.uf2 --target waveshare-rp2040-zero -size=short .
// tinygo flash --target waveshare-rp2040-zero -size=short -monitor .

package main

import (
	"fmt"
	"image/color"
	"machine"
	"math/rand"
	"strconv"
	"time"

	qrcode "github.com/skip2/go-qrcode"
	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinydraw"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

// 京都 [伏見稲荷大社](https://inari.jp/)の神籤
var oracle = [32][3]string{
	{"　　一番", "吉凶未分末大吉", "よしあしいまだわからずすえだいきち"},
	{"　　二番", "大大吉", "だいだいきち"},
	{"　　三番", "大吉", "だいきち"},
	{"　　四番", "凶後吉", "きょうのちきち"},
	{"　　五番", "大吉", "だいきち"},
	{"　　六番", "末吉", "すえきち"},
	{"　　七番", "凶後大吉", "きょうのちだいきち"},
	{"　　八番", "吉凶相交末吉", "きちきょうあいまじわりすえきち"},
	{"　　九番", "末大吉", "すえだいきち"},
	{"　　十番", "吉凶相半", "きちきょうあいなかばす"},
	{"　十一番", "末大吉", "すえだいきち"},
	{"　十二番", "凶後吉", "きょうのちきち"},
	{"　十三番", "大吉", "だいきち"},
	{"　十四番", "小凶後吉", "しょうきょうのちきち"},
	{"　十五番", "中吉", "ちゅうきち"},
	{"　十六番", "小吉", "しょうきち"},
	{"　十七番", "向大吉", "むかうだいきち"},
	{"　十八番", "吉", "きち"},
	{"　十九番", "末吉", "すえきち"},
	{"　二十番", "後吉", "のちきち"},
	{"二十一番", "大吉", "だいきち"},
	{"二十二番", "末大吉", "すえだいきち"},
	{"二十三番", "凶後大吉", "きょうのちだいきち"},
	{"二十四番", "吉凶不分末吉", "きちきょうわかたずすえきち"},
	{"二十五番", "凶後吉", "きょうのちきち"},
	{"二十六番", "凶後吉", "きょうのちきち"},
	{"二十七番", "吉凶相央", "きちきょうあいなかばす"},
	{"二十八番", "大吉", "だいきち"},
	{"二十九番", "向大吉", "むかうだいきち"},
	{"　三十番", "大吉", "だいきち"},
	{"三十一番", "吉", "きち"},
	{"三十二番", "大大吉", "だいだいきち"},
}

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
	// qr, err := qrcode.New(str, qrcode.Low)
	qr, err := qrcode.NewWithForcedVersion(str, 2, qrcode.Low)
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
	/*
		fmt.Printf("QR Code bitmap size: %d x %d\n", height, width)
		fmt.Println("len:", len(str), str)
		time.Sleep(50 * time.Millisecond)
		if height > 33 {
			fmt.Println("The data size is too large.", height*2)
			return 0
		}
	*/
	// 例: 全配列を表示（true: 1, false: 0として）
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if bitmap[y][x] {
				tinydraw.FilledRectangle(&display, int16(x*2+offset_x), int16(y*2+offset_y), 2, 2, black)
				// fmt.Print("1 ")
			} else {
				tinydraw.FilledRectangle(&display, int16(x*2+offset_x), int16(y*2+offset_y), 2, 2, white)
				// fmt.Print("0 ")
			}
		}
		// fmt.Println()
	}
	display.Display()

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

func dispMsg24(display ssd1306.Device, x int, y int, str string) {
	//	display.ClearBuffer()
	tinyfont.WriteLine(&display, &freemono.Bold24pt7b, int16(x), int16(y), str, black)
	display.Display()
}

func dispMsg18(display ssd1306.Device, x int, y int, str string) {
	//	display.ClearBuffer()
	tinyfont.WriteLine(&display, &freemono.Bold18pt7b, int16(x), int16(y), str, black)
	display.Display()
}

func dispMsg12(display ssd1306.Device, x int, y int, str string) {
	//	display.ClearBuffer()
	tinyfont.WriteLine(&display, &freemono.Bold12pt7b, int16(x), int16(y), str, black)
	display.Display()
}

func dispMsg09(display ssd1306.Device, x int, y int, str string) {
	//	display.ClearBuffer()
	tinyfont.WriteLine(&display, &freemono.Bold9pt7b, int16(x), int16(y), str, black)
	display.Display()
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
	display.SetRotation(drivers.Rotation180)
	display.ClearDisplay()
	fmt.Println("QR code Omikuji.")
	w, h := display.Size()
	fmt.Printf("display:%T, width:%d, height:%d\n", display, w, h)
	time.Sleep(50 * time.Millisecond)

	tinydraw.FilledRectangle(display, 0, 0, 128, 64, white)
	dispMsg24(*display, 12, 32, "QR")
	dispMsg12(*display, 12, 54, "omikuji")
	time.Sleep(5000 * time.Millisecond)

	tinydraw.FilledRectangle(display, 0, 0, 128, 64, white)

	for t := 5; t >= 0; t-- {
		tinydraw.FilledRectangle(display, 0, 0, 128, 64, white)
		dispMsg24(*display, 48, 48, strconv.Itoa(t))
		time.Sleep(1000 * time.Millisecond)
	}

	// tinydraw.FilledRectangle(&display, 0, 0, 64, 128, white)
	tinydraw.FilledRectangle(display, 0, 0, 128, 64, white)

	rand.Seed(time.Now().UnixNano()) // 現在の時刻をシードとして設定
	for i := 0; i < 32; i++ {
		randomInt := rand.Intn(32) // 0から31までの整数を生成
		ver := dispQR(*display, 32, 0, oracle[randomInt][1])
		fmt.Printf("%s\t%s\t%s\n",
			oracle[randomInt][0], oracle[randomInt][1], oracle[randomInt][2])
		if ver != 2 {
			fmt.Println(ver)
		}
		time.Sleep(500 * time.Millisecond)
	}
	for {
		tinydraw.FilledRectangle(display, 0, 0, 128, 64, white)
		dispQR(*display, 0, 0, "Press reset button and restart.")
		time.Sleep(1000 * time.Millisecond)
		tinydraw.FilledRectangle(display, 0, 0, 128, 64, white)
		dispQR(*display, 64, 0, "初期化釦を押し再起動")
		time.Sleep(1000 * time.Millisecond)
	}
}
