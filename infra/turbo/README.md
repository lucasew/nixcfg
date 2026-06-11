# Turbo Infrastructure Workspace

This directory contains the Terraform configuration (`gcp.tf`) for provisioning an ephemeral compute instance ("turbo") on Google Cloud Platform (GCP). It is designed to quickly spin up a machine (optionally with GPU acceleration for ML tasks) and spin it down when no longer needed to save costs.

## Examples

**1. Start a standard instance (50GB disk):**
This provisions a standard `e2-micro` instance by default.
```bash
terraform apply -var gcp_turbo_disk_size=50 -var gcp_turbo_stop=false
```

**2. Start a "Turbo" instance (50GB disk, high CPU, GPU enabled):**
Setting `gcp_turbo_modo_turbo=true` changes the machine type to `n1-highcpu-4` and attaches an NVIDIA Tesla T4 GPU (based on `gcp.tf` defaults).
```bash
terraform apply -var gcp_turbo_disk_size=50 -var gcp_turbo_stop=false -var gcp_turbo_modo_turbo=true
```

**3. Stop the instance:**
This scales the instance count to 0, effectively stopping it to prevent further billing, while allowing you to keep the terraform state. It also changes the base image to NixOS.
```bash
terraform apply -var gcp_turbo_disk_size=20 -var gcp_turbo_stop=true -var gcp_instance_image=nixos -var gcp_turbo_modo_turbo=false
```
