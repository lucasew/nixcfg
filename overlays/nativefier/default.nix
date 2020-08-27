self: super:
{
  stdenv = super.stdenv // {
    mkNativefier = import ./mkNativefier.nix;
  };
}
