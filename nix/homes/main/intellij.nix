{ pkgs, ... }:
{
  home.file = {
    ".var/app/com.jetbrains.IntelliJ-IDEA-Community/config/idea64.vmoptions".source = ../../../config/intellij/idea64.vmoptions;
    ".var/app/com.jetbrains.IntelliJ-IDEA-Ultimate/config/idea64.vmoptions".source = ../../../config/intellij/idea64.vmoptions;
  };
}
