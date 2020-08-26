self: super: {
    utils = {
        importAllIn = import ./importAllIn.nix;
        lsName = import ./lsName.nix;
        image2color = import ./image2color.nix;
    };
}