# Build AnyShake Observer

## Requirements

- Go matching the version in [go.mod](../go.mod).
- Node.js and npm.
- Git and GNU Make.
- A C compiler and CGO enabled for the race detector used by `make test`.

## Build from source

1. Clone the repository and enter its root directory:

   ```sh
   git clone https://github.com/anyshake/observer.git
   cd observer
   ```

2. Install frontend dependencies and build the web interface:

   ```sh
   cd web/src
   npm ci
   npm run build
   cd ../..
   ```

   This creates `web/dist`, which Go embeds in the application. Build the frontend before running Go tests or building the backend.

3. Run the Go unit tests (optional):

   ```sh
   make test
   ```

   This runs `go test -race ./... -count=1`, matching CI and disabling cached test results. Run tests on the host platform before setting cross-compilation variables.

4. Build the backend:

   ```sh
   make build
   ```

   The binary is written to `build/dist/observer`, and configuration assets are copied to `build/dist/assets`. The build disables CGO. Use `GO=/path/to/go` with Make to select a different Go executable.

## Cross-compilation

Set `GOOS` and `GOARCH` to select a target platform. For example:

```sh
GOOS=linux GOARCH=arm64 make build
GOOS=windows GOARCH=amd64 make build
```

Windows builds produce `build/dist/observer.exe`. Set `GOARM` or `GOMIPS` when required by the target architecture.

## Run the application

Edit `build/dist/assets/config.json` to configure the hardware connection, database, and server, then run:

```sh
./build/dist/observer --config ./build/dist/assets/config.json
```
