{ config, pkgs, lib, ... }:

{
  config = {
    virtualisation.oci-containers.containers.postgresql = {
      image = "postgresql-image";
      ports = [ "5432:5432" ];
      volumes = [ "/var/lib/postgresql:/var/lib/postgresql" ];
      config = {
        system.stateVersion = "22.11";
        services.postgresql = {
          enable = true;
          package = pkgs.postgresql_14; # Usando uma versão específica
          ensureDatabases = [ "nextcloud" ];
          ensureUsers = [{
            name = "nextcloud";
            ensureDBOwnership = true;
          }];
        };
      };
    };
  };
}
