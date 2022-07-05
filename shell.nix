#...
{ } :
  let
    # Pinned nixpkgs
    pkgs = import (builtins.fetchGit {
      name = "nixpkg-21.11";                         
      url = "https://github.com/NixOS/nixpkgs/";             
      ref = "refs/tags/21.05";           
      rev = "fefb0df7d2ab2e1cabde7312238026dcdc972441";                       
    }) {};                                       
  in
    pkgs.mkShell {
      nativeBuildInputs = with pkgs.buildPackages; [ 
        go
      ];
    }
