{ writeShellScriptBin }:

writeShellScriptBin "workspaced" ''
  dotfilesFolder=
  if [ -d ~/.dotfiles ]; then
    dotfilesFolder=~/.dotfiles
  elif [ -d /home/lucasew/.dotfiles ]; then
    dotfilesFolder=/home/lucasew/.dotfiles
  elif [ -d /etc/.dotfiles ]; then
    dotfilesFolder=/etc/.dotfiles
  fi
  if [ -z "$dotfilesFolder" ]; then
    echo "can't find dotfiles folder" >&2
    exit 1
  fi
  exec "$dotfilesFolder/bin/shim/workspaced" "$@"
''
