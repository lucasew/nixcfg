package workspaced

workspaced: {
	workspaces: {
		www:                      1
		calendar:                 2
		meta_ai_claude:           1500
		meta_ai_jules:            1501
		meta_dbeaver:             1600
		meta_steam:               1700
		meta_keepass:             1801
		meta_obsidian:            1802
		meta_fava_beancount:      1803
		meta_dotfiles:            1804
		meta_super_productivity:  1805
		comm_telegram:            1900
		comm_element:             1901
		comm_whatsapp:            1902
		comm_meet:                1903
		comm_discord_profissional: 1904
		comm_discord_pessoal:     1905
		comm_discord_video:       1906
		work_lewtec:              2000
		work_lewtec_vps:          2001
		work_lewtec_contapila:    2003
		work_lewtec_romaneiro:    2004
		work_lewtec_launcher:     2005
		work_lewtec_homepage:     2006
		work_lewtec_infra:        2007
		work_mestrado:            2200
		side_blog:                3000
		side_nixpkgs:             3002
		side_contapila:           3004
		side_rinha_backend:       3005
		side_curso_latex:         3006
		side_uilab:               3007
		side_miseci:              3008
		side_fetchurl:            3009
		stat_riverwood:           4000
		stat_whiterun:            4001
		stat_ravenrock:           4002
		stat_dgxa100:             4003
		video_obs_live:           8000
		waypipe:                  9999
	}

	desktop: wallpaper: {
		dir:     "~/.dotfiles/assets"
		default: "wall.jpg"
	}

	hosts: {
		riverwood: {
			ips:  ["100.82.35.120", "192.168.69.2"]
			port: 22
		}
		whiterun: {
			ips:  ["100.85.38.19", "192.168.69.1"]
			mac:  "a8:a1:59:9c:ab:32"
			port: 22
		}
		ravenrock: {
			ips:  ["100.122.87.59"]
			port: 22
		}
		phone: {
			ips:  ["100.76.88.29", "192.168.69.4"]
			port: 22
		}
	}

	backup: rsyncnet_user: "de3163@de3163.rsync.net"

	browser: {
		default: "zen"
		webapp:  "brave"
	}


	lazy_tools: {
		"github:zed-industries/zed": {
			version: "0.222.4"
			bins: ["zed"]
		}
		codex: {
			version: "0.104.0"
			bins: ["codex"]
			global: true
		}
		bun: {
			version: "1.3.8"
			bins: ["bun"]
		}
		"docker-compose": {
			version: "5.0.1"
			global: true
			bins: ["docker-compose"]
		}
		fd: {
			version: "10.3.0"
			global: true
			bins: ["fd"]
		}
		opencode: {
			version: "1.1.40"
			global: true
			bins: ["opencode"]
		}
		rclone: {
			version: "1.72.1"
			global: true
			bins: ["rclone"]
		}
		ripgrep: {
			version: "15.1.0"
			global: true
			bins: ["rg"]
		}
		fzf: {
			version: "0.67.0"
			global: true
			bins: ["fzf", "fzf-tmux"]
		}
		ruff: {
			version: "0.14.14"
			bins: ["ruff"]
		}
		shfmt: {
			version: "3.12.0"
			bins: ["shfmt"]
		}
		shellcheck: {
			version: "0.11.0"
			bins: ["shellcheck"]
		}
		sops: {
			version: "3.11.0"
			bins: ["sops"]
		}
		uv: {
			version: "0.9.28"
			global: true
			bins: ["uv", "uvx"]
		}
		"yt-dlp": {
			version: "2025.10.22"
			global: true
			bins: ["yt-dlp"]
		}
		terraform: {
			version: "1.9.0"
			bins: ["terraform"]
		}
		tflint: {
			version: "0.60.0"
			bins: ["tflint"]
		}
		go: {
			version: "1.24.12"
			bins: ["go", "gofmt"]
		}
		"golangci-lint": {
			version: "1.61.0"
			bins: ["golangci-lint"]
		}
		"github:docker/docker-language-server": {
			version: "0.20.1"
			global: true
			bins: ["docker-langserver"]
		}
		"github:ffbinaries/ffbinaries-prebuilt": {
			version: "6.1"
			bins: ["ffmpeg", "ffprobe"]
		}
		"github:helix-editor/helix": {
			version: "25.07.1"
			global: true
			bins: ["hx"]
		}
		"github:valentjn/ltex-ls": {
			version: "16.0.0"
			global: true
			bins: ["ltex-ls"]
		}
		"npm:bash-language-server": {
			version: "5.6.0"
			global: true
			bins: ["bash-language-server"]
		}
		"npm:vscode-langservers-extracted": {
			version: "4.10.0"
			bins: [
				"vscode-html-language-server",
				"vscode-css-language-server",
				"vscode-json-language-server",
				"vscode-eslint-language-server",
			]
		}
		"vfox:clang": {
			version: "18.1.8"
			bins: ["clang", "clang-cpp"]
		}
		jless: {
			version: "0.9.0"
			global: true
			bins: ["jless"]
		}
		gh: {
			version: "2.86.0"
			global: true
			bins: ["gh"]
		}
		"ast-grep": {
			version: "0.40.5"
			global: true
			bins: ["ast-grep", "sg"]
		}
		"github:boyter/scc": {
			version: "3.6.0"
			global: true
			bins: ["scc"]
		}
	}

	inputs: {
		self: {
			from: "self"
		}
		papirus: {
			from: "github:PapirusDevelopmentTeam/papirus-icon-theme"
			version: "HEAD"
		}
	}

	modules: {
		icons: {
			input: "core:base16-icons-linux"
			enable: !workspaced.runtime.is_phone
			config: {
				input_dir: "papirus:Papirus"
			}
		}

		fontconfig: {
			input: "self:modules/fontconfig"
			enable:      true
			config: {
				serif:       "Fira Code"
				sans_serif:  "Fira Code"
				monospace:   "Fira Code"
				emoji:       "Noto Color Emoji"
			}
		}
		webapp: {
			input: "self:modules/webapp"
			enable: true
			config: {
				apps: {
					"castable-iframe": {
						url:     "https://castable-iframe.vercel.app/"
						profile: "castable-iframe"
					}
					gemini: url: "https://gemini.google.com"
					teste:  url: "https://google.com"
					element: {
						url:          "https://app.element.io"
						profile:      "element"
						desktop_name: "Element Matrix"
					}
					reemo: {
						url:          "https://reemo.io"
						desktop_name: "Remote control"
					}
					clickup: {
						url:          "https://clickup.com"
						desktop_name: "ClickUp"
					}
					facebook: {
						url:          "https://facebook.com"
						profile:      "facebook"
						desktop_name: "Facebook"
					}
					pocketcasts: {
						url:          "https://pocketcasts.com"
						desktop_name: "PocketCasts"
						icon:         "sound"
					}
					twitter: {
						url:          "https://twitter.com"
						profile:      "twitter"
						desktop_name: "Twitter"
					}
					"discord-pessoal": {
						url:          "https://discord.com/channels/me"
						profile:      "discord-pessoal"
						desktop_name: "Discord (pessoal)"
					}
					"discord-profissional": {
						url:          "https://discord.com/channels/me"
						profile:      "discord-profissional"
						desktop_name: "Discord (profissional)"
					}
					"youtube-tv": {
						url:          "https://youtube.com/tv"
						profile:      "youtube-tv"
						desktop_name: "YouTube (UI para TV)"
					}
					whatsapp: {
						url:          "web.whatsapp.com"
						profile:      "zap"
						desktop_name: "WhatsApp"
					}
					notion: {
						url:          "notion.so"
						profile:      "notion"
						desktop_name: "Notion"
					}
					duolingo: {
						url:          "duolingo.com"
						desktop_name: "Duolingo"
					}
					"youtube-music": {
						url:          "music.youtube.com"
						profile:      "youtubemusic"
						desktop_name: "Youtube Music"
					}
					planttext: {
						url:          "https://www.planttext.com/"
						desktop_name: "PlantText"
					}
					rainmode: {
						url:          "https://youtu.be/mPZkdNFkNps"
						desktop_name: "Tocar som de chuva"
						icon:         "weather-showers"
					}
					gmail: {
						url:          "gmail.com"
						desktop_name: "GMail"
					}
					keymash: {
						url:          "https://keyma.sh/learn"
						desktop_name: "keyma.sh: Keyboard typing train"
					}
					calendar: {
						url:          "https://calendar.google.com/calendar/u/0/r/customday"
						desktop_name: "Calendário"
						icon:         "x-office-calendar"
					}
					"twitch-dashboard": url: "https://dashboard.twitch.tv/stream-manager"
					"trello-pessoal":   url: "https://trello.com/b/bjoRKSM2/pessoal"
					"trello-side-projects": url: "https://trello.com/b/36ncJYYV/side-projects"
					"trello-dashboard": url: "trello.com"
					"geforce-now":      url: "play.geforcenow.com"
				}
			}
		}
		base16: {
			input: "self:modules/base16"
			enable: true
			config: {
				base00: "282c34"
				base01: "353b45"
				base02: "3e4451"
				base03: "545862"
				base04: "565c64"
				base05: "abb2bf"
				base06: "b6bdca"
				base07: "c8ccd4"
				base08: "e06c75"
				base09: "d19a66"
				base0A: "e5c07b"
				base0B: "98c379"
				base0C: "56b6c2"
				base0D: "61afef"
				base0E: "c678dd"
				base0F: "be5046"
			}
		}
		"base16-shell":   {input: "self:modules/base16-shell", enable: true}
		"base16-helix":   {input: "self:modules/base16-helix", enable: true}
		"base16-vscode":  {input: "self:modules/base16-vscode", enable: true}
		"base16-sway":    {input: "self:modules/base16-sway", enable: true}
		"base16-gtk":     {input: "self:modules/base16-gtk", enable: true}
		"base16-rofi":    {input: "self:modules/base16-rofi", enable: true}
		"base16-dunst":   {input: "self:modules/base16-dunst", enable: true}
		"base16-tmux":    {input: "self:modules/base16-tmux", enable: true}
		"base16-opencode": {input: "self:modules/base16-opencode", enable: true}
		"base16-swaylock": {input: "self:modules/base16-swaylock", enable: true}
		"script-directory": {input: "self:modules/script-directory", enable: true}
		mise: {input: "self:modules/mise", enable: true}
	}
}
