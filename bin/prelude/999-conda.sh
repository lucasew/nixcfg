set +h # FIX: bash: hash: hashing está desabilitado

function loadConda {
    export PATH="$PATH:$HOME/.conda/bin"

    # >>> conda initialize >>>
    # !! Contents within this block are managed by 'conda init' !!
    __conda_setup="$('/home/lucasew/.conda/bin/conda' 'shell.bash' 'hook' 2> /dev/null)"
    if [ $? -eq 0 ]; then
        eval "$__conda_setup"
    else
        if [ -f "/home/lucasew/.conda/etc/profile.d/conda.sh" ]; then
            . "/home/lucasew/.conda/etc/profile.d/conda.sh"
        else
            export PATH="/home/lucasew/.conda/bin:$PATH"
        fi
    fi
    unset __conda_setup
    # <<< conda initialize <<<

}
