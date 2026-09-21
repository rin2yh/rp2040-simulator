# Agent Rules

このルールに記載の内容は必ず守ること。

## 基本方針

- RP2040やGPIO / I2C / SPIをエミュレートしない。RPCはLED、Display、Encoderなどのデバイス単位にする。
- 利用者のアプリコードはPCとTinyGo実機で共通にする。
- 利用者のアプリをシミュレーター固有のloop、callback、`internal/` packageへ依存させない。
- `machine/`と`driver/`には、デスクトップとTinyGoでimportを変えずに使える最小限の互換packageだけを置く。
- 実装は`internal/`へ置き、ディレクトリは原則一段階までとする。
- 二つ以上の実装で必要になるまで抽象化を追加しない。
- READMEは概要、初期セットアップ、基本コマンドに留める。

## Toolchain and tasks

バージョンと定型コマンドは`mise.toml`で管理する。Makefileは追加しない。

- Setup: `mise install`。
- PC emulator: `mise run run`。
- テストと静的検査: `mise run test`。
- GUI描画テスト: `mise run test-gui`。
- スクリーンショット: `mise run screenshot`。
- PC build: `mise run build`。

TinyGoの検証は対象サンプルを直接ビルドする。実機書き込みはユーザーが明示的に依頼した場合だけ行う。

```sh
tinygo build -target=waveshare-rp2040-zero -o build/led-blink.uf2 ./examples/led-blink
```

## Structure

```text
machine/                 TinyGo machine互換の公開package
driver/ws2812/           PC RPC / TinyGo実ドライバの切り替え
examples/led-blink/      利用者コードの例
internal/bridge/         RPC protocolとclient
internal/board/          ボードprofileとGUI描画
internal/emulator/       Ebitengine、入力、仮想デバイス、RPC server
```

ボード追加時は`internal/board/<name>.go`へprofileを追加する。ボードごとのサブディレクトリは作らない。

## zero-kb02

- OLED: SSD1306 128x64、I2C0、`0x3c`、SDA GPIO12、SCL GPIO13、400 kHz、Rotation180。
- Encoder: GPIO3 / GPIO4、button GPIO2、precision 4、active-low。
- Keys: columns GPIO5..8、rows GPIO9..11、settle 1 ms。
- Joystick: X GPIO29、Y GPIO28、button GPIO0。ADC min `0x3000`、center `0x8000`、max `0xc800`。Yは反転。
- Key LEDs: WS2812 x12、GPIO1。

ADC中央ずれ、encoderの物理的な回転方向、OLED表示方向は実機未検証である。調整値はboard profileへ置く。

## Emulator behavior

- Keys: clickまたは`QWER` / `ASDF` / `ZXCV`。
- Encoder: circular/vertical drag、knob上のwheel、`[` / `]`。clickは短い押下、right click / Space / Enterは保持。
- Joystick: dragまたはarrow keys。releaseで中央。clickは短い押下、right click / Shiftは保持。
- RESET: click / F5で仮想デバイスを初期化する。
- BOOT: clickで保持を切り替え、RESETでsimulated bootloaderへ入る。`B` + F5も同じ。
- OLED: click / Oで拡大する。focus loss、modal、RESETでは入力を解除する。
- Key LEDs: RPCで受信した12個のRGB値を各キー中央へ描画する。

## Verification

通常は`mise run test`を実行する。GUI / input変更では`mise run test-gui`、公開driver変更では対応するTinyGoサンプルのビルドも実行する。

SSD1306 pixel parity、buffer commit、encoder fraction、pointer capture、focus loss、joystick clamp、RESET、RPC LED frame、外部moduleからのimportを維持する。

完成画面は`internal/board/testdata/<board>.png`をgolden imageとして比較する。意図したGUI変更では`UPDATE_BOARD_GOLDEN=1 mise run screenshot`を実行し、`build/emulator.png`を目視確認してからgolden imageを更新する。

## Release

- Git tagをGo moduleのバージョンとして扱い、ソースコード内にバージョン定数を置かない。
- tagprでSemantic VersioningのtagとGitHub Releaseを作成する。通常はpatch、`tagpr:minor`と`tagpr:major`で更新幅を指定する。
- 実行バイナリは配布しない。利用者はGo moduleとして取得する。
- CIはLinux上のGo test / vet、macOS上のGUI golden image、Windows build、TinyGoサンプルを検証する。
- 依存関係の更新にはDependabotを使う。コンテナを配布しないためTrivyは追加しない。
