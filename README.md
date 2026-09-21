# rp2040-simulator

RP2040ボード向けTinyGoアプリを、実機へ書き込まずにPCで動かすためのエミュレータです。zero-kb02とconf2025badgeに対応しています。デスクトップではローカルRPC、TinyGoでは実ドライバへ接続します。

![zero-kb02 emulator](internal/board/testdata/zero-kb02.png)

SSD1306 OLED、キーLED、ロータリーエンコーダ、ジョイスティック、BOOT / RESETを画面上で操作できます。RP2040やGPIO / I2C / SPI自体はエミュレートしません。

## Usage

Go moduleをアプリへ追加します。

```sh
go get github.com/rin2yh/rp2040-simulator@latest
```

PCとTinyGoで共通に使えるpackageと最小構成のアプリは、[LED blink example](examples/led-blink)を参照してください。

conf2025badgeは `RP2040_SIMULATOR_BOARD=conf2025badge go run .` で選択します。アプリからの自動起動にも同じ環境変数を使います。キーLEDは `machine.GPIO0` に2個接続し、TinyGoでは `-target=xiao-rp2040` を指定します。ブザーの音声出力は未対応です。

## Development

### Setup

[mise](https://mise.jdx.dev/) をインストールしてから、開発ツールをセットアップします。

```sh
mise install
```

### Commands

```sh
mise run run         # エミュレータを起動
mise run run-conf2025badge # conf2025badgeを起動
mise run check       # テストと静的検査
mise run screenshot  # GUIのスクリーンショットを生成
mise run build       # PC向けバイナリをビルド
```

## License

[MIT License](LICENSE)
