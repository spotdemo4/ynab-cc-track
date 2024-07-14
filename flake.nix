{
  description = "A Nix-flake-based go development environment";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    templ.url = "github:a-h/templ";
  };

  outputs = { nixpkgs, flake-utils, ... }@inputs:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { 
          inherit system;
          config.allowUnfree = true;
          overlays = [
            inputs.templ.overlays.default
          ];
        };
      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [ 
            go
            gotools
            templ
            gopls
            air
            tailwindcss
          ];
        };
      }
    );
}
