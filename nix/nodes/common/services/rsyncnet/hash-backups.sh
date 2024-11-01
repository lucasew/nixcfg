#!/usr/bin/env bash

# TODO: comment out for prod
# function wssh {
#   ssh de3163@de3163.rsync.net -- "$@"
# }

function temp_file {
  echo "$(basename "$(mktemp -u)").$1"
}

OUT_DIR="$(pwd)"
# OUT_DIR="/tmp/rsyncnet"
mkdir -p "$OUT_DIR"

cd "$OUT_DIR" || exit

last_snapshots=($(ls | grep -e ^custom_daily | sort | tail -n 2))

echo "last_snapshots: ${last_snapshots[@]}"

if ! ((${#last_snapshots[@]} < 2)); then
  {
    echo Subject: Diferença entre os snapshots: ${last_snapshots[@]}
    echo
    git diff "${last_snapshots[0]}" "${last_snapshots[1]}" | awk '/^[+-][a-z0-9]/ { print $0 }' | grep -v git/zz | sort -k 2
  } | sendmail
fi

for snapshot in $(wssh ls .zfs/snapshot); do
  snapshot_file="${snapshot}.txt"
  if [[ -f "${snapshot}.txt"  ]]; then
    echo "$snapshot já dumpado!" >&2
    continue
  fi
  tmpfile="$(temp_file "$snapshot")"
  time wssh rclone sha1sum ".zfs/snapshot/${snapshot}" | pv | sort > "$tmpfile" && mv "$tmpfile" "$snapshot_file" && echo "$snapshot dumpado!" >&2
done
