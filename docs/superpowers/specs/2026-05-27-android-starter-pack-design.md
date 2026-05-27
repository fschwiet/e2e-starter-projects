# Android Starter Pack — Design

**Date:** 2026-05-27

## Overview

A minimal, well-configured Android app serving as a clean starting point for future development. Single Compose-based "Hello, World!" activity in Kotlin. No frameworks or dependency injection. Full local build and verification pipeline.

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
│   │   │   │   └── values/strings.xml
│   │   │   └── AndroidManifest.xml
│   │   ├── test/
│   │   │   └── java/com/example/newapplicationname/
│   │   │       └── ExampleUnitTest.kt
│   │   └── androidTest/
│   │       └── java/com/example/newapplicationname/
│   │           └── MainActivityTest.kt
│   └── build.gradle.kts
├── gradle/
│   ├── wrapper/
│   │   ├── gradle-wrapper.jar
│   │   └── gradle-wrapper.properties
│   └── libs.versions.toml
├── build.gradle.kts
├── settings.gradle.kts
├── gradlew
├── gradlew.bat
├── .editorconfig
├── .gitignore
└── README.md
```

---

## Build Configuration

**Build scripts:** Kotlin DSL (`.gradle.kts`) throughout.

**UI toolkit:** Jetpack Compose (current Android default). No XML layouts — the activity sets its content via `setContent { ... }`.

**SDK and JVM versions:**

- `compileSdk`: 35 (Android 15)
- `targetSdk`: 35
- `minSdk`: 24 (Android 7.0, ~97% of active devices as of now)
- `sourceCompatibility` / `targetCompatibility` / Kotlin `jvmTarget`: 17
- JDK required to build: 17 (matches AGP 8.7 requirement)

**Application identity:**

- `namespace` and `applicationId`: `com.example.newapplicationname`, declared in `app/build.gradle.kts` (`namespace` is required by AGP 8+)

**Release builds:** R8/minification is disabled in the starter (`isMinifyEnabled = false`). Forkers can enable it when they have something worth shrinking.

**Tool versions** (declared in `gradle/libs.versions.toml`):

- Kotlin: 2.0.x (current stable)
- Android Gradle Plugin: 8.7.x (current stable)
- Compose Compiler Gradle plugin: `org.jetbrains.kotlin.plugin.compose` (required with Kotlin 2.0+; replaces the old `composeOptions` block)
- Compose BOM: latest stable (drives all `androidx.compose.*` versions)
- ktlint Gradle plugin: `jlleitschuh/ktlint-gradle`
- detekt: `io.gitlab.arturbosch.detekt` with default ruleset
- Android Lint: built into AGP, no extra config needed

**Test dependencies** (also in the version catalog):

- `androidx.test.ext:junit`
- `androidx.compose.ui:ui-test-junit4` (versioned via the Compose BOM)
- `androidx.compose.ui:ui-test-manifest` (debug-only — enables `createAndroidComposeRule`)

All dependency versions live in `gradle/libs.versions.toml` — no hardcoded versions in build files.

---

## Application

`MainActivity.kt` — a `ComponentActivity` whose `onCreate` calls `setContent { ... }` containing a `MaterialTheme` wrapping a Compose `Text` that displays `stringResource(R.string.hello_world)`.

`strings.xml` — defines `hello_world` as `"Hello, World!"`. Both the Compose UI and the E2E test reference this same string resource.

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
@get:Rule
val composeTestRule = createAndroidComposeRule<MainActivity>()

@Test
fun helloWorldIsDisplayed() {
    val expected = composeTestRule.activity.getString(R.string.hello_world)
    composeTestRule.onNodeWithText(expected).assertIsDisplayed()
}
```

Uses `createAndroidComposeRule<MainActivity>()` to launch the activity and the Compose testing API to assert the Hello World text is on screen. The expected string is loaded from the same `hello_world` resource used by the UI, so the test stays correct if the resource value changes.

---

## Code Style

`.editorconfig` is present but **empty**. ktlint uses its built-in defaults; the file exists so forkers have an obvious place to add overrides without needing to know where ktlint looks for configuration.

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

### Chained commands for development

```
./gradlew precommit             # all checks, then apply formatting
./gradlew precommitConnected    # all checks, e2e tests, then apply formatting
```

These tasks **mutate source files** (`ktlintFormat` rewrites in place). For a read-only CI pipeline, invoke `./gradlew check` (which already includes `ktlintCheck`) instead.

### How the chained tasks are wired

`check` is the built-in Gradle lifecycle task. Its dependencies in this project come from:

- AGP contributes `lint` and `test` to `check`
- `ktlint-gradle` contributes `ktlintCheck` to `check`
- `detekt` contributes `detekt` to `check`

`precommit` depends on `check` and `ktlintFormat`, with `ktlintFormat` ordered after `check` via `mustRunAfter`. Formatting is only applied if all checks pass.

`precommitConnected` depends on `check`, `connectedAndroidTest`, and `ktlintFormat`, ordered: `check` → `connectedAndroidTest` → `ktlintFormat`. Formatting is the final step, applied only if everything passes.

---

## `.gitignore`

Standard Android entries:

```
.gradle/
build/
.idea/
local.properties
*.iml
.DS_Store
captures/
```

---

## README Sections

**Prerequisites** — Android Studio (with SDK), JDK 17, connected device or emulator (for `precommitConnected`).

**Install** — clone the repo; run `./gradlew build` to verify setup.

**Development** — launch the app on a connected device:

```
./gradlew installDebug
```

**Verification Pipeline** — lists all individual and chained commands as shown above.
