package module

module: {
	meta: {
		requires: ["base16"]
		recommends: []
	}

	config: {
		theme_name: string | *"base16"
		dconf: {
			"org/gnome/desktop/interface": {
				"gtk-theme": string | *theme_name
				"color-scheme": string
				if workspaced.modules.base16.config.dark_mode {
					"color-scheme": *"prefer-dark" | string
				}
				if workspaced.modules.base16.config.dark_mode == false {
					"color-scheme": *"prefer-light" | string
				}
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
