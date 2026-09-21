# Emulator Rules

## Behavior

- Keys: clickまたは`QWER` / `ASDF` / `ZXCV`。
- Encoder: circular/vertical drag、knob上のwheel、`[` / `]`。clickは短い押下、right click / Space / Enterは保持。
- Joystick: dragまたはarrow keys。releaseで中央。clickは短い押下、right click / Shiftは保持。
- RESET: click / F5で仮想デバイスを初期化する。
- BOOT: clickで保持を切り替え、RESETでsimulated bootloaderへ入る。`B` + F5も同じ。
- OLED: click / Oで拡大する。focus loss、modal、RESETでは入力を解除する。
- Key LEDs: RPCで受信した12個のRGB値を各キー中央へ描画する。

## Verification

SSD1306 pixel parity、buffer commit、encoder fraction、pointer capture、focus loss、joystick clamp、RESET、RPC LED frameを維持する。
