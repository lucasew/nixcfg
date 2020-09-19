self: super: {
  fetch = url: builtins.fetchurl {url = url;};
}
