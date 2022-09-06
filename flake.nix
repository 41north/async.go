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

      linters = with pkgs; [
        alejandra # https://github.com/kamadorueda/alejandra
        gofumpt # https://github.com/mvdan/gofumpt
        nodePackages.prettier # https://prettier.io/
        treefmt # https://github.com/numtide/treefmt
      ];

      # devshell command categories
      dev = pkgWithCategory "dev";
      linter = pkgWithCategory "linters";
      formatter = pkgWithCategory "formatters";
      util = pkgWithCategory "utils";
    in {
      checks = {
        format =
          pkgs.runCommandNoCC "treefmt" {
            nativeBuildInputs = linters;
          } ''
            # keep timestamps so that treefmt is able to detect mtime changes
            cp --no-preserve=mode --preserve=timestamps -r ${self} source
            cd source
            HOME=$TMPDIR treefmt --fail-on-change
            touch $out
          '';
      };

      devShell = pkgs.devshell.mkShell {
        packages = with pkgs;
          [
            delve # https://github.com/go-delve/delve
            go_1_19 # https://go.dev/
            gotools # https://go.googlesource.com/tools
            websocat # https://github.com/vi/websocat
          ]
          ++ linters;

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
