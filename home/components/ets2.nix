{pkgs, ...}:
let
    version = "1.28.0.10";
    storePath = "/nix/store/mgr9dj4mmwwdzn8hvhgrvkqsb01f8kki-Euro.Truck.Simulator.2.v${version}.Inclu.ALL.DLC";
    bin = pkgs.writeShellScriptBin "ets2" ''
        ${pkgs.wineFull}/bin/wine ${storePath}/bin/win_x86/eurotrucks2.exe 
    '';
in {
    home.packages = [bin];
}