self: super:
with super;
let
  version = "1.3.18";
  appimage = builtins.fetchurl {
    url = "https://cdn.documentnode.net/stable/${version}/DocumentNode-${version}-x86_64.AppImage";
    sha256 = "0px27r0pamjc5dnnyh14800xsmpjnpy5rn24c87ilsqd5826ag16";
  };
  app = pkgs.stdenv.mkDerivation {
    pname = "DocumentNode";
    version = version;
    dontUnpack = true;
    installPhase = ''
      cp ${appimage} $out
      chmod +x $out
    '';
  };
  bin = pkgs.writeShellScriptBin "DocumentNode" ''
    ${pkgs.appimage-run}/bin/appimage-run ${app} $*
  '';
  desktop = pkgs.makeDesktopItem {
    name = "DocumentNode";
    desktopName = "Document Node";
    type = "Application";
    icon = builtins.fetchurl {
      url = "https://documentnode.io/images/documentnode_simple.svg";
      sha256 = "14rlby0d4bvq0760k51gibjlwshgmxnsf8f5f86ap2sadcn7c3q7";
    };
    exec = "${bin}/bin/DocumentNode";
  };
in
{
  documentNode = desktop;
}
