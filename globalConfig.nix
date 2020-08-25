rec {
    machine_name = "acer";
    username = "lucasew";
    email = "lucas59356@gmail.com";
    selectedDesktopEnvironment = "xfce_i3";
    hostname = "acer-nix";
    wallpaper = builtins.fetchurl {url = "http://wallpaperswide.com/download/aurora_sky-wallpaper-1366x768.jpg";};

    overlaysPath = ./overlays;
    dotfileRootPath = builtins.toPath (overlaysPath + "/.."); 

    nixpkgs = builtins.fetchTarball {
        url = "https://github.com/NixOS/nixpkgs/archive/20.03.tar.gz";
        sha256 = "0182ys095dfx02vl2a20j1hz92dx3mfgz2a6fhn31bqlp1wa8hlq";
    };
    latestNixpkgs = builtins.fetchTarball {
      url = "https://github.com/NixOS/nixpkgs/archive/master.tar.gz";
    };
    defaultPkgs = import nixpkgs {};
    setupScript = ''
    export NIXPKGS_ALLOW_UNFREE=1
    export NIXOS_CONFIG=${dotfileRootPath}/machine/${machine_name}
    NIX_PATH=nixpkgs-overlays=${overlaysPath}:nixpkgs=${nixpkgs}:nixos-config=$NIXOS_CONFIG:dotfiles=${dotfileRootPath}

    alias nixos-rebuild="sudo -E nixos-rebuild"
    alias nixos-install="sudo -E nixos-install --system $NIXOS_CONFIG"
    '';
}
