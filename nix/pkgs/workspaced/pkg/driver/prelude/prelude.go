package prelude

import (
	_ "workspaced/pkg/driver/clipboard/termux"
	_ "workspaced/pkg/driver/clipboard/wlcopy"
	_ "workspaced/pkg/driver/clipboard/xclip"
	_ "workspaced/pkg/driver/screen/sway"
	_ "workspaced/pkg/driver/screen/x11"
	_ "workspaced/pkg/driver/tray/dbus"
	_ "workspaced/pkg/driver/wm/hyprland"
	_ "workspaced/pkg/driver/wm/i3ipc"
)
