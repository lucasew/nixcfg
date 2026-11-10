# shellcheck shell=bash
# aichat honors AICHAT_CONFIG_DIR on every OS. Pin it so the config written by
# modules/aichat (~/.config/aichat/config.yaml) is found on macOS too, where
# dirs::config_dir() would point at ~/Library/Application Support/aichat.
export AICHAT_CONFIG_DIR="$HOME/.config/aichat"
