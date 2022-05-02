{ stdenv
, unity3d
, fetchFromGitHub
, xvfb-run
, lib
}:
stdenv.mkDerivation {
  name = "Geographical-Adventures";
  version = "unstable-2022-5-2";
  src = fetchFromGitHub {
    owner = "SebLague";
    repo = "Geographical-Adventures";
    sha256 = "sha256-DXziwS9KuJoSao/IQ70kGCFtjaWgtHPCwrcPdrMY5AE=";
    rev = "c611e3455a012b8838faa81e47351a7fcfa2e449";
  };
  buildInputs = [
    unity3d
    xvfb-run
  ];
  installPhase = ''
    xvfb-run unity-editor -quit -batchmode -projectPath "$(pwd)" -executeMethod UnityBuilderAction.Builder.BuildProject -customBuildPath $out -buildTarget LinuxStandalone -logfile /dev/stdout
  '';
}
