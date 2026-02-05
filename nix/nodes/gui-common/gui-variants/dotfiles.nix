{pkgs, ...}: {
  systemd.user.services.dotfile-waybar = {
    path = with pkgs; [
      script-directory-wrapper
      custom.colorpipe
      bash
    ];
    script = ''
      mkdir ~/.config/waybar -p
      cat $(sdw d root)/nix/nodes/gui-common/gui-variants/hyprland/waybar/style.css | colorpipe > ~/.config/waybar/style.css
    '';
    restartTriggers = [
      ./hyprland/waybar/style.css
    ];
    wantedBy = ["default.target"];
  };

  systemd.user.services.dotfile-hyprland = {
    path = with pkgs; [
      script-directory-wrapper
      custom.colorpipe
      bash
    ];
    script = ''
      mkdir ~/.config/hypr -p
      cat $(sdw d root)/nix/nodes/gui-common/gui-variants/hyprland/hypr/hyprland.conf | colorpipe > ~/.config/hypr/hyprland.conf
    '';
    restartTriggers = [./hyprland/hypr/hyprland.conf];
    wantedBy = ["default.target"];
  };
}
