self: super:
with super;
let
    launcherZip = pkgs.requireFile {
        name = "ShiginimaSE_v4400.zip";
        sha1 = "61cb768106e6e449158ebb2608ad1327402d9fec";
        url = "https://teamshiginima.com/update/";
    };
    envLibPath = with pkgs; stdenv.lib.makeLibraryPath [
        curl
        libpulseaudio
        systemd
        alsaLib # needed for narrator
        flite # needed for narrator
        xorg.libXxf86vm # needed only for versions <1.13
        xorg.libX11
        libGLU libGL xorg.libXpm xorg.libXext alsaLib
    ];
    libPath = stdenv.lib.makeLibraryPath (with pkgs;[
        alsaLib
        atk
        cairo
        cups
        dbus
        expat
        fontconfig
        freetype
        gdk-pixbuf
        glib
        gnome2.GConf
        gnome2.pango
        gtk3-x11
        gtk2-x11
        nspr
        nss
        stdenv.cc.cc
        zlib
        libuuid
  ] ++
  (with xorg; [
    libX11
    libxcb
    libXcomposite
    libXcursor
    libXdamage
    libXext
    libXfixes
    libXi
    libXrandr
    libXrender
    libXtst
    libXScrnSaver
  ]));
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
            # makeWrapper ${pkgs.latest.oraclejre}/bin/java $out/bin/minecraft \
            makeWrapper ${pkgs.jre8}/bin/java $out/bin/minecraft \
                --add-flags "-jar $out/share/java/minecraft.jar" \
                      --prefix LD_LIBRARY_PATH : ${envLibPath}
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
