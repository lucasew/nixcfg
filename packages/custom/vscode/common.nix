{pkgs ? import <nixpkgs> {}, ...}:
{
  identifier = pkgs.lib.mkDefault "common";
  extensions = [
    {
      publisher = "MS-CEINTL";
      name = "vscode-language-pack-pt-BR";
      version = "1.61.4";
      sha256 = "sha256-qR7UD5SZBfW3Ihly0eP9ZHYMrKgxR4iaperP+jpI82s=";
    }
    {
      publisher = "arrterian";
      name = "nix-env-selector";
      version = "0.1.2";
      sha256 = "1n5ilw1k29km9b0yzfd32m8gvwa2xhh6156d4dys6l8sbfpp2cv9";
    }
    {
      publisher = "CoenraadS";
      name = "bracket-pair-colorizer";
      version = "1.0.61";
      sha256 = "0r3bfp8kvhf9zpbiil7acx7zain26grk133f0r0syxqgml12i652";
    }
    {
      publisher = "donjayamanne";
      name = "githistory";
      version = "0.6.8";
      sha256 = "0wc0wsnqwyg0pz0jrmw0038k6g1p564krqscrx3h8wpyfymcd68l";
    }
    {
      publisher = "file-icons";
      name = "file-icons";
      version = "1.0.25";
      sha256 = "0s6lr7s1a0alkknazmch5k2m0r16p5gnlzn3yyan9wl8k3579c25";
    }
  ];
  settings = {
    "workbench.iconTheme" = "file-icons";
    "workbench.colorTheme" = "One Dark Pro";
  };
}
