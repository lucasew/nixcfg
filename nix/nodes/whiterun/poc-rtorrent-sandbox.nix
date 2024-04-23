{ config, lib, ...}:

{
  systemd.services.rtorrent-sandbox-poc.serviceConfig = lib.mkMerge [
    config.systemd.services.rtorrent.serviceConfig
    {
      ExecStart = lib.mkForce "/run/rtorrent-poc-payload";  
      RuntimeDirectory = lib.mkForce "rtorrent-sandbox-poc";
      BindPaths = [ "/run" ];
    }
  ];
}