{
  description = "Deep Learning Project";

  inputs = {
    # TODO: Pin to specific commit
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url =
      "github:NixOS/nixpkgs/d934204a0f8d9198e1e4515dd6fec76a139c87f0";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = nixpkgs.legacyPackages.${system};
      in {
        packages.default = pkgs.buildGo122Module {
          name = "elf";
          version = "0.0.1";
          src = self;
          vendorHash = null;
        };

        devShells.default =
          pkgs.mkShell { packages = with pkgs; [ go_1_22 bruno ]; };
      });
}

