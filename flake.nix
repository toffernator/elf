{
  description = "Deep Learning Project";

  inputs = {
    # TODO: Pin to specific commit
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        # Place your own derivations here
      in {
        devShells.default = pkgs.mkShell { packages = with pkgs; [ just go ]; };
      });
}

