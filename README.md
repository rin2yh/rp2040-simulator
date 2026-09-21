# rp2040-simulator

TinyGoアプリを実機へ書き込まずにPCで動かすためのzero-kb02エミュレータです。デスクトップではローカルRPC、TinyGoでは実ドライバへ接続します。

![zero-kb02 emulator](internal/board/testdata/zero-kb02.png)

SSD1306 OLED、キーLED、ロータリーエンコーダ、ジョイスティック、BOOT / RESETを画面上で操作できます。RP2040やGPIO / I2C / SPI自体はエミュレートしません。

## Setup

[mise](https://mise.jdx.dev/) をインストールしてから、ツールをセットアップします。

```sh
mise install
```

## Run

```sh
go run ./examples/led-blink
```

PCではアプリの起動時にエミュレータも開きます。画面上の部品をクリックまたはドラッグして操作できます。

別のGo moduleから利用する場合は依存関係を追加し、[examples/led-blink](examples/led-blink)と同じpackageをimportします。

```sh
go get github.com/rin2yh/rp2040-simulator@v0.1.0
```

## LED blink on zero-kb02

```sh
tinygo build -target=waveshare-rp2040-zero -o build/led-blink.uf2 ./examples/led-blink

# 実機へ書き込む場合
tinygo flash -target=waveshare-rp2040-zero ./examples/led-blink
```

## Development

```sh
mise run test
```

## License

[MIT License](LICENSE)
