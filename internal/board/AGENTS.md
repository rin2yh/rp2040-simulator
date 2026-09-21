# Board Rules

ボードprofileは`<name>.go`へ追加し、ボードごとのサブディレクトリは作らない。

## zero-kb02

- OLED: SSD1306 128x64、I2C0、`0x3c`、SDA GPIO12、SCL GPIO13、400 kHz、Rotation180。
- Encoder: GPIO3 / GPIO4、button GPIO2、precision 4、active-low。
- Keys: columns GPIO5..8、rows GPIO9..11、settle 1 ms。
- Joystick: X GPIO29、Y GPIO28、button GPIO0。ADC min `0x3000`、center `0x8000`、max `0xc800`。Yは反転。
- Key LEDs: WS2812 x12、GPIO1。

ADC中央ずれ、encoderの物理的な回転方向、OLED表示方向は実機未検証である。調整値はboard profileへ置く。
