{...}: {
  # The workspaced daemon is now managed automatically via shell prelude
  # using 'workspaced daemon --try'. This allows it to run on Android
  # and other environments without systemd.
  config = {};
}
