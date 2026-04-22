[CmdletBinding()]
param(
    [string]$Repo = "bodrovis/lokex-cli",
    [string]$BinName = "lokex-cli",
    [string]$Version = "latest",
    [string]$InstallDir = "$env:LOCALAPPDATA\Programs\lokex-cli\bin"
)

$ErrorActionPreference = "Stop"
Set-StrictMode -Version 2.0

try {
    [Net.ServicePointManager]::SecurityProtocol = `
        [Net.SecurityProtocolType]::Tls12 -bor `
        [Net.SecurityProtocolType]::Tls13
} catch {
    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
}

function Get-ArchName {
    $arch = $env:PROCESSOR_ARCHITEW6432
    if (-not $arch) {
        $arch = $env:PROCESSOR_ARCHITECTURE
    }

    switch ($arch.ToUpperInvariant()) {
        "AMD64" { return "x86_64" }
        "X86"   { return "i386" }
        "ARM64" { return "arm64" }
        default { throw "unsupported architecture: $arch" }
    }
}

function Normalize-Version([string]$InputVersion) {
    if ($InputVersion -eq "latest") {
        return "latest"
    }

    if ($InputVersion.StartsWith("v")) {
        return $InputVersion
    }

    return "v$InputVersion"
}

function Invoke-Download {
    param(
        [Parameter(Mandatory = $true)][string]$Url,
        [Parameter(Mandatory = $true)][string]$OutFile
    )

    Invoke-WebRequest -UseBasicParsing -Uri $Url -Headers @{ "User-Agent" = "lokex-cli-installer" } -OutFile $OutFile
}

function Get-ReleaseJson {
    param(
        [Parameter(Mandatory = $true)][string]$Url
    )

    $response = Invoke-WebRequest -UseBasicParsing -Uri $Url -Headers @{ "User-Agent" = "lokex-cli-installer" }
    return ($response.Content | ConvertFrom-Json)
}

function Get-ExpectedChecksum {
    param(
        [Parameter(Mandatory = $true)][string]$ChecksumsPath,
        [Parameter(Mandatory = $true)][string]$AssetName
    )

    foreach ($line in Get-Content -LiteralPath $ChecksumsPath) {
        if ($line -match '^(?<hash>[A-Fa-f0-9]{64})\s+(?<file>.+)$') {
            if ($Matches["file"] -eq $AssetName) {
                return $Matches["hash"].ToLowerInvariant()
            }
        }
    }

    throw "checksum for $AssetName not found in checksums.txt"
}

function Add-ToUserPath {
    param(
        [Parameter(Mandatory = $true)][string]$Dir
    )

    $dirNorm = $Dir.TrimEnd('\')

    $machinePath = [Environment]::GetEnvironmentVariable("Path", "Machine")
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")

    $allEntries = @()
    if ($machinePath) { $allEntries += ($machinePath -split ';') }
    if ($userPath)    { $allEntries += ($userPath -split ';') }

    $alreadyPresent = $allEntries |
        Where-Object { $_ -and $_.Trim() } |
        ForEach-Object { $_.Trim().TrimEnd('\') } |
        Where-Object { $_ -ieq $dirNorm } |
        Select-Object -First 1

    if ($alreadyPresent) {
        return $false
    }

    $newUserPath = if ($userPath -and $userPath.Trim()) {
        "$userPath;$Dir"
    } else {
        $Dir
    }

    [Environment]::SetEnvironmentVariable("Path", $newUserPath, "User")
    $env:Path = "$env:Path;$Dir"

    return $true
}

$goreArch = Get-ArchName
$normalizedVersion = Normalize-Version $Version

$tmpDir = Join-Path ([System.IO.Path]::GetTempPath()) ([System.IO.Path]::GetRandomFileName())
$null = New-Item -ItemType Directory -Path $tmpDir

try {
    if ($normalizedVersion -eq "latest") {
        $apiUrl = "https://api.github.com/repos/$Repo/releases/latest"
        $release = Get-ReleaseJson -Url $apiUrl

        $asset = $release.assets |
            Where-Object { $_.name -match "^$([regex]::Escape($BinName))_.+_Windows_$([regex]::Escape($goreArch))\.zip$" } |
            Select-Object -First 1

        $checksums = $release.assets |
            Where-Object { $_.name -eq "checksums.txt" } |
            Select-Object -First 1

        if (-not $asset) {
            throw "could not find release asset for Windows/$goreArch"
        }

        if (-not $checksums) {
            throw "could not find checksums.txt for latest release"
        }

        $assetUrl = $asset.browser_download_url
        $checksumsUrl = $checksums.browser_download_url
    }
    else {
        $assetUrl = "https://github.com/$Repo/releases/download/$normalizedVersion/${BinName}_${normalizedVersion}_Windows_${goreArch}.zip"
        $checksumsUrl = "https://github.com/$Repo/releases/download/$normalizedVersion/checksums.txt"
    }

    $assetName = Split-Path -Leaf $assetUrl
    $archivePath = Join-Path $tmpDir $assetName
    $checksumsPath = Join-Path $tmpDir "checksums.txt"
    $extractDir = Join-Path $tmpDir "extract"

    Write-Host "Downloading $assetName..."
    Invoke-Download -Url $assetUrl -OutFile $archivePath
    Invoke-Download -Url $checksumsUrl -OutFile $checksumsPath

    $expectedChecksum = Get-ExpectedChecksum -ChecksumsPath $checksumsPath -AssetName $assetName
    $actualChecksum = (Get-FileHash -Algorithm SHA256 -LiteralPath $archivePath).Hash.ToLowerInvariant()

    if ($actualChecksum -ne $expectedChecksum) {
        throw @"
checksum mismatch for $assetName
expected: $expectedChecksum
actual:   $actualChecksum
"@
    }

    Write-Host "Checksum verified."

    Expand-Archive -LiteralPath $archivePath -DestinationPath $extractDir -Force

    $exeName = "$BinName.exe"
    $binary = Get-ChildItem -LiteralPath $extractDir -Recurse -File |
        Where-Object { $_.Name -ieq $exeName } |
        Select-Object -First 1

    if (-not $binary) {
        throw "binary $exeName not found in archive"
    }

    if (-not (Test-Path -LiteralPath $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }

    $targetPath = Join-Path $InstallDir $exeName
    Copy-Item -LiteralPath $binary.FullName -Destination $targetPath -Force

    $pathAdded = Add-ToUserPath -Dir $InstallDir

    Write-Host ""
    Write-Host "Installed to $targetPath"

    if ($pathAdded) {
        Write-Host "Added $InstallDir to your user PATH."
        Write-Host "Open a new terminal to use $exeName everywhere."
    }
    else {
        Write-Host "$InstallDir is already in PATH."
    }

    Write-Host ""
    Write-Host "Ready to use: $exeName"
}
finally {
    if (Test-Path -LiteralPath $tmpDir) {
        Remove-Item -LiteralPath $tmpDir -Recurse -Force
    }
}