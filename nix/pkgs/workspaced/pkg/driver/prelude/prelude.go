package prelude

import (
	_ "workspaced/pkg/driver/audio/pulse"
	_ "workspaced/pkg/driver/battery/linux"
	_ "workspaced/pkg/driver/brightness/brightnessctl"
	_ "workspaced/pkg/driver/clipboard/termux"
	_ "workspaced/pkg/driver/clipboard/wlcopy"
	_ "workspaced/pkg/driver/clipboard/xclip"
	_ "workspaced/pkg/driver/media/dbus"
	_ "workspaced/pkg/driver/notification/notify_send"
	_ "workspaced/pkg/driver/notification/termux"
	_ "workspaced/pkg/driver/screen/sway"
	_ "workspaced/pkg/driver/screen/x11"
	_ "workspaced/pkg/driver/screenshot/grim"
	_ "workspaced/pkg/driver/screenshot/maim"
	_ "workspaced/pkg/driver/tray/dbus"
	_ "workspaced/pkg/driver/wallpaper/feh"
	_ "workspaced/pkg/driver/wallpaper/swaybg"
	_ "workspaced/pkg/driver/wm/hyprland"
	_ "workspaced/pkg/driver/wm/i3ipc"
)
