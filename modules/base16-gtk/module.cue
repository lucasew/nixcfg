package module

module: {
	meta: {
		requires: ["base16"]
		recommends: []
	}

	config: {
		enable: bool | *true
		theme_name: string | *"base16"
		dconf: {
			"org/gnome/desktop/interface": {
				"gtk-theme": string | *theme_name
			}
		}
	}

	drivers: {
		"workspaced/pkg/driver/notification.Driver": {
			notification_dbus:        100
			notification_notify_send: 10
		}
	}
}
