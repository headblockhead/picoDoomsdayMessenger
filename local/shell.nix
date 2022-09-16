{ pkgs ? import <nixpkgs> { } }:

with pkgs;

stdenv.mkDerivation {
  name = "local-build-enviroment";

  buildInputs = [ pkgs.gcc pkgs.pkg-config pkgs.zlib pkgs.gtk3 pkgs.mesa ];
}
