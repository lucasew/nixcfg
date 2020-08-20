self: super: {
  dotwrap = super.writeScriptBin "dotwrap" ''
    #!/usr/bin/env bash
    ${super.globalConfig.setupScript}

    [ "$1" == pushd ] && cd "${super.globalConfig.dotfileRootPath}" && shift

    $*

  '';
}
