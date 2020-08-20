{pkgs, ...}:
let
    version = "1.28.0.10";
    game = "/run/media/lucasew/Dados/DADOS/Jogos/Euro.Truck.Simulator.2.v1.28.0.10.Inclu.ALL.DLC";
    bin = pkgs.writeShellScriptBin "ets2" ''
        ${pkgs.wineFull}/bin/wine ${game}/bin/win_x86/eurotrucks2.exe 
    '';
in {
    home.packages = [bin];
}