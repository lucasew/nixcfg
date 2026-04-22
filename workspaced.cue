package workspaced

workspaced: {

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

	drivers: {
		"workspaced/pkg/driver/terminal.Driver": {
			"terminal_kitty": 70
		}
	}


	lazy_tools: {
		zed: {
			ref: "github:zed-industries/zed"
			bins: ["zed"]
		}
		codex: {
			ref: "mise:codex"
			bins: ["codex"]
			global: true
		}
		bun: {
			ref: "mise:bun"
			bins: ["bun"]
		}
		docker_compose: {
			ref: "github:docker/compose"
			global: true
			bins: ["docker-compose"]
		}
		fd: {
			ref: "github:sharkdp/fd"
			global: true
			bins: ["fd"]
		}
		opencode: {
			ref: "github:anomalyco/opencode"
			global: true
			bins: ["opencode"]
		}
		rclone: {
			ref: "github:rclone/rclone"
			global: true
			bins: ["rclone"]
		}
		rtk: {
			ref: "github:rtk-ai/rtk"
			global: true
			bins: ["rtk"]
		}
		ripgrep: {
			ref: "github:burntsushi/ripgrep"
			global: true
			bins: ["rg"]
		}
		fzf: {
			ref: "github:junegunn/fzf"
			global: true
			bins: ["fzf", "fzf-tmux"]
		}
		ruff: {
			ref: "github:astral-sh/ruff"
			bins: ["ruff"]
		}
		shfmt: {
			ref: "github:patrickvane/shfmt"
			bins: ["shfmt"]
		}
		shellcheck: {
			ref: "github:koalaman/shellcheck"
			bins: ["shellcheck"]
		}
		sops: {
			ref: "github:getsops/sops"
			bins: ["sops"]
		}
		uv: {
			ref: "github:astral-sh/uv"
			global: true
			bins: ["uv", "uvx"]
		}
		node: {
			ref: "mise:node"
			global: true
			bins: ["node", "npm", "npx"]
		}
		yt_dlp: {
			ref: "mise:yt-dlp"
			global: true
			bins: ["yt-dlp"]
		}
		terraform: {
			ref: "github:hashicorp/terraform"
			bins: ["terraform"]
		}
		tflint: {
			ref: "github:terraform-linters/tflint"
			bins: ["tflint"]
		}
		go: {
			ref: "mise:go"
			bins: ["go", "gofmt"]
		}
		golangci_lint: {
			ref: "github:golangci/golangci-lint"
			bins: ["golangci-lint"]
		}
		docker_language_server: {
			ref: "github:docker/docker-language-server"
			global: true
			bins: ["docker-langserver"]
		}
		ffmpeg: {
			ref: "github:ffbinaries/ffbinaries-prebuilt"
			bins: ["ffmpeg", "ffprobe"]
		}
		helix: {
			ref: "github:helix-editor/helix"
			global: true
			bins: ["hx"]
		}
		ltex_ls: {
			ref: "github:valentjn/ltex-ls"
			global: true
			bins: ["ltex-ls"]
		}
		bash_language_server: {
			ref: "mise:npm:bash-language-server"
			global: true
			bins: ["bash-language-server"]
		}
		vscode_langservers: {
			ref: "mise:npm:vscode-langservers-extracted"
			bins: [
				"vscode-html-language-server",
				"vscode-css-language-server",
				"vscode-json-language-server",
				"vscode-eslint-language-server",
			]
		}
		clang: {
			ref: "mise:vfox:clang"
			bins: ["clang", "clang-cpp"]
		}
		jless: {
			ref: "jless"
			global: true
			bins: ["jless"]
		}
		gh: {
			ref: "github:cli/cli"
			global: true
			bins: ["gh"]
		}
		ast_grep: {
			ref: "github:ast-grep/ast-grep"
			global: true
			bins: ["ast-grep", "sg"]
		}
		scc: {
			ref: "github:boyter/scc"
			global: true
			bins: ["scc"]
		}
		tirith: {
			ref: "github:sheeki03/tirith"
			global: true
			bins: ["tirith"]
		}
		jujutsu: {
			ref: "github:jj-vcs/jj"
			global: true
			bins: ["jj"]
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
		hermes: {input: "self:modules/hermes", enable: true}
	}
}
