---
name: ebitengine-screenshot
description: Capture and visually inspect the rendered UI of this project's Ebitengine emulator through its deterministic GUI test. Use when reviewing emulator appearance, validating a visual change, or producing a local emulator screenshot. Do not use for browser screenshots.
---

# Ebitengine screenshot

Run:

```sh
mise run screenshot
```

This renders representative device state, validates the composed frame, and writes `build/emulator.png`. Inspect the resulting file with the local image viewer. Report visible layout or rendering problems, not only whether the test passed.

Keep screenshots as local artifacts unless the user asks to publish or commit them.

When adding or repairing screenshot support:

- Run Ebitengine on the main goroutine; use `TestMain` when a test owns `RunGame`.
- Capture a completed frame after application pixels and overlays have been composed.
- Drive representative state from the test rather than depending on manual input.
- Keep ordinary headless tests usable; gate GUI capture behind a dedicated build tag or task.
- Make the capture terminate by itself and return a test failure when rendering is blank or incomplete.
