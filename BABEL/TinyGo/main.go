// tinygo build -o BABEL.uf2 --target pico -size=short .
// tinygo flash --target pico -size=short -monitor .

package main

import (
	"fmt"
	"machine"
	"time"
	"crypto/subtle"
)

// TrimLastChar は文字列の最後のルーンを削除する。
func TrimLastChar(s string) string {
	if s == "" {
		return ""
	}
	// 文字列をルーンのスライスに変換する
	runes := []rune(s)
	// スライスの最後の要素を除外して、新しい文字列として返す
	return string(runes[:len(runes)-1])
}

// 文字列(str)を引数として受け取り、1文字つづ表示していく関数
func printString(str string, wait int) {
    runes := []rune(str)    // []runeにキャスト（文字ごとに分解）
	interval := time.Millisecond *  time.Duration(wait) // 1文字を表示する毎に、waitを定義
	// ループで1文字ずつ表示
	for _, r := range runes {
        fmt.Printf("%c", r)
		time.Sleep(interval)
	}
}

// getPassword はコンソールから文字列を受け取り、
// 不要な前後の空白を取り除いた文字列を返す。
func getPassword() string {
	enter_flag := false // 改行コードのチェック用フラグ
	var readbuffer string
	for { // キー入力待ち
		// PCからの受信データをチェック
		if machine.Serial.Buffered() > 0 {
			c, err := machine.Serial.ReadByte()
			if err == nil {
				if c < 32 {
					switch c {
					case '\r':
						enter_flag = true // machine.Serial.WriteByte('\r')
					case '\n':
						enter_flag = true // machine.Serial.WriteByte('\n')
					case '\b':
						if len(readbuffer) > 0 { // バックスペースで、最後尾の１文字を削除
							machine.Serial.WriteByte('\b') // 表示部分の最後の1文字を消去
							machine.Serial.WriteByte(' ')
							machine.Serial.WriteByte('\b')
							readbuffer = TrimLastChar(readbuffer) // すでに取り込んでいる文字列データの最後の1文字を消去
						}
					default:
						// println(c)	  -- >  BS=8, Enter=13
						// Convert nonprintable control characters to
						// ^A, ^B, etc.
						machine.Serial.WriteByte('^')
						machine.Serial.WriteByte(c + '@')
					}
				} else if c >= 127 {
					// Anything equal or above ASCII 127, print ^?.
					machine.Serial.WriteByte('^')
					machine.Serial.WriteByte('?')
				} else {
					// Echo the printable character back to the
					// host computer.
					machine.Serial.WriteByte(c)
					// 読み込んだ文字をエコーバックし、文字列バッファーに保存する。
					readbuffer = readbuffer + string(c)
				}
			}
		}
		// This assumes that the input is coming from a keyboard
		// so checking 120 times per second is sufficient. But if
		// the data comes from another processor, the port can
		// theoretically receive as much as 11000 bytes/second
		// (115200 baud). This delay can be removed and the
		// Serial.Read() method can be used to retrieve
		// multiple bytes from the receive buffer for each
		// iteration.
		if true == enter_flag {
			// 改行コードを検出したら、ループを抜け、取り込んだ文字列を返す。
			fmt.Printf("\n")
			break
		}
	}
	return readbuffer
}

func checkPassword(input, actual string) bool {
    // 文字列をバイトスライスに変換して比較
    return subtle.ConstantTimeCompare([]byte(input), []byte(actual)) == 1
}

func main() {
	// 文字列の定義
	str_password := "E.HOBA"
	str_Genesis_11_7 := `
 Go to, let us go down, and there confound
their language, that they may not understand
one another's speech.
`
	str_babel := " BABEL "
	// str_babel := " 𒁀𒀊𒅋𒌋  "  // アッカド語 では「神の門」を表す。 一方聖書によると、ヘブライ語の「balal（ごちゃ混ぜ、混乱）」から来ているとされる。
	time.Sleep(time.Millisecond * 3000)
	fmt.Printf("\033[2J\033[H")		// 画面クリア＋カーソルホーム（左上）へ移動
	fmt.Printf("\033[37m")  // 文字を白色に
	printString("attach cd 01 /\n", 50)
	time.Sleep(time.Millisecond * 5000)
	printString("enter author password\n", 50)

	//	入力待ち用のリーダー
	var flag_password bool = false
	for count := 0; count < 3; count++ {
		printString("pass:", 50)
		input := getPassword()
		flag_password = checkPassword(input, str_password)
		if true != flag_password {
			if 2 > count {
				printString("Sorry, try again.\n", 25)
			}
			continue
		}
		break
	}
	if true != flag_password {
		printString("3 incorrect password attempts\n", 25)
		// os.Exit(1)
		for {
			time.Sleep(time.Second * 30)
		}
	}

	time.Sleep(time.Millisecond *  3000)
	printString(str_Genesis_11_7, 25)
	time.Sleep(time.Millisecond *  5000)
	printString("\n\n\n\n\n", 1000)
	fmt.Printf("\033[31m ")  // 文字を赤色に
	var count int64 = 0	// int64 の最大値は, 9223372036854775807 
	for {
		printString(str_babel, 8)
		count++
		if count == 0 {
			break
		}
	}
	fmt.Printf("\033[37m")  // 文字を白色に
}
