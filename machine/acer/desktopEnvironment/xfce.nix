{config, pkgs, ...}:
{
    services.xserver = {
        desktopManager.xfce.enable = true;
        xautolock = {
            enable = true;
            time = 10;
            killtime = 24*60;
        };
    };
}
