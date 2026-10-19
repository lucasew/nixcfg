package module

#hex: {
	for k, v in modules.base16.config if k =~ "^base" {
		(k): "#\(v)"
	}
}

modules: "base16-herdr": config: {}

// Herdr 0.8.0 Palette fields. Roles match modules/README.md.
// There is no built-in base16 theme name; these overrides sit on theme.name.
file: home: {
	".config/herdr/config.toml": {
		type: "toml"
		values: {
			theme: custom: {
				accent:      #hex.base0D
				panel_bg:    #hex.base00
				surface0:    #hex.base02
				surface1:    #hex.base03
				surface_dim: #hex.base01
				overlay0:    #hex.base04
				overlay1:    #hex.base06
				text:        #hex.base05
				subtext0:    #hex.base04
				mauve:       #hex.base0E
				green:       #hex.base0B
				yellow:      #hex.base0A
				red:         #hex.base08
				blue:        #hex.base0D
				teal:        #hex.base0C
				peach:       #hex.base09
			}
		}
	}
}
