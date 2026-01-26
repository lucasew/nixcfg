#!/usr/bin/env bash

# /**
#  * Shared configuration and utilities for rsync.net operations.
#  *
#  * This script handles:
#  * - Configuring the `rsyncnet_user` variable (with defaults and validation).
#  * - Defining the `update_status` function for notifications.
#  *
#  * Usage:
#  *   NOTIFICATION_TITLE="My Task"
#  *   source bin/lib/rsyncnet.sh
#  */

: "${NOTIFICATION_TITLE:?NOTIFICATION_TITLE must be set before sourcing this script}"

rsyncnet_user="${RSYNCNET_USER:-}"
if [ -z "$rsyncnet_user" ]; then
	echo "WARNING: Using hardcoded rsync.net credentials. Set RSYNCNET_USER to override." >&2
	rsyncnet_user=de3163@de3163.rsync.net
fi

if [[ "$rsyncnet_user" == -* ]]; then
    echo "ERROR: Invalid RSYNCNET_USER: cannot start with a hyphen." >&2
    exit 1
fi

# /**
#  * Sends a notification and prints a status message to stderr.
#  *
#  * Uses `sd source_me notification` to send a desktop notification.
#  * The notification title is determined by the global `NOTIFICATION_TITLE`.
#  *
#  * @param args The message to display/log.
#  */
function update_status {
	sd source_me notification --id $$ --title "$NOTIFICATION_TITLE" --message "$*"
	echo '[*]' "$@" >&2
}
