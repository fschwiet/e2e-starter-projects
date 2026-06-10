# Starts an Android emulator (clean boot) and waits until it is fully booted.
# Leaves the emulator running. Idempotent: no-op if a device is already connected.
# Usage: ./scripts/start-emulator.ps1 [-Avd <name>]
# Exit codes: 0 ready | 2 missing tool / no AVD | 3 boot timeout
param(
    [string]$Avd
)

function Resolve-Adb {
    if ($env:ANDROID_HOME) {
        $candidate = Join-Path $env:ANDROID_HOME 'platform-tools\adb.exe'
        if (Test-Path $candidate) { return $candidate }
    }
    $onPath = Get-Command adb -ErrorAction SilentlyContinue
    if ($onPath) { return $onPath.Source }
    return $null
}

function Resolve-Emulator {
    if ($env:ANDROID_HOME) {
        $candidate = Join-Path $env:ANDROID_HOME 'emulator\emulator.exe'
        if (Test-Path $candidate) { return $candidate }
    }
    return $null
}

$adb = Resolve-Adb
if (-not $adb) {
    Write-Error "adb not found. Set ANDROID_HOME or add platform-tools to PATH. Run scripts/check-prerequisites.ps1 first."
    exit 2
}

$emu = Resolve-Emulator
if (-not $emu) {
    Write-Error "emulator not found. Set ANDROID_HOME (expected `$env:ANDROID_HOME\emulator\emulator.exe)."
    exit 2
}

# Idempotency: if a device is already connected, do nothing.
$existing = @()
foreach ($line in (& $adb devices 2>&1)) {
    if ($line -match '^(\S+)\s+device$') { $existing += $Matches[1] }
}
if ($existing.Count -gt 0) {
    Write-Host "Device already connected: $($existing[0]). Nothing to start."
    exit 0
}

# Pick an AVD.
$avds = & $emu -list-avds 2>$null
if (-not $avds) {
    Write-Error "No AVDs found. Create one in Android Studio's Device Manager."
    exit 2
}
if (-not $Avd) {
    $Avd = ($avds | Select-Object -First 1)
} elseif ($avds -notcontains $Avd) {
    Write-Error "AVD '$Avd' not found. Available: $($avds -join ', ')"
    exit 2
}

# Launch detached so the emulator keeps running after this script exits.
# -no-snapshot = a true clean boot every time (no snapshot load OR save). Using
# -no-snapshot-save instead would *load* the quick-boot snapshot, which hangs if that
# snapshot was corrupted by a previous force-close.
Write-Host "Starting emulator '$Avd' (clean boot)..."
Start-Process -FilePath $emu -ArgumentList @('-avd', $Avd, '-no-snapshot')

# Wait for the device to register with adb, then for boot to complete.
& $adb wait-for-device
$deadline = (Get-Date).AddSeconds(180)
while ((Get-Date) -lt $deadline) {
    $booted = (& $adb shell getprop sys.boot_completed 2>$null | Out-String).Trim()
    if ($booted -eq '1') {
        $serialLine = (& $adb devices) | Select-String -Pattern '^(\S+)\s+device$' | Select-Object -First 1
        $serial = if ($serialLine) { $serialLine.Matches[0].Groups[1].Value } else { 'unknown' }
        Write-Host "Emulator ready: $serial"
        exit 0
    }
    Start-Sleep -Seconds 3
}

Write-Error "Emulator did not finish booting within 180s."
exit 3
