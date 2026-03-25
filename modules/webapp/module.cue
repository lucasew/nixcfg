package module

module: {
	meta: {
		requires: []
		recommends: []
	}

	config: {
		enable: bool | *true
		apps?: [string]: {
			url: string
			profile?: string
			desktop_name?: string
			icon?: string
		}
	}
}
