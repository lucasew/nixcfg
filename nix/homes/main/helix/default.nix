{
  lib,
  config,
  pkgs,
  ...
}:
{
  config = lib.mkIf config.programs.helix.enable {
    home.packages = with pkgs.unstable; [
      typos-lsp
      yaml-language-server
      docker-compose-language-service
      nodePackages.bash-language-server
      nodePackages.svelte-language-server
      emmet-language-server
      vscode-langservers-extracted
      marksman
      gopls
      ltex-ls
      jdt-language-server

      pkgs.unstable.helix
    ];
  };
}
