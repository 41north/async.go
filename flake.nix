{
  description = "A collection of utilities for async code in Go.";

  nixConfig = {
    substituters = [
      "https://cache.nixos.org"
      "https://nix-community.cachix.org"
    ];
    trusted-public-keys = [
      "cache.nixos.org-1:6NCHdD59X431o0gWypbMrAURkbJ16ZPMQFGspcDShjY="
      "nix-community.cachix.org-1:mB9FSh9qf2dCimDSUo8Zy7bkq5CX+/rkCWyvRCYg3Fs="
    ];
  };

  inputs = {
    devshell = {
      url = "github:numtide/devshell";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    pre-commit-hooks = {
      url = "github:cachix/pre-commit-hooks.nix";
      inputs.flake-utils.follows = "flake-utils";
    };
  };

  outputs = {
    self,
    devshell,
    nixpkgs,
    flake-utils,
    pre-commit-hooks,
    ...
  } @ inputs: let
    inherit (flake-utils.lib) eachSystem flattenTree mkApp;
  in
    eachSystem
    [
      "aarch64-linux"
      "aarch64-darwin"
      "x86_64-darwin"
      "x86_64-linux"
    ]
    (system: let
      pkgs = import nixpkgs {
        inherit system;
        overlays = [
          devshell.overlay
        ];
      };

      inherit (pkgs) dockerTools buildGoModule;
      inherit (pkgs.stdenv) isLinux;
      inherit (pkgs.lib) lists fakeSha256 licenses platforms;

      pkgWithCategory = category: package: {inherit package category;};

      # devshell command categories
      dev = pkgWithCategory "dev";
      linter = pkgWithCategory "linters";
      formatter = pkgWithCategory "formatters";
      util = pkgWithCategory "utils";
    in {
      devShell = pkgs.devshell.mkShell {
        env = [
          # disable CGO for now
          {
            name = "CGO_ENABLED";
            value = "0";
          }
        ];
        packages = with pkgs; [
          alejandra # https://github.com/kamadorueda/alejandra
          delve # https://github.com/go-delve/delve
          go_1_19 # https://go.dev/
          gofumpt # https://github.com/mvdan/gofumpt
          gotools # https://go.googlesource.com/tools
          nodePackages.prettier # https://prettier.io/
          treefmt # https://github.com/numtide/treefmt
          websocat # https://github.com/vi/websocat
        ];

        commands = with pkgs; [
          (formatter alejandra)
          (formatter gofumpt)
          (formatter nodePackages.prettier)

          (linter golangci-lint)

          (util jq)
          (util just)
        ];
      };
    });
}
