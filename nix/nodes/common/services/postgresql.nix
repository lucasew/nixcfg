{ lib, config, ... }:
{
  options = {
    services.postgresql.userSpecificDatabases = lib.mkOption {
      description = lib.mdDoc "Extra databases and users for specific reasons";

      type = with lib.types; attrsOf (listOf str);

      # $user_$database for all databases + $user so psql works out of the box 
      # one will still need to set passwords for users
      apply = lib.mapAttrs (k: v: (map (item: "${k}_${item}") v) ++ [k]);
    };

    services.postgresql.testDatabases = lib.mkOption {
      description = lib.mdDoc "Extra databases to be created and granted to the test user. Shouldn't contain production data and will be prefixed with `test_`";
      type = with lib.types; listOf str;
      default = [];
    };
  };
  config = lib.mkIf config.services.postgresql.enable {

    services.postgresql.userSpecificDatabases.test = config.services.postgresql.testDatabases;

    services.postgresqlBackup = {
      enable = true;
      databases = [ "postgres" ];
    };

    services.postgresql = {
      ensureDatabases = lib.flatten (lib.attrValues config.services.postgresql.userSpecificDatabases);
      ensureUsers = (map ({name, value}: {
        inherit name;
        ensurePermissions = lib.listToAttrs (map (database: {
          name = "DATABASE \"${database}\"";
          value = "ALL PRIVILEGES";
        }) value);
        ensureClauses = {
          login = true;
        };
        
      }) (lib.attrsToList config.services.postgresql.userSpecificDatabases));
    };
  };
}
