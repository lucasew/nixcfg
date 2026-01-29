{ lib, buildGoModule }:

buildGoModule {
  pname = "workspaced";
  version = "0.0.1";

  src = ./.;

  vendorHash = "sha256-4gQjC18gSf/rnjUq6L161dTUkxqUlf79fzG8+QOR4B4=";

  meta = with lib; {
    description = "Workspace manager daemon";
    license = licenses.mit;
  };
}
