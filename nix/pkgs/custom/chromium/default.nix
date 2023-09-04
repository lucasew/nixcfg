{ ccacheStdenv
, chromium
, lib
}:

lib.pipe chromium [
  (drv: drv.override { stdenv = ccacheStdenv; })
  (drv: drv.overrideAttrs (old: {
    patches = (old.patches or []) ++ [ ./hide-tabs.patch ];
  }))
]
