self: super: 
{
  usb_tixati = super.pkgs.buildFHSUserEnv {
    name = "usb_tixati";
    targetPkgs = pkgs: with super.pkgs; [
      glib
      zlib
      dbus
      dbus-glib
      gtk2
      gdk-pixbuf
      cairo
      pango
    ];
    runScript = "/run/media/lucasew/Dados/PortableApps/PROGRAMAS/Tixati_portable/tixati_Linux64bit";
  }; 
}
