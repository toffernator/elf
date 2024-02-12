{
  description = "Deep Learning Project";

  inputs = {
    # TODO: Pin to specific commit
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url =
      "github:NixOS/nixpkgs/d934204a0f8d9198e1e4515dd6fec76a139c87f0";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = nixpkgs.legacyPackages.${system};
      in {
        devShells.default =
          pkgs.mkShell { packages = with pkgs; [ just go_1_22 bruno ]; };
      });
}

