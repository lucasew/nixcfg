{
  lib,
  stdenvNoCC,
  fetchFromGitHub,
  makeWrapper,
  bash,
  coreutils,
  findutils,
  gawk,
  git,
  gnugrep,
  gnused,
  gum,
  jq,
  procps,
  quickshell,
  systemd,
}:

stdenvNoCC.mkDerivation (finalAttrs: {
  pname = "omarchy";
  version = "4.0.0-unstable-2026-08-15";

  src = fetchFromGitHub {
    owner = "basecamp";
    repo = "omarchy";
    rev = "f32ebbdb730c4e8fe11e4046cef4267e466264ea";
    hash = "sha256-uBoJVjt0l/iERvd48GJdEER+fvYfvZL70MetWYcj/TM=";
  };

  nativeBuildInputs = [ makeWrapper ];

  dontConfigure = true;
  dontBuild = true;

  # Upstream launcher only probes Hyprland. Keep hyprctl, also accept Sway.
  postPatch = ''
    substituteInPlace bin/omarchy-launch-shell \
      --replace-fail 'hyprctl -j monitors >/dev/null 2>&1 && return 0' \
      'if command -v hyprctl >/dev/null && hyprctl -j monitors >/dev/null 2>&1; then return 0; fi
        if command -v swaymsg >/dev/null && swaymsg -t get_outputs >/dev/null 2>&1; then return 0; fi'
  '';

  installPhase = ''
    runHook preInstall

    mkdir -p $out/share/omarchy $out/bin
    cp -a bin shell default themes version $out/share/omarchy/
    for extra in icon.png icon.txt logo.svg logo.txt LICENSE; do
      if [ -e "$extra" ]; then
        cp -a "$extra" $out/share/omarchy/
      fi
    done
    if [ -d config ]; then
      cp -a config $out/share/omarchy/
    fi
    if [ -d applications ]; then
      cp -a applications $out/share/omarchy/
    fi

    runtime=${lib.makeBinPath [
      bash
      coreutils
      findutils
      gawk
      git
      gnugrep
      gnused
      gum
      jq
      procps
      quickshell
      systemd
    ]}

    for src in $out/share/omarchy/bin/omarchy $out/share/omarchy/bin/omarchy-*; do
      [ -f "$src" ] && [ -x "$src" ] || continue
      name=$(basename "$src")
      makeWrapper "$src" "$out/bin/$name" \
        --set OMARCHY_PATH "$out/share/omarchy" \
        --prefix PATH : "$out/bin:$runtime"
    done

    runHook postInstall
  '';

  meta = {
    description = "Omarchy Quattro CLI, plugin manager, and Quickshell desktop host";
    homepage = "https://github.com/basecamp/omarchy";
    license = lib.licenses.mit;
    platforms = lib.platforms.linux;
    mainProgram = "omarchy";
  };
})
