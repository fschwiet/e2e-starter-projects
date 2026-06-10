# Checks that an Android device/emulator is connected and ready to install/test.
# Exit codes: 0 ready | 2 no device connected | 3 device asleep/locked

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

# Find a device in 'device' state (skip header, offline, unauthorized).
$lines = & $adb devices 2>&1
$ready = @()
foreach ($line in $lines) {
    if ($line -match '^(\S+)\s+device$') { $ready += $Matches[1] }
}

if ($ready.Count -eq 0) {
    Write-Host "No device connected (in 'device' state)."
    $emu = Resolve-Emulator
    if ($emu) {
        $avds = & $emu -list-avds 2>$null
        if ($avds) {
            Write-Host "Available AVDs:"
            $avds | ForEach-Object { Write-Host "  $_" }
            Write-Host "Start one with: ./scripts/start-emulator.ps1"
            Write-Host "  (or a specific AVD: ./scripts/start-emulator.ps1 -Avd $($avds | Select-Object -First 1))"
        } else {
            Write-Host "No AVDs found. Create one in Android Studio's Device Manager."
        }
    }
    exit 2
}

$serial = $ready[0]
Write-Host "Device connected: $serial"

# Wakefulness: Awake | Asleep | Dozing | Dreaming
$power = & $adb -s $serial shell dumpsys power 2>$null
$wakeMatch = $power | Select-String -Pattern 'mWakefulness=(\w+)' | Select-Object -First 1
$wakefulness = if ($wakeMatch) { $wakeMatch.Matches[0].Groups[1].Value } else { 'Unknown' }

# Keyguard (lock screen) state.
$kg = & $adb -s $serial shell dumpsys window 2>$null
$kgMatch = $kg | Select-String -Pattern 'mDreamingLockscreen=(\w+)|isKeyguardShowing.*?(true|false)|KeyguardShowing=(true|false)' | Select-Object -First 1
$locked = $false
if ($kgMatch) {
    $val = ($kgMatch.Matches[0].Groups | Where-Object { $_.Value -in @('true', 'false') } | Select-Object -First 1).Value
    $locked = ($val -eq 'true')
}

$lockText = if ($locked) { 'locked' } else { 'unlocked' }
Write-Host "State: $wakefulness, $lockText"

if ($wakefulness -ne 'Awake' -or $locked) {
    Write-Warning "Device is not test-ready (must be Awake and unlocked). Wake and unlock it with:"
    Write-Warning "  & `"$adb`" -s $serial shell input keyevent KEYCODE_WAKEUP"
    Write-Warning "  & `"$adb`" -s $serial shell input keyevent 82   # dismiss keyguard"
    exit 3
}

Write-Host "Device is ready."
exit 0
