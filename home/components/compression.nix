{config, pkgs, ...}:
{
    home.packages = with pkgs; [
        xarchiver
        unzip
        p7zip
    ];
}
