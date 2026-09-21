# Internal Rules

完成画面は`board/testdata/<board>.png`をgolden imageとして比較する。意図したGUI変更では`UPDATE_BOARD_GOLDEN=1 mise run screenshot`を実行し、`build/emulator.png`を目視確認してからgolden imageを更新する。
