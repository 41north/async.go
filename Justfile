# just is a handy way to save and run project-specific commands.
#
# https://github.com/casey/just

# list all tasks
default:
  just --list

# Format the code
fmt:
  treefmt
alias f := fmt

# Start IDEA in this folder
idea:
  nohup idea-ultimate . > /dev/null 2>&1 &

# Start VsCode in this folder
code:
  code .

# Checks the source with nix
nix-check:
  nix flake check
