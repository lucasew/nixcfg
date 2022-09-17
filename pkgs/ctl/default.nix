{ lib, pkgs }:
lib.climod {
  imports = [
    ./deploy
    # ./options # upstream has broken stuff and builtins.tryEval only deals with assertions
  ];
  name = "ctl";
  description = "lucasew's control CLI";
}
