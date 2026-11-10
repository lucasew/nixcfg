package module

options: {}

// aichat talking to the aperture OpenAI-compatible gateway.
// aichat reads this directory on every OS because the shell exports
// AICHAT_CONFIG_DIR (config/.bashrc.d.tmpl/20-env-aichat.sh); upstream
// would use ~/Library/Application Support/aichat on macOS.
file: home: {
	".config/aichat/config.yaml": {
		type: "text"
		values: {
			content: """
				clients:
				  - type: openai-compatible
				    name: aperture
				    api_base: https://ai.stargazer-shark.ts.net/v1
				    api_key: unused
				    models:
				      - name: deepseek_coding_flash
				        supports_function_calling: true
				    patch:
				      chat_completions:
				        '.*':
				          headers:
				            x-opencode-session: aichat-aperture
				default_model: aperture:deepseek_coding_flash
				""" + "\n"
		}
	}
}
