#! /usr/bin/env bash

export USERNAME=lucasew

export COMMAND=$1; shift

case "$COMMAND" in 
    "install") export ROOTFS="/mnt" ;; # install config files to /mnt
    "apply") export ROOTFS="" ;; # install config files to /
    *) echo "No such command. Supported: install apply" && exit 255 ;;
esac

export MACHINE=$1 ;shift

[ -z "$MACHINE" ] && echo "Empty machine" && exit 255
[ ! -d "$(pwd)/machine/$MACHINE" ] && echo "No such machine: $MACHINE" && exit 255

deploy_file() {
    dotfile_dir=$1; shift
    item=$1; shift
    destination_dir=$1; shift
    destination_file=$1; shift

    mkdir -p $destination_dir
    echo "let items = import $dotfile_dir; in items.$item" > /tmp/nixtemp
    mv /tmp/nixtemp $ROOTFS/$destination_dir/$destination_file
}

echo "Deploying ${USERNAME}@${MACHINE}..."

deploy_file "$(pwd)" "home" "/home/$USERNAME/.config/nixpkgs" "home.nix"
deploy_file "$(pwd)" "homeConfig" "/home/$USERNAME/.config/nixpkgs" "config.nix"
deploy_file "$(pwd)" "machine" "/etc/nixos" "configuration.nix"

echo "The config files are where they should be. It's time to let nix do the rest"
echo "All the configs are pointing to $(pwd)." 
echo "If you need to change the files to another location run me again as you are running me now"

