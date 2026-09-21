# LED blink on zero-kb02

zero-kb02の12個のキーLEDを点滅させる最小構成のアプリです。同じコードをPCエミュレータとTinyGo実機で使用できます。

## PC

リポジトリのルートから実行します。アプリの起動時にエミュレータも開きます。

```sh
go run ./examples/led-blink
```

## TinyGo

```sh
tinygo build -target=waveshare-rp2040-zero -o build/led-blink.uf2 ./examples/led-blink
```

実機へ書き込む場合だけ、次のコマンドを実行します。

```sh
tinygo flash -target=waveshare-rp2040-zero ./examples/led-blink
```
