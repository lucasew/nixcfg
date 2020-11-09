{
  description,
  command,
  enable ? true
}:
{
  Unit = {
    Description = description;
    PartOf = ["graphical-session.target"];
  };
  Service = {
    Type = "exec";
    ExecStart = command;
    Restart = "on-failure";
  };
  Install = {
    WantedBy = if enable then [
      "default.target"
    ] else [];
  };
}
