# Stops the running Android emulator.
# Exit codes: 0 ok (including no-op when nothing is running) | 2 adb not found
function Resolve-Adb {
    if ($env:ANDROID_HOME) {
        $candidate = Join-Path $env:ANDROID_HOME 'platform-tools\adb.exe'
        if (Test-Path $candidate) { return $candidate }
    }
    $onPath = Get-Command adb -ErrorAction SilentlyContinue
    if ($onPath) { return $onPath.Source }
    return $null
}

$adb = Resolve-Adb
if (-not $adb) {
    Write-Error "adb not found. Set ANDROID_HOME or add platform-tools to PATH."
    exit 2
}

$serial = $null
foreach ($line in (& $adb devices 2>&1)) {
    if ($line -match '^(emulator-\S+)\s+device$') { $serial = $Matches[1]; break }
}

if (-not $serial) {
    Write-Host "No emulator running. Nothing to stop."
    exit 0
}

Write-Host "Stopping emulator: $serial"
& $adb -s $serial emu kill
exit 0
