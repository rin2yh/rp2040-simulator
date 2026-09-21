# rp2040-simulator

RP2040ボード向けTinyGoアプリを、実機へ書き込まずにPCで動かすためのエミュレータです。zero-kb02とconf2025badgeに対応しています。デスクトップでは最初のデバイス操作時にエミュレータを自動起動してローカルRPCへ接続し、TinyGoでは実ドライバへ接続します。

![zero-kb02 emulator](internal/board/testdata/zero-kb02.png)

SSD1306 OLED、キーLED、ロータリーエンコーダ、ジョイスティック、BOOT / RESETを画面上で操作できます。RP2040やGPIO / I2C / SPI自体はエミュレートしません。

## Usage

Go moduleをアプリへ追加します。

```sh
go get github.com/rin2yh/rp2040-simulator/machine@latest
```

zero-kb02以外を使う場合は、デバイスを使う前に対象ボードを設定します。`Configure`を呼ばない場合はzero-kb02を使用します。

```go
if err := simulator.Configure(simulator.Config{
	Board: simulator.BoardConf2025Badge,
}); err != nil {
	log.Fatal(err)
}
```

PCとTinyGoで共通に使えるpackageと最小構成のアプリは、[LED blink example](examples/led-blink)を参照してください。

## Development

### Setup

[mise](https://mise.jdx.dev/) をインストールしてから、開発ツールをセットアップします。

```sh
mise install
```

### Commands

```sh
mise run check       # テストと静的検査
mise run screenshot  # GUIのスクリーンショットを生成
```

## License

[MIT License](LICENSE)
