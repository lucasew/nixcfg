{...}:
let
  nixgramGit = builtins.fetchGit {
    url = "https://github.com/lucasew/nixgram";
  };
in import "${nixgramGit}"
