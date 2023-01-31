{ self, config, pkgs, lib, ... }:
let
  inherit (lib) concatStringsSep attrValues mapAttrs;

  mkDate = dateStr:
    let
      dateChars = lib.stringToCharacters dateStr;
      step = value: stepVal:
      if builtins.typeOf stepVal == "string" then
        (step (value + stepVal))
      else if builtins.typeOf stepVal == "int" then
        (step (value + (builtins.elemAt dateChars stepVal)))
      else value;
    in step "";

  mkInput = inputName:
  let
    input = self.inputs.${inputName};
    revDate = if (input.sourceInfo or null) != null then
      mkDate input.sourceInfo.lastModifiedDate 0 1 2 3 "/" 4 5 "/" 6 7 " " 8 9 ":" 10 11 ":" 12 13 null
      else "unknown";
    fullRev = "${inputName}@${input.shortRev} (${revDate})";
  in ''<span class="btn btn-light"><b>${inputName}</b> <span class="hidden-part">${input.sourceInfo.lastModifiedDate or "unknown"}-${input.shortRev}</span></span>'';

  template = ''
<!DOCTYPE html>
  <html>
    <head>
      <meta charset="utf-8">
      <title>${config.networking.hostName}</title>
      <style>
      :root {
        ${concatStringsSep "\n" (attrValues (mapAttrs 
          (k: v: ''
            --var-${k}: #${v};
          '') (pkgs.custom.colors.colors)
          ))}
          --bs-body-color: var(--base00);
          --bs-body-bg: var(--base05);
          --bs-secondary-color: var(--base01);
          --bs-secondary-bg: var(--base06);
          --bs-tertiary-color: var(--base02);
          --bs-tertiary-bg: var(--base07);
          --bs-emphasis-color: var(--base01);
          --bs-border-color: var(--base04);
          --bs-primary: var(--base00);
          --bs-primary-bg-subtle: var(--base0D);
          --bs-primary-border-subtle: var(--base0C);
          --bs-primary-text: var(--base00);
          --bs-success-bg-subtle: var(--base0B);
          --bs-danger-bg-subtle: var(--base08);
          --bs-warning-bg-subtle: var(--base0A);
          --bs-info-bg-subtle: var(base0D);
        }
            a:hover > .hidden-part, span:hover > .hidden-part {
              display: inherit;
            }
            a:not(:hover) > .hidden-part, span:not(:hover) > .hidden-part {
              display: none;
            }
      </style>
      <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha1/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-GLhlTQ8iRABdZLl6O3oVMWSktQOp6b7In1Zl3/Jr59b6EGGoI1aFkw7cmDA6j6gD" crossorigin="anonymous">
    </head>

    <body class="mx-auto" style="max-width: max(80vw, 20rem)">
      <section id="hello" class="my-1 d-flex flex-row align-items-center justify-content-center">
        <img style="height: 4rem; width: auto;" src="/nix-logo.png">
        <h1 style="font-size: 4rem;">${config.networking.hostName}</h1>
      </section>
      <section id="nginx" class="my-1">
        <h2>Nginx hosts</h2>
        ${concatStringsSep "\n" (attrValues (mapAttrs
          (k: v: ''
            <a class="btn btn-light" target="_blank" href="http://${k}"><b>${k}</b></a>
          '') (config.services.nginx.virtualHosts)
        ))}
      </section>

      <section id="versions" class="my-1 flex">
        <h2>Inputs</h2>
          <span class="btn btn-light"><b>nixcfg</b> <span class="hidden-part">${self.shortRev}  (${mkDate self.sourceInfo.lastModifiedDate 0 1 2 3 "/" 4 5 "/" 6 7 " " 8 9 ":" 10 11 ":" 12 13 null})</span></span>

          ${builtins.concatStringsSep " " (map (mkInput) (builtins.sort (a: b: a < b)(builtins.attrNames self.inputs)))}
      </section>


    </body>

  </html>
  '';
in
{
  environment.etc."rootdomain/index.html".source = pkgs.writeText "template.html" template;
  environment.etc."rootdomain/favicon.ico".source = pkgs.fetchurl {
    url = "https://nixos.org/favicon.ico";
    sha256 = "sha256-59/+37K66dD+ySkvZ5JS/+CyeC2fDD7UAr1uiShxwYM=";
  };
  environment.etc."rootdomain/nix-logo.png".source = "${pkgs.nixos-icons}/share/icons/hicolor/1024x1024/apps/nix-snowflake.png";

  services.nginx.virtualHosts."${config.networking.hostName}.${config.networking.domain}" = {
    locations."/".root = "/etc/rootdomain";
  };
}
