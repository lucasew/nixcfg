flake: final: prev:
let
  inherit (final) callPackage;
in
let
  cp = f: (callPackage f) { };
in
{
  unstable = flake.lib.mkPkgs {
    nixpkgs = flake.inputs.nixpkgs-unstable;
    inherit (prev) system;
  };
  inherit flake;
  bumpkin = rec {
    bumpkin = cp flake.inputs.bumpkin;
    inputs = bumpkin.loadBumpkin {
      inputFile = ../bumpkin.json;
      outputFile = ../bumpkin.json.lock;
    };
    unpacked = (cp ./lib/unpackRecursive.nix) inputs;
  };

  wallabag = prev.wallabag.overrideAttrs (old: {
    postFixup = (old.postFixup or "") + ''
      # exit 1
      echo $out/**/*.yml
      substituteInPlace $out/app/config/services{,_test}.yml \
        --replace-fail '../../src/Wallabag' "$out/src/Wallabag"
    '';
  });

  nbr = import "${flake.inputs.nbr}" { pkgs = final; };

  blender-bin = flake.inputs.blender-bin.packages.${prev.system}; 

  inherit (flake.inputs.nix-alien.packages.${prev.system}) nix-alien;

  pythonPackagesExtensions = [
    (final: prev: {
      pyctcdecode = final.callPackage ./pkgs/python/pyctcdecode { };
      kenlm = final.callPackage ./pkgs/python/kenlm { };
      pyvirtualcam = final.callPackage ./pkgs/python/pyvirtualcam { };
    })
  ];

  python3PackagesBin = prev.python3Packages.overrideScope (
    self: super: {
      torch = super.pytorch-bin;
      # torch = super.torch-bin // {
      #   inherit (super.torch) cudaCapabilities cxxdev;
      #   cudaSupport = true;
      # };
      torchaudio = super.torchaudio-bin;
      torchvision = super.torchvision-bin;
    }
  );

  python3PackagesCuda = prev.python3Packages.overrideScope (_: _: { cudaSupport = true; });

  lib = prev.lib.extend (
    final: prev: {
      jpg2png = cp ./lib/jpg2png.nix;
      buildDockerEnv = cp ./lib/buildDockerEnv.nix;
      climod = cp flake.inputs.climod;
    }
  );

  devenv = final.writeShellScriptBin "devenv" ''
    nix run ${flake.inputs.devenv} -- "$@"
  '';

  ctl = cp ./pkgs/ctl;
  cmdlinegl = cp ./pkgs/cmdlinegl;

  personal-utils = cp ./pkgs/personal-utils.nix;
  fhsctl = cp ./pkgs/fhsctl.nix;
  pkg = cp ./pkgs/pkg.nix;
  text2image = cp ./pkgs/text2image.nix;
  wrapWine = cp ./pkgs/wrapWine.nix;
  ollama-webui = cp ./pkgs/ollama-webui;
  home-manager = cp "${flake.inputs.home-manager}/home-manager";

  prev = prev;
  requireFileSources = [ flake.inputs.nix-requirefile-data ];

  appimage-wrap = final.nbr.appimage-wrap;

  dotenv = cp flake.inputs.dotenv;
  p2k = cp flake.inputs.pocket2kindle;
  pytorrentsearch = cp flake.inputs.pytorrentsearch;
  redial_proxy = cp flake.inputs.redial_proxy;
  send2kindle = cp flake.inputs.send2kindle;
  nixgram = cp flake.inputs.nixgram;
  go-annotation = cp flake.inputs.go-annotation;
  ts-proxy = flake.inputs.ts-proxy.packages.${prev.system}.default;
  wrapVSCode = args: import flake.inputs.nix-vscode (args // { pkgs = prev; });
  wrapEmacs = args: import flake.inputs.nix-emacs (args // { pkgs = prev; });

  instantngp = cp ./pkgs/instantngp.nix;

  # nix-option = callPackage "${flake.inputs.nix-option}" {
  #   nixos-option = (callPackage "${flake.inputs.nixpkgs}/nixos/modules/installer/tools/nixos-option" { }).overrideAttrs (attrs: attrs // {
  #     meta = attrs.meta // {
  #       platforms = lib.platforms.all;
  #     };
  #   });
  # };

  nur = import flake.inputs.nur { pkgs = prev; };

  wineApps = {
    cs_extreme = cp ./pkgs/wineApps/cs_extreme.nix;
    dead_space = cp ./pkgs/wineApps/dead_space.nix;
    gta_sa = cp ./pkgs/wineApps/gta_sa.nix;
    among_us = cp ./pkgs/wineApps/among_us.nix;
    ets2 = cp ./pkgs/wineApps/ets2.nix;
    mspaint = cp ./pkgs/wineApps/mspaint.nix;
    pinball = cp ./pkgs/wineApps/pinball.nix;
    sosim = cp ./pkgs/wineApps/sosim.nix;
    tora = cp ./pkgs/wineApps/tora.nix;
    nfsu2 = cp ./pkgs/wineApps/nfsu2.nix;
    flatout2 = cp ./pkgs/wineApps/flatout2.nix;
    watchdogs2 = cp ./pkgs/wineApps/watchdogs2.nix;
    rimworld = cp ./pkgs/wineApps/rimworld.nix;
    skyrim = cp ./pkgs/wineApps/skyrim.nix;
  };
  custom = rec {
    kodi = final.kodi.withPackages (kpkgs: with kpkgs; [
      vfs-sftp
      sponsorblock
      joystick
      sendtokodi
    ]);
    i3pystatus = cp ./pkgs/i3pystatus.nix;
    colorpipe = cp ./pkgs/colorpipe;
    chromium = cp ./pkgs/custom/chromium;
    ncdu = cp ./pkgs/custom/ncdu.nix;
    neovim = cp ./pkgs/custom/neovim;
    emacs = cp ./pkgs/custom/emacs;
    firefox = cp ./pkgs/custom/firefox;
    tixati = cp ./pkgs/custom/tixati.nix;
    vscode = cp ./pkgs/custom/vscode;
    rofi_xorg = cp ./pkgs/custom/rofi.nix;
    rofi = final.custom.rofi_xorg;
    rofi_wayland = prev.callPackage ./pkgs/custom/rofi.nix { rofi = final.rofi-wayland; };
    pidgin = cp ./pkgs/custom/pidgin.nix;
    send2kindle = cp ./pkgs/custom/send2kindle.nix;
    retroarch = cp ./pkgs/custom/retroarch.nix;
    loader = cp ./pkgs/custom/loader/default.nix;
    polybar = cp ./pkgs/custom/polybar.nix;
    colors-lib-contrib = import "${flake.inputs.nix-colors}/lib/contrib" { pkgs = prev; };
    # wallpaper = ./wall.jpg;
    wallpaper = colors-lib-contrib.nixWallpaperFromScheme {
      scheme = final.custom.colors;
      width = 1366;
      height = 768;
      logoScale = 2;
    };
    inherit (flake) colors;
  };

  script-directory-wrapper = final.writeShellScriptBin "sdw" ''
    set -eu
    export SD_CMD=
    export SD_ROOT="$(${flake}/bin/source_me sd d root)"
    [ -d "$SD_ROOT" ]
    "$SD_ROOT/bin/source_me" sd "$@"
  '';

  opencv4Full = prev.python3Packages.opencv4.override {
    pythonPackages = prev.python3Packages;
    enablePython = true;
    enableContrib = true;
    enableTesseract = true;
    enableOvis = true;
    enableUnfree = true;
    enableGtk3 = true;
    enableGPhoto2 = true;
    enableFfmpeg = true;
    enableGStreamer = false;
    enableIpp = true;
    enableTbb = true;
    enableDC1394 = true;
  };

  opencv4FullCuda = final.opencv4Full.override {
    enableCuda = true;
    enableCudnn = true;
  };

  intel-ocl = prev.intel-ocl.overrideAttrs (old: {
    src = prev.fetchzip {
      url = "https://github.com/lucasew/nixcfg/releases/download/debureaucracyzzz/SRB5.0_linux64.zip";
      sha256 = "sha256-4qaX7wTqxKSrRWeQv1Zrs6eTT0fKJ6g9QBFocugwd2E=";
      stripRoot = false;
    };
  });

  cached-nix-shell = callPackage flake.inputs.src-cached-nix-shell { pkgs = prev; };

  # rio = flake.inputs.rio.packages.${prev.system}.default;

  arcan = prev.arcan.override { useTracy = false; };

  movfuscator = prev.callPackage ./pkgs/movfuscator { };

  regex101 = prev.callPackage flake.inputs.regex101 { };

  elixir_edge = final.unstable.elixir.override {
    erlang = final.unstable.erlang_27;
  };

  conda = prev.conda.override {
    extraPkgs = with final; [
      libGL
      libGLU
      xorg.libX11
      xorg.libXi
      qt5.qtbase
    ];
  };
}
