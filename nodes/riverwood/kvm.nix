{self, global, lib, pkgs, ...}:
let
  inherit (global) username;
in
{
  users.users.${username} = {
    extraGroups = [ "kvm" "libvirt" ];
  };
  virtualisation = {
    kvmgt.enable = true;
    libvirtd.enable = true;
  };
  environment.systemPackages = with pkgs; [
    virt-manager
  ];
}
