package module

module: {
	meta: {
		requires: []
		recommends: []
	}

	config: {
		apps?: [string]: {
			url: string
			profile?: string
			desktop_name?: string
			icon?: string
		}
	}
}
