# Android Starter Pack — Design

**Date:** 2026-05-27

## Overview

A minimal, well-configured Android app serving as a clean starting point for future development. Single "Hello, World!" activity in Kotlin. No frameworks or dependency injection. Full local build and verification pipeline.

Use `new-application-name` as the placeholder for the app/project name throughout — easy to find-and-replace when forking.

---

## Project Structure

```
new-application-name/
├── app/
│   ├── src/
│   │   ├── main/
│   │   │   ├── java/com/example/newapplicationname/
│   │   │   │   └── MainActivity.kt
│   │   │   ├── res/
│   │   │   │   ├── layout/activity_main.xml
│   │   │   │   └── values/strings.xml
│   │   │   └── AndroidManifest.xml
│   │   ├── test/
│   │   │   └── java/com/example/newapplicationname/
│   │   │       └── ExampleUnitTest.kt
│   │   └── androidTest/
│   │       └── java/com/example/newapplicationname/
│   │           └── MainActivityTest.kt
│   └── build.gradle.kts
├── build.gradle.kts
├── settings.gradle.kts
├── gradle/libs.versions.toml
└── README.md
```

---

## Build Configuration

**Build scripts:** Kotlin DSL (`.gradle.kts`) throughout.

**SDK versions:**
- `compileSdk`: 35 (Android 15)
- `targetSdk`: 35
- `minSdk`: 24 (Android 7.0, ~97% of active devices)

**Tool versions** (declared in `gradle/libs.versions.toml`):
- Kotlin: 2.0.x (current stable)
- Android Gradle Plugin: 8.7.x (current stable)
- ktlint Gradle plugin: `jlleitschuh/ktlint-gradle`
- detekt: `io.gitlab.arturbosch.detekt` with default ruleset
- Android Lint: built into AGP, no extra config needed

All dependency versions live in `gradle/libs.versions.toml` — no hardcoded versions in build files.

---

## Application

`MainActivity.kt` — a single Activity that sets its content view to `activity_main.xml`.

`activity_main.xml` — a simple layout containing a `TextView` displaying the string resource `@string/hello_world`.

`strings.xml` — defines `hello_world` as `"Hello, World!"`. Both the app layout and the E2E test reference this same string value.

---

## Tests

### Unit Test — `ExampleUnitTest.kt`

```kotlin
@Test
fun addition_isCorrect() {
    assertEquals(4, 2 + 2)
}
```

Verifies the unit test infrastructure is wired up correctly.

### E2E Test — `MainActivityTest.kt`

```kotlin
@Test
fun helloWorldIsDisplayed() {
    onView(withText("Hello, World!"))
        .check(matches(isDisplayed()))
}
```

Uses `ActivityScenarioRule` to launch `MainActivity` and Espresso to assert the Hello World text is visible on screen.

---

## Verification Pipeline

### Individual commands

```
./gradlew detekt                # static analysis
./gradlew lint                  # Android lint
./gradlew test                  # unit tests
./gradlew connectedAndroidTest  # e2e tests (requires connected device/emulator)
./gradlew ktlintCheck           # check formatting
./gradlew ktlintFormat          # apply formatting
```

### Composite commands

```
./gradlew ci                    # all checks, then apply formatting
./gradlew ciConnected           # all checks, e2e tests, then apply formatting
```

### How the composite tasks are wired

`check` (built-in Gradle lifecycle task, extended by plugins) runs: `ktlintCheck`, `detekt`, `lint`, `test`.

`ci` depends on `check` and `ktlintFormat`, with `ktlintFormat` ordered after `check`. Formatting is only applied if all checks pass.

`ciConnected` depends on `check`, `connectedAndroidTest`, and `ktlintFormat`, ordered: `check` → `connectedAndroidTest` → `ktlintFormat`. Formatting is the final step, applied only if everything passes.

---

## README Sections

**Prerequisites** — Android Studio (with SDK), Java JDK, connected device or emulator (for `ciConnected`).

**Install** — clone the repo; run `./gradlew build` to verify setup.

**Development** — launch the app on a connected device:
```
./gradlew installDebug
```

**Verification Pipeline** — lists all individual and composite commands as shown above.
