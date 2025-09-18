$outputDirectory = "C:\src\up_bins\"

# The output directory exists
if (-Not (Test-Path -Path $outputDirectory)) {
    New-Item -ItemType Directory -Path $outputDirectory
}

$env:GOOS = "windows"
$env:GOARCH = "amd64"

# Build the Go program
go build -ldflags="-s -w" -trimpath -o "$outputDirectory\crun.exe" .
if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed with exit code $LASTEXITCODE"
    exit $LASTEXITCODE
} else {
    Write-Host "Build succeeded. Output: $outputDirectory\crun.exe"
}
exit 0