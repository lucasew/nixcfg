{...}:
let
    pkgs = import <nixpkgs> {};
    paint = pkgs.fetchzip {
        url = "https://archive.org/download/MSPaintWinXP/mspaint%20WinXP%20English.zip";
        sha256 = "119c7304szbky9n0d7761qvl09fmg9wh4ilna7fzcj691igly562";
    };
    dll = builtins.fetchurl {
        url = "https://www.dlldump.com/dllfiles/M/mfc42u.dll";
        sha256 = "12mi28j78p8350pn38iqkmcxxz69xmbz7k9ws76i7xv825siv8gi";
    };
    theDerivation =  pkgs.stdenv.mkDerivation {
        name = "mspaint-xp-base";
        version = "1.0";    

        src = paint;

        installPhase = ''
        mkdir -p "$out"
        cp -r "$src/mspaint.exe" "$out/mspaint.exe"
        cp "${dll}" "$out/mfc42u.dll"
        '';
        enablePatchelf = false;
    };
    bin = pkgs.writeShellScriptBin "mspaint" ''
        ${pkgs.wineStable}/bin/wine ${theDerivation}/mspaint.exe
    '';
in {
    home.packages = [bin];
}