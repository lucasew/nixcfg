{pkgs, ...}:
{
    allowUnfree = true;
    packageOverrides = import ./common/apps;
}

