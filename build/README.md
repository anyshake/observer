# Build AnyShake Observer

In this document manual procedures are explained (without `make` and continuous integration).

## Prerequisites

 - [Go 1.19 or later](https://golang.org/doc/install)
 - [NPM 9.3 or later](https://www.npmjs.com/get-npm)
 - [Node.js 18.14 or later](https://nodejs.org/en/download/)
 - [Git 2.20 or later](https://git-scm.com/downloads) (Optional)

**For Windows developers, make sure you're using Windows 10 or later. And if you have [Git for Windows](https://git-scm.com/download/win) installed, you can use Git Bash to run the commands in this document. Otherwise, please use PowerShell instead of CMD.**

## Get Source Code

Clone the repository with `--depth=1` to avoid downloading the entire history.

``` bash
$ git clone --depth=1 https://github.com/anyshake/observer
```

If you don't have `git` installed, you can download the source code as a `.zip` or `.tar.gz` archive, extract it and `cd` into the directory. However, you will not be able to update the source code with `git pull`.

 - [github.com/anyshake/observer/archive/refs/heads/master.zip](https://github.com/anyshake/observer/archive/refs/heads/master.zip)
 - [github.com/anyshake/observer/archive/refs/heads/master.tar.gz](https://github.com/anyshake/observer/archive/refs/heads/master.tar.gz)

## Build Frontend

Remove the `frontend/dist` to avoid conflicts.

``` bash
$ rm -rf frontend/dist
```

Install dependencies and build the frontend.

``` bash
$ cd frontend/src
$ npm install
$ npm run build
```

Move the output `build` directory to parent, rename it to `dist`.

``` bash
$ mv build ../dist
```

## Build Backend

Go back to the repository root directory to build backend binary.

``` bash
$ CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o observer cmd/*.go
```

If you're using Windows with PowerShell, please use the following command instead.

```powershell
$ $env:CGO_ENABLED=0
$ go build -ldflags="-s -w" -trimpath -o observer.exe cmd/*.go
```

The output file will be `observer` (Windows is `observer.exe`), test it with `-version` or `-help` flag.
