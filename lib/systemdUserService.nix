{description, command}:
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
    WantedBy = [
      "default.target"
    ];
  };
}
