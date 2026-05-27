# new-application-name

A minimal Android starter pack: single Compose "Hello, World!" activity in Kotlin, with a full local verification pipeline.

Rename the project: search-and-replace `new-application-name` (folder/project) and `com.example.newapplicationname` (package/applicationId) across the repo, then move the source under `app/src/*/java/` to match the new package.

## Prerequisites

- **Android Studio** with the Android SDK (compileSdk 35 / Android 15 installed).
- **JDK 17 or later** as the active JDK (`java -version`). Android Gradle Plugin 8.7 requires 17+; JDK 21 works.
- **Connected device or emulator** with API level 24+, screen unlocked — only required for `precommitConnected` and `connectedAndroidTest`.

## Install

```bash
git clone <repo-url>
cd new-application-name
./gradlew build
```

The first build downloads dependencies and may take several minutes.

## Development

Install and launch the debug APK on a connected device:

```bash
./gradlew installDebug
```

## Verification Pipeline

Individual commands:

```bash
./gradlew detekt                # static analysis
./gradlew lint                  # Android lint
./gradlew test                  # unit tests
./gradlew connectedAndroidTest  # e2e tests (requires connected device/emulator)
./gradlew ktlintCheck           # check formatting
./gradlew ktlintFormat          # apply formatting
```

Chained:

```bash
./gradlew precommit             # all checks, then apply formatting
./gradlew precommitConnected    # all checks, e2e tests, then apply formatting
./gradlew ci                    # all checks + e2e tests, no formatting (used by GitHub Actions)
```

`precommit` / `precommitConnected` **mutate source files** because they run `ktlintFormat` at the end. `ci` is the read-only equivalent used in CI — it runs `check` (which includes `ktlintCheck`, `detekt`, `lint`, and `test`) plus `connectedAndroidTest`.
