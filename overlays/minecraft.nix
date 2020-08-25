self: super:
with super;
let
    launcherZip = pkgs.requireFile {
        name = "ShiginimaSE_v4400.zip";
        sha1 = "61cb768106e6e449158ebb2608ad1327402d9fec";
        url = "https://teamshiginima.com/update/";
    };
    drv = pkgs.stdenv.mkDerivation rec {
        name = "minecraft";
        src = launcherZip;
        dontUnpack = true;
        nativeBuildInputs = with pkgs; [
            makeWrapper
            unzip
        ];
        buildInputs = with pkgs; [
            unzip
        ];
        installPhase = ''
            unzip ${src}
            ls -lha
            mkdir -p $out/share/java $out/bin
            for file in linux_osx/*.jar;
            do
                cat "$file" > $out/share/java/minecraft.jar
            done
            makeWrapper ${jre8}/bin/java $out/bin/minecraft --add-flags "-jar $out/share/java/minecraft.jar"
        '';
        meta = {
            homepage = "https://teamshiginima.com/update/";
            description = "Minecraft";
            # license = stdenv.licences.proprietary;
            platforms = stdenv.lib.platforms.unix;
        };
    };
in
{
    minecraft =  pkgs.makeDesktopItem {
        name = "minecraft";
        desktopName = "Shiginima Minecraft";
        type = "Application";
        exec = "${drv}/bin/minecraft $*";
    };
}