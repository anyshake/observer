# ./Makefile.ps1

param(
    [string]$target
)

function Build {
    Clean
    Pre-Build
    npm run build -- --GENERATE_SOURCEMAP=false
    Remove-Item -Recurse -Force "../dist"
    Move-Item "./build" "../dist"
}

function Run {
    Pre-Build
    $env:BROWSER = "none"
    npm run start
}

function Clean {
    if (Test-Path "./build") {
        Remove-Item -Recurse -Force "./build"
    }
    if (Test-Path "../dist") {
        Remove-Item -Recurse -Force "../dist"
    }
    New-Item -ItemType Directory -Force -Path "../dist"
}

function Pre-Build {
    $version = Get-Content "../../VERSION"
    $release = Get-Date -Format "yyyyMMddHHmmss"
    $commit = git rev-parse --short HEAD
    "$env:REACT_APP_VERSION=$version" > .env
    "$env:REACT_APP_RELEASE=$commit-$release" >> .env
}

switch ($target) {
    "build" { Build }
    "run" { Run }
    "clean" { Clean }
    "pre-build" { Pre-Build }
    default { Write-Host "No valid target specified. Use 'build', 'run', 'clean', or 'pre-build'." }
}
