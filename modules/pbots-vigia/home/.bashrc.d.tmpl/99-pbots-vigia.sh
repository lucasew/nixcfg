# shellcheck shell=bash
# Start personal-bots vigia on grokbot. vigia-up is idempotent.
if [[ -x "${HOME}/.personal/personal-bots/bin/vigia-up" ]]; then
	"${HOME}/.personal/personal-bots/bin/vigia-up"
fi
