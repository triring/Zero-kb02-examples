# Heatbeat Sound

## このプログラムについて
Raspberry pi のアクセスランプ（緑）は、ACT LEDは設定により外部LEDとしての機能を割り当てることができます。  
デフォルトでは、ACT LEDはSDへのアクセスで点滅するように設定されていますが、これをheartbeatモードに設定すると一定間隔で「短い点滅が2回、その後に少し長い休止」の点滅し続けるようになります。  

```bash
> sudo sh -c 'echo heartbeat  >/sys/class/leds/led0/trigger'
```

この点滅周期は、システムの負荷（ロードアベレージ）と連動しており、以下のように変動します。

* 低負荷: 点滅の間隔はゆっくりになる。
* 高負荷: 点滅の間隔が速くなり、最大で毎分180回程度の速い点滅になる。

今回は、これを模倣し、心拍を、フルカラーLEDで表現してみました。  
英語の医学的・科学的な文脈では、心音をLub-dubと表現します。  
これは、心臓の弁が閉じる音を正確に表したもので、「lub」が心室の収縮（第一心音）、「dub」が心室の弛緩（第二心音）に対応しています。  
今回は、Lub-dubを以下の音階で表現しています。  
* 5 octaves ラ->lub:心室の収縮(第一心音)
* 4 octaves ラ->dub:心室の弛緩(第二心音)

### コンパイル方法  

必要に応じて、以下のパッケージの導入して下さい。  

```bash
> go mod tidy
> go get github.com/tinygo-org/pio/rp2-pio
> go get github.com/tinygo-org/pio/rp2-pio/piolib
```

ソースコードは、[main.go](./main.go) です。    
このソースコードのあるディレクトリに移動して、以下のコマンドを実行して下さい。コンパイルが完了すると、生成した実行用バイナリがマイコンボードに転送されます。  

```bash
> tinygo flash --target waveshare-rp2040-zero --size short -monitor .
```

また、実行用バイナリを転送できない場合は、以下のコマンドで、実行用バイナリを作成し、手作業で、実行用バイナリをzero-kb02に転送して下さい。  

```bash
> tinygo build -o heartbeat_led.uf2 -target=waveshare-rp2040-zero -size short .
```

### 使い方

医療ドラマの病室等で鳴っているあの音です。  
変数OneBeatで、１拍の音の長さを決めているので、この部分の設定を変えると心拍数を上げたり下げたりできます。  
いろいろと試して見て下さい。  
