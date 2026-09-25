# zero-kb02 device check

Run `go run ./examples/zero-kb02` to open the desktop emulator. QWER/ASDF/ZXCV light the corresponding keys and OLED blocks. Rotate the knob to move the center pixel; move the joystick to move the lower pixel. Press the knob or stick to light a bottom corner.

The same application API is compiled for the real board by `mise run build-tinygo`; hardware pin configuration is implemented in `driver/zerokb02`. USB HID and host serial are outside this example.
