# Consistently Ignored Changes

This file lists patterns of changes that have been consistently rejected by human reviewers. All agents MUST consult this file before proposing a new change. If a planned change matches any pattern described below, it MUST be abandoned.

---

## IGNORE: Command Injection in Calculator

**- Pattern:** Do not "fix" the command injection vulnerability in the `repl` function of the calculator script.
**- Justification:** The script is a developer tool where the user intentionally passes Python code to be executed. The `exec(arg)` call is the intended functionality, not a flaw. Attempts to remove it are consistently rejected as they break the tool's core purpose.
**- Files Affected:** `bin/.calc/commoncalc/__common.py`

---

## IGNORE: Command Injection from arp-scan

**- Pattern:** Do not report command injection vulnerabilities related to parsing the output of `arp-scan`.
**- Justification:** This script is intended to be run in a trusted local network environment. The risk of a malicious actor manipulating ARP responses to inject commands is considered an acceptable, low-probability risk. Hardening this script is not a priority and pull requests for it will be closed.
**- Files Affected:** `bin/misc/dns-cgi`

---

## IGNORE: Purely Stylistic Refactoring of Scripts

**- Pattern:** Do not refactor shell scripts for purely stylistic reasons or to change code structure without altering functionality.
**- Justification:** Changes that only rearrange code, such as extracting logic into functions without fixing a bug or adding a feature, are often rejected. The owner may have a specific reason for the existing structure, and such changes create noise. The rejected refactoring of the phone backup logic in `bin/backup` is a key example.
**- Files Affected:** `bin/backup`

---
