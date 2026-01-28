{ lib, buildGoModule }:

buildGoModule {
  pname = "workspaced";
  version = "0.0.1";

  src = ./.;

  # vendorHash = lib.fakeHash; # update this after first failed build
  vendorHash = "sha256-qqCV2U3qd24QTW2EsCb7nn5Ulg+ffUspNZkvGrGhffU=";

  meta = with lib; {
    description = "Workspace manager daemon";
    license = licenses.mit;
  };
}
