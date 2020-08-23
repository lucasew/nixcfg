self: super:
let
  pkgs = super.pkgs;
  nurRepo = builtins.fetchTarball {
    url = "https://github.com/nix-community/NUR/archive/67fb3de1cf678b614cc618cbf9e221361bf1dd0c.tar.gz";
    sha256 = "15jkyjwllmzgclg4y3fq0lam0l9jm99idl8c9pjs6dm1vkdsbajn";
  };
in
{
  nur = import nurRepo {
    inherit pkgs;
  };
}
