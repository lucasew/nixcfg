{pkgs, ...}:
let
  pluginNocapsquit = pkgs.vimUtils.buildVimPlugin {
    name = "nocapsquit";
    src = pkgs.fetchFromGitHub {
        owner = "lucasew";
        repo = "nocapsquit.vim";
        rev = "4418b78b635e797eab915bc54380a2a7e66d2e84";
        sha256 = "1jwwiq321b86bh1z3shcprgh2xs5n1xjy9s364zxlxy8qhwfsryq";
    };
  };
  customNeovim = pkgs.neovim.override {
    viAlias = true;
    vimAlias = true;
    configure = {
      plug.plugins = with pkgs.vimPlugins; [
        LanguageClient-neovim
        auto-pairs
        echodoc
        indentLine
        lightline-vim
        onedark-vim
        pluginNocapsquit
        vim-commentary
        vim-nix
        vim-startify
      ];
      customRC = ''
      let g:LanguageClient_serverCommands = ${builtins.toJSON (import ./langservers.nix {inherit pkgs;})}
      set completefunc=LanguageClient#compltete
      ${builtins.readFile ./rc.vim}
      '';
    };
  };
in
{
  environment.systemPackages = [customNeovim];
  environment.variables.EDITOR = "nvim";
}
