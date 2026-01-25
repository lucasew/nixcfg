function notification {
	local notification_id="$RANDOM"
	# echo $notification_id
	local title="Notification"
	local message="Notification message"
	local progress=""

	while [[ $# -gt 0 ]]; do
		case "$1" in
		-h | --help)
			cat <<EOF
notification: cria notificação usando notify-send ou termux

  -t, --title: Título da notificação
  -m, --message: Mensagem da notificação
  -i, --id: ID para atualização de notificação, só reusar para substituir
  -p, --progress: Porcentagem de barra de progresso
EOF
			return 0
			;;
		-t | --title)
			if [[ -n "${2:-}" ]]; then
				title="$2"
				shift
			else
				echo "notification: argument required for $1" >&2
				return 1
			fi
			;;
		-m | --message)
			if [[ -n "${2:-}" ]]; then
				message="$2"
				shift
			else
				echo "notification: argument required for $1" >&2
				return 1
			fi
			;;
		-i | --id)
			if [[ -n "${2:-}" ]]; then
				notification_id="$2"
				shift
			else
				echo "notification: argument required for $1" >&2
				return 1
			fi
			;;
		-p | --progress)
			if [[ -n "${2:-}" ]]; then
				progress="$2"
				shift
			else
				echo "notification: argument required for $1" >&2
				return 1
			fi
			;;
		esac
		shift
	done

	if has_binary termux-notification; then
		if [[ -n "$progress" ]]; then
			title="$title ($progress%)"
		fi
		termux-notification -i "$notification_id" -t "$title" -c "$message"
	else
		local extra_args=()
		if [[ -n "$progress" ]]; then
			extra_args+=(-h "int:value:$progress")
		fi
		notify-send -r "$notification_id" "$title" "$message" "${extra_args[@]}"
	fi
}
