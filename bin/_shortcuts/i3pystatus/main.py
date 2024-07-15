#!/usr/bin/env -S sd nix shell
#!nix-shell -p i3pystatus -i python3
from i3pystatus import Status
status = Status()

status.register("clock",
    format="%a %-d %b %X KW%V",)

status.run()
