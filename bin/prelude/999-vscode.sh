function code {
  flatpak --socket=wayland run com.visualstudio.code --enable-features=UseOzonePlatform --ozone-platform-hint=auto
 "$@"
}
