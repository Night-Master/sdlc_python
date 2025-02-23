# Target platforms
$platforms = @("windows/amd64", "linux/amd64")

# Loop through each platform
foreach ($platform in $platforms) {
    # Parse platform information
    $platform_split = $platform -split '/'
    $GOOS = $platform_split[0]
    $GOARCH = $platform_split[1]

    # Set output file name
    $output_name = "gosdlc_$GOOS-$GOARCH"
    if ($GOOS -eq "windows") {
        $output_name += ".exe"
    }

    # Set build environment variables
    $env:GOOS = $GOOS
    $env:GOARCH = $GOARCH

    # Enable CGO for Windows builds only
    if ($GOOS -eq "windows") {
        $env:CGO_ENABLED = 1
    } else {
        $env:CGO_ENABLED = 0
    }

    # Perform the build
    go build -o $output_name

    # Check for build errors
    if ($LASTEXITCODE -ne 0) {
        Write-Host "An error has occurred! Aborting the script execution..."
        exit 1
    }
}