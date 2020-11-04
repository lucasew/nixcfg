self: super: 
let
  dotenvGit = builtins.fetchGit {
    url = "https://github.com/lucasew/dotenv";
  };
in {
  dotenv = import "${dotenvGit}";
}
