rec {
    cfg = import ./config.nix;
    home = import ./home;
    homeConfig = import ./home/config.nix;
    machine = import "./machine/${cfg.machine_name}";
}