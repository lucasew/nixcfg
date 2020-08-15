{pkgs, ...}: 
    # FIXME: Can't hear that lovely music and the sound effects
let 
    pkgs = import <nixpkgs> {};
    pinballSource = pkgs.stdenv.mkDerivation rec {
        name = "mspinball";
        version = "1.0";
        src = pkgs.fetchurl {
            url = "https://archive.org/download/SpaceCadet_Plus95/Space_Cadet.rar";
            sha256 = "3cc5dfd914c2ac41b03f006c7ccbb59d6f9e4c32ecfd1906e718c8e47f130f4a";
        };
        unpackPhase = ''
        mkdir -p $out
        cd $out && ${pkgs.unrar}/bin/unrar x ${src}
        '';
        installPhase = ''
        cp -r $src $out
        '';
        enablePatchElf = false;
    };
    bin = pkgs.writeShellScriptBin "pinball" ''
    cd "${pinballSource}"; ${pkgs.wine}/bin/wine "${pinballSource}/PINBALL.exe"
    '';
in
{
    home.packages = [bin];
}