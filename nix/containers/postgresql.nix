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
      # Esta configuração é para o pg_hba.conf
      authentication = ''
        host all all 0.0.0.0/0 md5
      '';
    };
    # A imagem precisa de uma rede funcional e de um nome de host
    networking.hostName = "postgresql";
    # Abrir a porta do postgresql no firewall *dentro* do contêiner
    networking.firewall.allowedTCPPorts = [ 5432 ];
  };

  # 2. Construir a imagem OCI
  postgresql-image = pkgs.dockerTools.buildImage {
    name = "nixos-postgresql";
    tag = "latest";
    config = {
      # O comando padrão para iniciar o contêiner NixOS
      Cmd = [ "${pkgs.nixos-container}/bin/nixos-container" "run-unconfigured-container" ];
    };
    # A configuração NixOS a ser usada para construir a imagem
    nixOSConfiguration = postgresql-nixos-config;
  };

in
{
  # 3. Definir o contêiner para usar a imagem construída
  virtualisation.oci-containers.containers.postgresql = {
    imageFile = postgresql-image;
    ports = [ "5432:5432" ];
    # Mapeia para um diretório de dados diferente no host para evitar conflitos com o postgresql do host
    volumes = [ "/var/lib/postgresql_container_data:/var/lib/postgresql" ];
  };
}
