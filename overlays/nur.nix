self: super:
let
  pkgs = super.pkgs;
  nurRepo = builtins.fetchTarball {
    # Get the revision by choosing a version from https://github.com/nix-community/NUR/commits/master
    url = "https://github.com/nix-community/NUR/archive/67fb3de1cf678b614cc618cbf9e221361bf1dd0c.tar.gz";
    # Get the hash by running `nix-prefetch-url --unpack <url>` on the above url
    sha256 = "04387gzgl8y555b3lkz9aiw9xsldfg4zmzp930m62qw8zbrvrshd";
  };
in
{
  nur = import nurRepo {
    inherit pkgs;
  };
}
