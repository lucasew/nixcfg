{ lib, buildGoModule }:

buildGoModule {
  pname = "workspaced";
  version = "0.0.1";

  src = ./.;

  # vendorHash = lib.fakeHash; # update this after first failed build
  vendorHash = "sha256-EuXLV+pBZxKrPjqyUYXkI9dJNeTIDvPVYOTO+CWr0mc=";

  meta = with lib; {
    description = "Workspace manager daemon";
    license = licenses.mit;
  };
}
