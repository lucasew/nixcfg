{ ... }:

{
  programs.ghostty = {
    enableBashIntegration = true;
    settings = {
      window-decoration = false;
      theme = "base16-custom";
      cursor-style = "block";
    };
  };
}
