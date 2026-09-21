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
- TinyGo互換性検証では`mise.toml`に指定されたfixtureを使い、`examples/`をfixtureにしない。

## Verification

検証コマンドは`mise.toml`を参照する。実機書き込みはユーザーが明示的に依頼した場合だけ行う。

外部moduleからのimportを維持する。

Git tagをGo moduleのバージョンとして扱い、ソースコード内にバージョン定数を置かない。実行バイナリやコンテナは配布しない。
