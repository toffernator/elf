{ pkgs, lib, config, inputs, ... }:

let appName = "main";
in {
  env = {
    GOOSE_DRIVER = "sqlite3";
    GOOSE_DBSTRING = "elf.db";
    GOOSE_MIGRATION_DIR = "db/migrations";
  };

  # FIXME: Install `templ`.
  packages = with pkgs; [ git tailwindcss air goose ];

  languages.go.enable = true;

  scripts = {
    db_create-migration = {
      exec = ''
        ${pkgs.goose}/bin/goose create $1 sql
      '';
    };

    db_up = {
      exec = ''
        ${pkgs.goose}/bin/goose up
      '';
    };

    db_up-one = {
      exec = ''
        ${pkgs.goose}/bin/goose up-by-one
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

    build = {
      exec = ''
        tailwind-build
        templ-generate
        ${pkgs.go}/bin/go build -ldflags "-X main.Environment=production" -o ./bin/${appName} ./cmd/${appName}/main.go
      '';
    };

    test = {
      exec = ''
        ${pkgs.go}/bin/go test -race -v -timeout 30s 
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

  };

  processes = {
    tailwind-watch = {
      exec = ''
        ${pkgs.tailwindcss}/bin/tailwindcss -i views/css/app.css -o public/styles.css --watch
      '';
    };

    templ-watch = {
      exec = ''
        nix run github:a-h/templ -- generate --watch --proxy http://localhost:4000
      '';
    };

    air = {
      exec = ''
        ${pkgs.air}/bin/air
      '';
    };
  };

}
