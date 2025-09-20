function codew {
  flatpak --socket=wayland run com.visualstudio.code --enable-features=UseOzonePlatform --ozone-platform-hint=auto "$@"
}

function code {
  "$(which code)"  --enable-features=UseOzonePlatform --ozone-platform-hint=auto "$@"
}
