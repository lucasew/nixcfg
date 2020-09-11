let
  cfg = rec {
      machine_name = "acer";
      username = "lucasew";
      email = "lucas59356@gmail.com";
      # selectedDesktopEnvironment = "xfce_i3";
      selectedDesktopEnvironment = "gnome";
      hostname = "acer-nix";
      wallpaper = builtins.fetchurl {
        url = "http://wallpaperswide.com/download/armenia_syunik_khustup_hayk_k13-wallpaper-1366x768.jpg";
        sha256 = "1z2439f0d8hpqwjp07xhwkcp7svzvbhljayhdfssmvi619chlc0p";
      };
      # wallpaper = builtins.fetchurl {
      #   url = "http://wallpaperswide.com/download/aurora_sky-wallpaper-1366x768.jpg";
      #   sha256 = "1gk4bw5mj6qgk054w4g0g1zjcnss843afq5h5k0qpsq9sh28g41a";
      # };


      dotfileRootPath = 
      let
        env = builtins.getEnv "DOTFILES";
        envNotNull = assert (env != ""); env;
        envExists = assert (builtins.pathExists envNotNull); envNotNull;
      in envExists;

      nixpkgs = "${builtins.getEnv "HOME"}/.nix-defexpr/channels/nixos-unstable";
      home-manager = "${builtins.getEnv "HOME"}/.nix-defexpr/channels/home-manager";
      nur = "${builtins.getEnv "HOME"}/.nix-defexpr/channels/nur";
      pkgs = "${dotfileRootPath}/pkgs.nix";
      setupScript = import ./home/components/dotfiles/gen.nix cfg;
  };
in cfg
