{
  pkgs,
  lib,
  config,
  ...
}:
{
  services.tailscale.enable = lib.mkDefault true;

  networking.firewall = lib.mkIf config.services.tailscale.enable {
    trustedInterfaces = [ "tailscale0" ];

    # allow the Tailscale UDP port through the firewall
    allowedUDPPorts = [ config.services.tailscale.port ];
  };

  # create a oneshot job to authenticate to Tailscale
  systemd.services.tailscale-autoconnect = lib.mkIf config.services.tailscale.enable {
    description = "Automatic connection to Tailscale";

    # make sure tailscale is running and the network is online
    after = [ "network-online.target" "tailscale.service" ];
    wants = [ "network-online.target" "tailscale.service" ];
    wantedBy = [ "multi-user.target" ];

    # set this service as a oneshot job
    serviceConfig.Type = "oneshot";
    serviceConfig.Restart = "on-failure";
    serviceConfig.RestartSec = "5s";


    # have the job run this shell script
    script = with pkgs; ''
      # Retry logic
      for i in $(seq 1 5); do
        # check if we are already authenticated to tailscale
        status="$(${tailscale}/bin/tailscale status -json | ${jq}/bin/jq -r .BackendState)"
        if [ $status = "Running" ]; then # if so, then do nothing
          exit 0
        fi

        # otherwise authenticate with tailscale
        #
        # IMPORTANT: This service requires a Tailscale auth key to be securely
        # provisioned using sops-nix. The NixOS build will fail until this is set up.
        #
        # 1. Create a file `secrets/tailscale-auth-key.yaml` with the content:
        #    TAILSCALE_AUTHKEY: "your-tskey-auth-..."
        #
        # 2. Encrypt the file in-place by running `sops encrypt -i secrets/tailscale-auth-key.yaml`
        #    in an environment with the `sops` command available.
        #
        # 3. Add the following to your `nix/nodes/common/sops.nix` file:
        #    sops.secrets."tailscale-authkey" = {
        #      sopsFile = ../../../secrets/tailscale-auth-key.yaml;
        #      key = "TAILSCALE_AUTHKEY"; # This line is important
        #    };
        #
        # This script will then use the provisioned key.
        ${tailscale}/bin/tailscale up -authkey "$(cat ${config.sops.secrets."tailscale-authkey".path})"

        # Check status again
        status="$(${tailscale}/bin/tailscale status -json | ${jq}/bin/jq -r .BackendState)"
        if [ $status = "Running" ]; then
            exit 0
        fi

        # If connection fails, restart the daemon and try again
        ${systemd}/bin/systemctl restart tailscale.service
        sleep 5
      done
      exit 1
    '';
  };
}
