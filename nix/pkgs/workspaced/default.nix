{ lib, buildGoModule }:

buildGoModule {
  pname = "workspaced";
  version = "0.0.1";

  src = ../../../bin/workspaced;

  # vendorHash = lib.fakeHash; # update this after first failed build
  vendorHash = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";

  meta = with lib; {
    description = "Workspace manager daemon";
    license = licenses.mit;
  };
}
