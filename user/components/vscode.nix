{config, pkgs, ...}:
{
    programs.vscode = {
        enable = true;
        package = pkgs.vscode;
        userSettings = {
            "workbench.iconTheme" = "file-icons";
            "workbench.colorTheme" = "One Dark Pro";
        };
        extensions = [
            pkgs.vscode-extensions.bbenoist.Nix
            pkgs.vscode-extensions.vscodevim.vim
        ] ++ pkgs.vscode-utils.extensionsFromVscodeMarketplace [
            {
                name = "file-icons";
                publisher = "file-icons";
                version = "1.0.25";
                sha256 = "0s6lr7s1a0alkknazmch5k2m0r16p5gnlzn3yyan9wl8k3579c25";
            }
            {
                name = "material-theme";
                publisher = "zhuangtongfa";
                version = "3.8.5";
                sha256 = "1fdhykyddzghzs8j701q04lb2rhfrr0sbz0ib0js0shj8v31n8aa";
            }
            {
                name = "remote-containers";
                publisher = "ms-vscode-remote";
                version = "0.117.1";
                sha256 = "0kq3wfwxjnbhbq1ssj7h704gvv1rr0vkv7aj8gimnkj50jw87ryd";
            }
        ];
    };

}