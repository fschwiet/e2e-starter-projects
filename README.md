# new-application-name

A minimal Android starter pack: single Compose "Hello, World!" activity in Kotlin, with a full local verification pipeline.

Rename the project: search-and-replace `new-application-name` (folder/project) and `com.example.newapplicationname` (package/applicationId) across the repo, then move the source under `app/src/*/java/` to match the new package.

## Prerequisites

- **Android Studio** with the Android SDK (compileSdk 35 / Android 15 installed).
  - Run .\scripts\check-prerequisites.ps1 to verify environment variables have been set
    - (JAVA_HOME, ANDROID_HOME, related path variables)
- **JDK 17 or later** as the active JDK (`java -version`). Android Gradle Plugin 8.7 requires 17+; JDK 21 works.
- **Connected device or emulator** with API level 24+, screen unlocked — only required for `precommitConnected` and `connectedAndroidTest`.

## Getting Started

```bash
./gradlew build
```

The first build downloads dependencies and may take several minutes.

## Run the project on a connected device

Install and launch the debug APK on a connected device:

```bash
./gradlew installDebug
```

## Verification Pipeline

The full verification pipeline — the individual gradle commands, the chained `precommit` /
`precommitConnected` / `ci` aggregates, and guidance on which to run when — lives in
[AGENTS.md](AGENTS.md#verification-pipeline), the canonical reference for the development loop.
