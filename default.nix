{
    cfg = import ./config.nix;
    home = import ./home;
    homeConfig = import ./home/config.nix;
    apps = import ./apps;
    machine = import ./machine/${cfg.machine_name};
}