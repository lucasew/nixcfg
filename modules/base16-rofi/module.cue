package module

module: {
	meta: {
		requires: ["base16"]
		recommends: []
	}

	config: {
		enable: bool | *true
	}

	drivers: {
		"workspaced/pkg/driver/notification.Driver": {
			notification_dbus:        100
			notification_notify_send: 10
		}
	}
}
