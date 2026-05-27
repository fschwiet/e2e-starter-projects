$failed = $false

if (-not $env:JAVA_HOME) {
    Write-Error "JAVA_HOME is not set. Set it with:`n[System.Environment]::SetEnvironmentVariable('JAVA_HOME', 'C:\Program Files\Android\Android Studio\jbr', 'User')"
    $failed = $true
}

if (-not $env:ANDROID_HOME) {
    Write-Error "ANDROID_HOME is not set. Set it with:`n[System.Environment]::SetEnvironmentVariable('ANDROID_HOME', `"`$env:LOCALAPPDATA\Android\Sdk`", 'User')"
    $failed = $true
}

if (-not (Get-Command java -ErrorAction SilentlyContinue)) {
    Write-Error "java is not on PATH. Add %JAVA_HOME%\bin to your PATH environment variable."
    $failed = $true
} else {
    $versionLine = & java -version 2>&1 | Select-String -Pattern '"(\d+)(?:\.(\d+))?'
    if ($versionLine) {
        $major = [int]$versionLine.Matches[0].Groups[1].Value
        if ($major -eq 1) {
            $major = [int]$versionLine.Matches[0].Groups[2].Value
        }
        if ($major -lt 17) {
            Write-Error "Java $major is below the required minimum of 17. Add %JAVA_HOME%\bin to your PATH to use the correct version."
            $failed = $true
        }
    } else {
        Write-Error "Could not determine Java version. Add %JAVA_HOME%\bin to your PATH environment variable."
        $failed = $true
    }
}

if (-not (Get-Command adb -ErrorAction SilentlyContinue)) {
    Write-Error "adb is not on PATH. Add %ANDROID_HOME%\platform-tools to your PATH environment variable."
    $failed = $true
}

if ($failed) {
    exit 2
}

Write-Host "All prerequisites satisfied."
exit 0
