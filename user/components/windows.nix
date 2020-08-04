{pkgs, config, ...}:
{
    home.packages = with pkgs; [
        wine
    ];
}