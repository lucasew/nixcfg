{...}:
let
    common = import ./../common;
    selected = common.selectedDesktopEnvironment;
    selectedPath = ./. + "/${selected}.nix";
in {
    imports = [selectedPath];
}
