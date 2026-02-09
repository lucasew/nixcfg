{...}: {
  # tmux configuration is now managed by workspaced templates
  # See: config/.config/tmux/tmux.conf.tmpl
  programs.tmux = {
    enable = true;
    extraConfig = ''
      source-file ~/.config/tmux/tmux.conf
    '';
  };
}
