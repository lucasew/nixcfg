{pkgs, ...}: 
{
  virtualisation.anbox = {
    enable = true;
    # image = import ../../../utils/anboxImage.nix;
  };
}
