# Elf

## Directory Structure

`cmd` contains all the executables:
- `server` a REST server that performs server-side rendering using `htmx` and `templ` for templating.

`db` contains database related files:
- `migrations` managed using `goose`.
- `seeds` managed using `goose`.

`internal` contains the implementation of `elf`:
- `core` contains models the data along with validation rules of these. These files should change the least.
- `sqlite` contains solutions for persisting data.
- `service` describes the actions on data models that are defined in `core`.
- `rest` the transport layer implementation that `cmd/server` relies on.

## Development

Elf relies on `devenv` for setting up and configuring a local development environment.
`devenv.nix` is the authorative source on dependencies and their versions.

```sh
devenv shell
```

