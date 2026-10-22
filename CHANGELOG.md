# Changelog

## [v0.2.1](https://github.com/rin2yh/rp2040-simulator/compare/v0.2.0...v0.2.1) - 2026-10-22

- feat: zero-kb02 の入力と OLED をデスクトップから利用可能にする by @rin2yh in https://github.com/rin2yh/rp2040-simulator/pull/14
- fix: update Go and TinyGo for security fixes by @rin2yh in https://github.com/rin2yh/rp2040-simulator/pull/20
- Bump Songmu/tagpr from 1.20.3 to 1.21.0 by @dependabot[bot] in https://github.com/rin2yh/rp2040-simulator/pull/16
- Bump acifani/setup-tinygo from 3.0.0 to 3.0.1 by @dependabot[bot] in https://github.com/rin2yh/rp2040-simulator/pull/18
- Bump github.com/hajimehoshi/ebiten/v2 from 2.10.2 to 2.10.4 by @dependabot[bot] in https://github.com/rin2yh/rp2040-simulator/pull/17

## [v0.2.0](https://github.com/rin2yh/rp2040-simulator/compare/v0.1.0...v0.2.0) - 2026-09-24

- Bump jdx/mise-action from 4.2.3 to 4.3.0 by @dependabot[bot] in https://github.com/rin2yh/rp2040-simulator/pull/1
- Bump golang.org/x/image from 0.45.0 to 0.46.0 by @dependabot[bot] in https://github.com/rin2yh/rp2040-simulator/pull/2
- Refine CI responsibilities by @rin2yh in https://github.com/rin2yh/rp2040-simulator/pull/5
- Add conf2025badge board and XIAO RP2040 support by @rin2yh in https://github.com/rin2yh/rp2040-simulator/pull/6
- Remove root emulator entrypoint by @rin2yh in https://github.com/rin2yh/rp2040-simulator/pull/9
- テーブル駆動テストを整理する by @rin2yh in https://github.com/rin2yh/rp2040-simulator/pull/10
- Use setup-go caching in CI by @rin2yh in https://github.com/rin2yh/rp2040-simulator/pull/11
- docs: remove board-specific README screenshot by @rin2yh in https://github.com/rin2yh/rp2040-simulator/pull/12

## [v0.1.0](https://github.com/rin2yh/rp2040-simulator/commits/v0.1.0) - 2026-09-21

- Add the zero-kb02 desktop emulator with OLED, keys, encoder, joystick, LEDs, BOOT, and RESET controls.
- Add desktop and TinyGo-compatible `machine` and `driver/ws2812` packages.
- Add the LED blink example and CI for Go, TinyGo, GUI rendering, and cross-platform builds.
