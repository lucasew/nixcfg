globalConfig: with globalConfig; ''
  export DOTFILES=${dotfileRootPath}
  export NIXPKGS_ALLOW_UNFREE=1
  export NIXOS_CONFIG=${dotfileRootPath}/machine/${machine_name}
  NIX_PATH=nixpkgs=${nixpkgs}:nixpkgs-overlays=${dotfileRootPath}/overlays:nixos-config=$NIXOS_CONFIG:dotfiles=${dotfileRootPath}

  alias nixos-rebuild="sudo -E nixos-rebuild"
  alias nixos-install="sudo -E nixos-install --system $NIXOS_CONFIG"
''
