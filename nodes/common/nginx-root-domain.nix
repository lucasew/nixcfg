{ self, config, pkgs, lib, ... }:

with pkgs.custom.colors.colors;
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
  in ''<div><b>${inputName}</b> <span class="hidden-part">${input.sourceInfo.lastModifiedDate or "unknown"}-${input.shortRev}</span></div>'';

  template = ''
<!DOCTYPE html>
  <html>
    <head>
      <meta charset="utf-8">
      <title>${config.networking.hostName}</title>
      <style>
:root{
  --color-blossom: #1d7484;
  --color-fade: #982c61;

  --color-bg: #${base00};
  --color-bg-alt: #${base01};

  --color-text: #${base05};
  --font-size-base: 1.8rem;

  --font-family-base: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, "Noto Sans", sans-serif;
  --font-family-heading: --font-family-base;
}

/* Body */
html {
  font-size: 62.5%; / So that root size becomes 10px*/
  font-family: var(--font-family-base);
}

body {
  /* var(--font-size-base must be a rem value */
  font-size: var(--font-size-base);
  line-height: 1.618;
  max-width: 38em;
  margin: auto;
  color: var(--color-text);
  background-color: var(--color-bg);
  padding: 13px;
}

@media (max-width: 684px) {
  body {
  font-size: var(--font-size-base * 0.85);
  }
}

@media (max-width: 382px) {
  body {
  font-size: var(--font-size-base * 0.75);
  }
}

/* added Nav */
nav ul {
  list-style-type: none;
  margin: 0;
  padding: 0;
  overflow: hidden;
  background-color: #666;
  position: fixed;
  top: 0;
  width: 100%;
}

nav ul li {
  float: left;
  border-right: 1px solid #bbb;
}

nav ul li a {
  display: block;
  color: white;
  text-align: center;
  padding: 14px 16px;
  text-decoration: none;
}

/* Change the link color to #111 (black) on hover */
nav ul li a:hover {
  background-color: #111;
}

:root {
  --word-wrap: {
    overflow-wrap: break-word;
    word-wrap: break-word;
    -ms-word-break: break-all;
    word-break: break-word;
  }
}

h1, h2, h3, h4, h5, h6 {
  line-height: 1.1;
  font-family: var(--font-family-heading);
  font-weight: 700;
  margin-top: 3rem;
  margin-bottom: 1.5rem;
  @apply --word-wrap;
}

h1 { font-size: 2.35em }
h2 { font-size: 2.00em }
h3 { font-size: 1.75em }
h4 { font-size: 1.5em }
h5 { font-size: 1.25em }
h6 { font-size: 1em }

p {
  margin-top: 0px;
  margin-bottom: 2.5rem;
}

small, sub, sup {
  font-size: 75%;
}

hr {
  border-color: var(--color-blossom);
}

a {
  text-decoration: none;
  color: var(--color-blossom);

  & hover {
      color: var(--color-fade);
      border-bottom: 2px solid var(--color-text);
  }

  & visited {
      color: darken(var(--color-blossom, 10%));
  }

}

ul {
  padding-left: 1.4em;
  margin-top: 0px;
  margin-bottom: 2.5rem;
}

li {
  margin-bottom: 0.4em;
}

blockquote {
  margin-left: 0px;
  margin-right: 0px;
  padding-left: 1em;
  padding-top: 0.8em;
  padding-bottom: 0.8em;
  padding-right: 0.8em;
  border-left: 5px solid var(--color-blossom);
  margin-bottom: 2.5rem;
  background-color: var(--color-bg-alt);
}

blockquote p {
  margin-bottom: 0;
}

img, video {
  height: auto;
  max-width: 100%;
  margin-top: 0px;
  margin-bottom: 2.5rem;
}

/* Pre and Code */
pre {
  background-color: var(--color-bg-alt);
  display: block;
  padding: 1em;
  overflow-x: auto;
  margin-top: 0px;
  margin-bottom: 2.5rem;
}

code {
  font-size: 0.9em;
  padding: 0 0.5em;
  background-color: var(--color-bg-alt);
  white-space: pre-wrap;
}

pre > code {
  padding: 0;
  background-color: transparent;
  white-space: pre;
}

/* Tables */
table {
  text-align: justify;
  width: 100%;
  border-collapse: collapse;
}

td, th {
  padding: 0.5em;
  border-bottom: 1px solid var(--color-bg-alt);
}

/* Buttons, forms and input */
input, textarea {
  border: 1px solid var(--color-text);

  & focus {
      border: 1px solid var(--color-blossom);
  }

}

textarea {
  width: 100%;
}

.button, button, input[type="submit"], input[type="reset"], input[type="button"] {
  display: inline-block;
  padding: 5px 10px;
  text-align: center;
  text-decoration: none;
  white-space: nowrap;

  background-color: var(--color-blossom);
  color: var(--color-bg;
  border-radius: 1px;
  border: 1px solid var(--color-blossom);
  cursor: pointer;http://155.138.194.207:8080/sty/
  box-sizing: border-box;

  &[disabled] {
      cursor: default;
      opacity: .5;
  }

  & focus:enabled, & hover:enabled {
      background-color: var(--color-fade);
      border-color: var(--color-fade);
      color: var(--color-bg);
      outline: 0;
  }

}

textarea, select, input {
  color: var(--color-text);
  padding: 6px 10px; /* The 6px vertically centers text on FF, ignored by Webkit */
  margin-bottom: 10px;
  background-color: var(--color-bg-alt);
  border: 1px solid var(--color-bg-alt);
  border-radius: 4px;
  box-shadow: none;
  box-sizing: border-box;

  & focus {
      border: 1px solid var(--color-blossom);
      outline: 0;
  }

}

input[type="checkbox"]:focus {
  outline: 1px dotted var(--color-blossom);
}

label, legend, fieldset {
  display: block;
  margin-bottom: .5rem;
  font-weight: 600;
}

div:hover .hidden-part {
  display: inherit;
}

div:not(:hover) .hidden-part {
  display: none;
}

.small-cards-container {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-around;
  flex-direction: row;
}

.small-cards-container > * {
  display: inline-block;
}

section#hello {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
}

section#hello > * {
  margin: 0;
  padding-left: 1rem;
  display: inline-block;
}

a {
  font-weight: bold;
}

      </style>
    </head>

    <body>
      <section id="hello">
        <img style="height: 4rem; width: auto;" src="/nix-logo.png">
        <h1 style="font-size: 4rem;">${config.networking.hostName}</h1>
      </section>
      <section id="nginx">
        <h2>Nginx hosts</h2>
        <div class="small-cards-container">
        ${concatStringsSep "\n" (attrValues (mapAttrs
          (k: v: ''
            <a class="btn btn-light" target="_blank" href="http://${k}">${k}</a>
          '') (config.services.nginx.virtualHosts)
        ))}
        </div>
      </section>

      <section id="versions">
        <h2>Inputs</h2><br>
            <div class="small-cards-container">
              <div><b>nixcfg</b> <span class="hidden-part">${self.shortRev}  (${mkDate self.sourceInfo.lastModifiedDate 0 1 2 3 "/" 4 5 "/" 6 7 " " 8 9 ":" 10 11 ":" 12 13 null})</span></div>

              ${builtins.concatStringsSep " " (map (mkInput) (builtins.sort (a: b: a < b)(builtins.attrNames self.inputs)))}
            </div>
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
