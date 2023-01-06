{global, pkgs, lib, self, config, ...}:
let
  inherit (pkgs) vim gitMinimal tmux xclip writeShellScriptBin;
  inherit (global) username;
in {
  imports = [
    /home/lucasew/WORKSPACE/nixpkgs/nixos/modules/services/monitoring/cockpit.nix
    ./flake-etc.nix
    ./nix.nix
    ./zerotier.nix
    ./user.nix
    ./ssh.nix
    ./colors.nix
    ./motd.nix
    ./screenkey.nix
    (builtins.fetchurl {
      url = "https://raw.githubusercontent.com/NixOS/nixpkgs/60c696e31b14797a346241e4f553399d92ba2b69/nixos/modules/config/dotd.nix";
      sha256 = "0n66xqb2vlv97fcfd3s74qv3dh9yslnvhxhzx3p3rq0vmsq4i2ml";
    })
  ];
  services.nginx.virtualHosts."cockpit.${config.networking.hostName}.${config.networking.domain}" = lib.mkIf config.services.cockpit.enable {
    locations."/" = {
      proxyPass = "http://127.0.0.1:${toString config.services.cockpit.port}";
    };
  };

  boot.cleanTmpDir = true;
  i18n.defaultLocale = "pt_BR.UTF-8";
  time.timeZone = "America/Sao_Paulo";
  environment.systemPackages = [
    vim
    gitMinimal
    tmux
    xclip
  ];
  environment.variables = {
    EDITOR = "nvim";
    PATH = "$PATH";
  };
  programs.bash = {
    promptInit = builtins.readFile ./bash_init.sh;
  };
  networking.domain = lib.mkDefault "lucao.net";
}
