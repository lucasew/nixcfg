{ config, pkgs, lib, ... }:

let
  # 1. Definir a configuração NixOS para a imagem
  postgresql-nixos-config = {
    system.stateVersion = "22.11";
    services.postgresql = {
      enable = true;
      package = pkgs.postgresql_14;
      ensureDatabases = [ "nextcloud" ];
      ensureUsers = [{
        name = "nextcloud";
        ensureDBOwnership = true;
      }];
      # É importante permitir conexões TCP/IP para que o contêiner do Nextcloud possa se conectar
      authentication = ''
        host all all 0.0.0.0/0 md5
      '';
    };
    networking.hostName = "postgresql";
    # Abrir a porta do postgresql no firewall do contêiner
    networking.firewall.allowedTCPPorts = [ 5432 ];
  };

  # 2. Construir a imagem OCI
  postgresql-image = pkgs.dockerTools.buildImage {
    name = "nixos-postgresql";
    tag = "latest";
    config = {
      Cmd = [ "${pkgs.nixos-container}/bin/nixos-container" "run-unconfigured-container" ];
    };
    nixOSConfiguration = postgresql-nixos-config;
  };

in
{
  # 3. Definir o contêiner para usar a imagem construída
  virtualisation.oci-containers.containers.postgresql = {
    imageFile = postgresql-image;
    ports = [ "5432:5432" ];
    volumes = [ "/var/lib/postgresql_data:/var/lib/postgresql" ]; # Mapeia para um diretório de dados diferente no host para evitar conflitos
  };
}
