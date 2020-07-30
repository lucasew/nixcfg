#! /usr/bin/env bash

export ROOTFS="" # /mnt in nix installer
export USERNAME=lucasew
export MACHINE=$1 ;shift

[ -z "$MACHINE" ] && echo "Empty machine" && exit 255
[ ! -d "$(pwd)/machine/$MACHINE" ] && echo "No such machine: $MACHINE" && exit 255

echo "Deploying ${USERNAME}@${MACHINE}..."

mkdir -p $ROOTFS/home/$USERNAME/.config
ln -sfn  $(pwd)/user/ $ROOTFS/home/$USERNAME/.config/nixpkgs
ln -sfn $(pwd)/common/ $(pwd)/user/common
sudo ln -sfn $(pwd)/machine/$MACHINE/ $ROOTFS/etc/nixos

