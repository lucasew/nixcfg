{pkgs, config, ...}:
{
    home.packages = with pkgs; [
        rclone
        rclone-browser
    ];
}