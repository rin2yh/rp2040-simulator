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

## Sources of truth

- ツールのバージョンと定型コマンドは`mise.toml`を参照する。Makefileは追加しない。
- CIとリリース設定は`.github/workflows/`、`.tagpr`、`.github/release.yml`、`.github/dependabot.yml`を参照する。
- packageとディレクトリの構成はソースツリーを参照する。
- ボードprofileは`internal/board/<name>.go`へ追加し、ボードごとのサブディレクトリは作らない。
- TinyGo互換性検証では`mise.toml`に指定されたfixtureを使い、`examples/`をfixtureにしない。

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

検証コマンドは`mise.toml`を参照する。実機書き込みはユーザーが明示的に依頼した場合だけ行う。

SSD1306 pixel parity、buffer commit、encoder fraction、pointer capture、focus loss、joystick clamp、RESET、RPC LED frame、外部moduleからのimportを維持する。

完成画面は`internal/board/testdata/<board>.png`をgolden imageとして比較する。意図したGUI変更では`UPDATE_BOARD_GOLDEN=1 mise run screenshot`を実行し、`build/emulator.png`を目視確認してからgolden imageを更新する。

Git tagをGo moduleのバージョンとして扱い、ソースコード内にバージョン定数を置かない。実行バイナリやコンテナは配布しない。
