#!/usr/bin/env bash
# Centralized Error Reporting Utility
# Enforces the project convention that all code paths handling unexpected errors
# MUST funnel through a single, centralized error-reporting function.

reportError() {
    local error_message="$1"
    local exit_code="${2:-1}"
    local timestamp
    timestamp="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"

    local caller_script="${BASH_SOURCE[1]:-${0}}"
    local caller_line="${BASH_LINENO[0]:-unknown}"
    local caller_func="${FUNCNAME[1]:-unknown}"

    # Construct the error payload
    local error_payload="[${timestamp}] ERROR: ${error_message}"
    error_payload+=" | Context: ${caller_script}:${caller_line} (in ${caller_func})"

    # Report to stderr (and potentially Sentry if available in the future)
    echo "${error_payload}" >&2

    # If a specific error reporting backend like Sentry is configured, it would be called here.
    # e.g., if [[ -n "$SENTRY_DSN" ]]; then sentry-cli send-event -m "${error_payload}" ; fi

    return "${exit_code}"
}
