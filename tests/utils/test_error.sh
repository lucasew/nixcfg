#!/usr/bin/env bash

# Source the error utility
source "$(dirname "$0")/../../src/utils/error.sh"

test_reportError_output() {
    local output
    output=$(reportError "Test error message" 2>&1 || true)

    if [[ ! "$output" == *"ERROR: Test error message"* ]]; then
        echo "FAIL: Expected 'ERROR: Test error message' in output, got: $output"
        return 1
    fi

    if [[ ! "$output" == *"Context:"* ]]; then
        echo "FAIL: Expected 'Context:' in output, got: $output"
        return 1
    fi

    echo "PASS: test_reportError_output"
}

test_reportError_exit_code() {
    reportError "Test error message" 42 2>/dev/null
    local exit_code=$?

    if [[ $exit_code -ne 42 ]]; then
        echo "FAIL: Expected exit code 42, got: $exit_code"
        return 1
    fi

    echo "PASS: test_reportError_exit_code"
}

run_tests() {
    echo "Running error reporting tests..."
    test_reportError_output
    test_reportError_exit_code
    echo "All tests passed."
}

run_tests
