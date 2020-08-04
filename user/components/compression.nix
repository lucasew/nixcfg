{config, pkgs, ...}:
{
    home.packages = with pkgs; [
        xarchiver
        unzip
    ];
}