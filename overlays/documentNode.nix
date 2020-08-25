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

in
{
  documentNode = bin;
}
