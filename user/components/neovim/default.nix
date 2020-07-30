{config, pkgs, ...}:

{
    programs.neovim = {
            enable = true;
            viAlias = true;
            vimAlias = true;
            vimdiffAlias = true;
            extraConfig =  ''
source ${pkgs.vimPlugins.vim-plug}/share/vim-plugins/vim-plug/plug.vim
${builtins.readFile ./neovim.vim}
            '';
    };


}