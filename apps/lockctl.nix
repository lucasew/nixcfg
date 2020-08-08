{pkgs, config}: {
    pkgs.writeShellScriptBin "lockctl" ''
        export COMMAND=$1; shift
        case "$COMMAND" in
            "lock") 
            xset dpms force off
            xautolock -locknow
            ;;
        esac
    ''
}