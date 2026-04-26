# Junim (Android Tweaks)

Notes on custom tweaks applied via `adb` or manually to the device.

## Increase Bluetooth Audio Volume Limit
To bypass the default volume cap for A2DP Bluetooth devices, run this via ADB:

```sh
adb shell content insert --uri content://settings/system --bind value:s:20 --bind name:s:volume_music_bt_a2dp
```

## Other Modifications
- Lowered display DPI to `100` for more screen real estate.
- Installed `AnySoftKeyboard` as a lightweight keyboard, stripped down to bare essentials.
