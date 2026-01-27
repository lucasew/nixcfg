{ lib, buildGoModule }:

buildGoModule {
  pname = "workspaced";
  version = "0.0.1";

  src = ../../../bin/workspaced;

  vendorHash = null;

  meta = with lib; {
    description = "Workspace manager daemon";
    license = licenses.mit;
  };
}
