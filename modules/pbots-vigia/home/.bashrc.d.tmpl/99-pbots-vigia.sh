# shellcheck shell=bash
# Start personal-bots vigia once per shell init on grokbot.
# vigia.py's own lock prevents duplicate loops.
if [[ -f "${HOME}/.personal/personal-bots/bin/vigia.py" ]]; then
	mkdir -p "${HOME}/.cache/pbots"
	(cd "${HOME}/.personal/personal-bots" && setsid -f mise x -- bin/vigia.py --loop 1800 >>"${HOME}/.cache/pbots/vigia.log" 2>&1)
fi
