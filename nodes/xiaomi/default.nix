{pkgs, config, options, ...}:
{
  system.stateVersion = "21.05";
  time.timeZone = "America/Sao_Paulo";
  home-manager = {
    useGlobalPkgs = true;
    config = {...}: {
      home.stateVersion = "21.05";
    };
  };
  nix = {
    package = pkgs.nix;
    extraConfig = ''
      experimental-features = nix-command flakes
    '';
  };
  environment.packages = with pkgs; [
    nix-option
    custom.neovim
    custom.emacs
    git
    htop
    pkg
    # some defaults on the default dotfile
    hostname
    gnugrep
    gnused
    gnutar
    gzip
    xz
    zip
    unzip
  ];
}
