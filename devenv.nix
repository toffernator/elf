{ pkgs, lib, config, inputs, ... }:

let appName = "main";
in {
  # FIXME: Install `templ`.
  packages = with pkgs; [ git tailwindcss air ];

  enterShell = ''
    exec zsh
  '';

  languages.go.enable = true;

  scripts = {
    tailwind-watch = {
      exec = ''
        ${pkgs.tailwindcss}/bin/tailwindcss -i views/css/app.css -o public/styles.css --watch
      '';
    };

    tailwind-build = {
      exec = ''
        ${pkgs.tailwindcss}/bin/tailwindcss -i views/css/app.css -o public/styles.css --minify
      '';
    };

    templ-generate = {
      exec = ''
        nix run github:a-h/templ -- generate
      '';
    };
    templ-watch = {
      exec = ''
        nix run github:a-h/templ -- generate --watch --proxy http://localhost:4000
      '';
    };
    dev = {
      exec = ''
        tailwind-watch &
        templ-watch &
        ${pkgs.air}/bin/air
      '';
    };
    build = {
      exec = ''
        tailwind-build
        templ-generate
        ${pkgs.go}/bin/go build -ldflags "-X main.Environment=production" -o ./bin/${appName} ./cmd/${appName}/main.go
      '';
    };
    vet = {
      exec = ''
        ${pkgs.go}/bin/go vet ./...
      '';
    };
    staticcheck = {
      exec = ''
        ${pkgs.gotools}/bin/staticcheck ./...
      '';
    };
    test = {
      exec = ''
        ${pkgs.go}/bin/go test -race -v -timeout 30s 
      '';
    };
  };
}
