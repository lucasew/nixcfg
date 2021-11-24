{config, ...}:
{
  initEl.pre = ''
    (let
      (
        ;; Temporarily increase the GC threshold to avoid GCs on initialization
        (gc-cons-threshold most-positive-fixnum)
        ;; Avoid analyzing files when loading remote files
        (file-name-handler-alist nil))
  '';
  initEl.pos = ''
    )
  '';
}
