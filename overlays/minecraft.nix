self: super:
with super;
let
  launcherZip = pkgs.requireFile {
    name = "ShiginimaSE_v4400.zip";
    sha1 = "61cb768106e6e449158ebb2608ad1327402d9fec";
    url = "https://teamshiginima.com/update/";
  };
  envLibPath = with pkgs; stdenv.lib.makeLibraryPath [
    alsaLib # needed for narrator
    curl
    flite # needed for narrator
    libGL
    libGLU
    libpulseaudio
    systemd
    xorg.libX11
    xorg.libXext
    xorg.libXpm
    xorg.libXxf86vm # needed only for versions <1.13
  ];
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
      # makeWrapper ${pkgs.oraclejre}/bin/java $out/bin/minecraft \
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
  minecraft = pkgs.makeDesktopItem {
    name = "minecraft";
    desktopName = "Shiginima Minecraft";
    type = "Application";
    icon = pkgs.fetch "https://icons.iconarchive.com/icons/blackvariant/button-ui-requests-2/1024/Minecraft-2-icon.png";
    exec = "${drv}/bin/minecraft $*";
  };
}
