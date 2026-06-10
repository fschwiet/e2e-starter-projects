## Dev helper scripts (`scripts/`)

These PowerShell scripts wrap the environment checks and emulator lifecycle that recur in the development loop. Use them to avoid permission denials on the more general underlying commands.

| Script                    | Purpose                                                                                                | Exit codes                                              |
| ------------------------- | ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------- |
| `check-prerequisites.ps1` | Static toolchain gate: `JAVA_HOME`, `ANDROID_HOME`, `java >= 17`, `adb` on PATH.                       | `0` ok · `2` something missing                          |
| `check-device.ps1`        | Runtime readiness: is an emulator/device connected, awake, and unlocked? Lists AVDs if none connected. | `0` ready · `2` no device · `3` asleep/locked           |
| `start-emulator.ps1`      | Start an emulator (clean boot), wait until booted, leave it running. Idempotent. `-Avd <name>` optional.| `0` ready · `2` missing tool / no AVD · `3` boot timeout |
| `stop-emulator.ps1`       | Kill the running emulator (no-op if none running).                                                     | `0` ok · `2` adb not found                              |

## Verify a change

Run these in order before considering a change done:

1. `./scripts/check-prerequisites.ps1` — toolchain + env vars (static).
2. `./scripts/check-device.ps1` — is a device connected & ready? (runtime). If it reports
   **no device** and you need connected/e2e tests, start one with
   `./scripts/start-emulator.ps1`, then re-run `check-device.ps1`.
3. Run the appropriate verification task (see below).

Stop the emulator when finished with `./scripts/stop-emulator.ps1`.

## Verification pipeline

Standard Gradle wrapper tasks (already allowlisted as `./gradlew *`).

Individual commands:

- `./gradlew detekt` — static analysis.
- `./gradlew lint` — Android lint.
- `./gradlew test` — JVM unit tests.
- `./gradlew connectedAndroidTest` — e2e tests (requires a ready device/emulator).
- `./gradlew ktlintCheck` — check formatting.
- `./gradlew ktlintFormat` — apply formatting.
- `./gradlew :app:assembleDebug` — build the debug APK.
- `./gradlew installDebug` — install and launch the debug APK on a connected device.

Chained (which to run when):

- `./gradlew precommit` — all checks, then apply formatting. Use for changes that don't need a device.
- `./gradlew precommitConnected` — all checks + e2e tests, then apply formatting. Use when a device is ready and the change touches device behavior.
- `./gradlew ci` — all checks + e2e tests, no formatting. The read-only equivalent used by GitHub Actions.

`precommit` / `precommitConnected` **mutate source files** (they run `ktlintFormat` at the end). `ci` is the read-only equivalent — it runs `check` (`ktlintCheck`, `detekt`, `lint`, `test`) plus `connectedAndroidTest`.

## Starting / stopping an emulator manually

`start-emulator.ps1` / `stop-emulator.ps1` wrap this. To do it by hand:

```powershell
& "$env:ANDROID_HOME\emulator\emulator.exe" -avd Pixel_6 -no-snapshot   # clean boot
# wait for boot:
& "$env:ANDROID_HOME\platform-tools\adb.exe" wait-for-device
# stop it:
& "$env:ANDROID_HOME\platform-tools\adb.exe" -s emulator-5554 emu kill
```
