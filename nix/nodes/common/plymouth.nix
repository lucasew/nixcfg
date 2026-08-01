{ url, sha256 }:
{ pkgs, ... }:
{
  boot.plymouth = {
    enable = true;
    theme = "breeze";
    logo = pkgs.plymouthSvgLogo {
      inherit url sha256;
    };
  };
}
