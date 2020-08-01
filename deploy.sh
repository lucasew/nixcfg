#! /usr/bin/env bash

export USERNAME=lucasew
export SELECTED_DE=xfce
export HOSTNAME=acer-nix

export COMMAND=$1; shift

case "$COMMAND" in 
    "install") export ROOTFS="/mnt" ;; # install config files to /mnt
    "apply") export ROOTFS="" ;; # install config files to /
    *) echo "No such command. Supported: install apply" && exit 255 ;;
esac

export MACHINE=$1 ;shift

[ -z "$MACHINE" ] && echo "Empty machine" && exit 255
[ ! -d "$(pwd)/machine/$MACHINE" ] && echo "No such machine: $MACHINE" && exit 255

echo "Deploying ${USERNAME}@${MACHINE}..."

mkdir -p $ROOTFS/home/$USERNAME/.config
ln -sfn  $(pwd)/user/ $ROOTFS/home/$USERNAME/.config/nixpkgs
ln -sfn $(pwd)/common/ $(pwd)/user/common
ln -sfn $(pwd)/common/ $(pwd)/machine/$MACHINE/common
sudo ln -sfn $(pwd)/machine/$MACHINE/ $ROOTFS/etc/nixos

cat $(pwd)/common/default.nix.example | \
    sed s/%USERNAME%/$USERNAME/ | \
    sed s/%DE%/$SELECTED_DE/ | \
    sed s/%HOSTNAME%/$HOSTNAME/ \
    > $(pwd)/common/default.nix

echo "The config files are where they should be. It's time to let nix do the rest"
echo "All the configs are pointing to $(pwd)." 
echo "If you need to change the files to another location run me again as you are running me now"

