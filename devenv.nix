{
  pkgs,
  inputs,
  ...
}:
let
  pkgs-unstable = import inputs.nixpkgs-unstable { system = pkgs.stdenv.system; };
in
{
  packages =
    with pkgs;
    [
      git
      go
      golangci-lint-langserver
      gopls
      clang
      entr
      prettier
      vscode-css-languageserver
    ]
    ++ [
      # golangci-lint moves fast, and I like to have its latest rules.
      pkgs-unstable.golangci-lint
    ];

  scripts = {
    run.exec = ''
      find . | entr -c -r bash -c 'CGO_ENABLE=1 go run cmd/deck/main.go'
    '';
    lint.exec = ''
      golangci-lint run
      clang-format -i ./internal/media/bridge_darwin.*
    '';
  };
}
