self: super:
let
  nixgramGit = builtins.fetchGit {
    url = "https://github.com/lucasew/nixgram";
  };
in {
  nixgram = import "${nixgramGit}";
}
