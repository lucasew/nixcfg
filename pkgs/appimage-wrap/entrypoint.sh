if [ $# == 0 ]; then
    echo "No AppImage provided"
    exit 1
fi

APPIMAGE="$1"; shift
PATH=$PATH:@fhs@/bin
OFFSET="$(appimage-env "$APPIMAGE" --appimage-offset | head -n 1)"
LOOP=$(udisksctl loop-setup -f "$APPIMAGE" --offset "$OFFSET" | sed 's;[ \.];\n;g' | grep '/dev/loop')
MOUNTPOINT=$(udisksctl mount -b $LOOP | sed 's;[ \.];\n;g' | grep '/media')
if [ ! -z "$MOUNTPOINT" ]; then
  appimage-env "$MOUNTPOINT/AppRun" "$@"
fi
if [ ! -z "$LOOP" ]; then
  while true; do
    if [[ ! "$(udisksctl unmount -b "$LOOP" 2>&1 || true)" =~ "Error" ]]; then
      break
    fi
    sleep 1
  done
  while true; do
    udisksctl loop-delete -b "$LOOP" && break || true
    sleep 1
  done
fi
