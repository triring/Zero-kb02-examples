// go run .\main.go
// go build -o BABEL.exe .\main.go

package main

import (
	"bufio"
	"fmt"
	"time"
	"os"	
	"strings"
	"crypto/subtle"
)

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
func getPassword() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		// 入力された文字列の前後空白を取り除く
		return strings.TrimSpace(scanner.Text()), nil
	}
	return "", scanner.Err()
}

// 入力されたPasswordをチェックする関数
func checkPassword(input, actual string) bool {
    // 文字列をバイトスライスに変換して比較
	if subtle.ConstantTimeCompare([]byte(input), []byte(actual)) == 1 {
		return true
	} else {
		return false
	}
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
	fmt.Printf("\033[37m")  // 文字を白色に
	fmt.Printf("\033[2J\033[H")		// 画面クリア＋カーソルホーム（左上）へ移動

	printString("attach cd 01 /\n", 50)
	time.Sleep(time.Millisecond * 5000)
	printString("enter author password\n", 50)

	//	入力待ち用のリーダー
	var flag_password bool = false
	for count := 0; count < 3; count++ {
		printString("pass:", 50)
		input,_  := getPassword()
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
		os.Exit(1)
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
