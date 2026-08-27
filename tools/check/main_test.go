package main

import (
	"strings"
	"testing"
)

func stepNamesOf(steps []step) []string {
	names := make([]string, 0, len(steps))
	for _, s := range steps {
		names = append(names, s.name)
	}

	return names
}

func TestSelectStepsWithNoArgsRunsWholePipeline(t *testing.T) {
	steps, err := selectSteps(nil)
	if err != nil {
		t.Fatalf("selectSteps(nil) returned error: %v", err)
	}

	got := strings.Join(stepNamesOf(steps), ",")
	want := "format,lint,unit,build,e2e"

	if got != want {
		t.Errorf("step names = %q, want %q", got, want)
	}
}

func TestSelectStepsWithOneNameRunsOnlyThatStep(t *testing.T) {
	steps, err := selectSteps([]string{"lint"})
	if err != nil {
		t.Fatalf("selectSteps([lint]) returned error: %v", err)
	}

	got := strings.Join(stepNamesOf(steps), ",")
	if want := "lint"; got != want {
		t.Errorf("step names = %q, want %q", got, want)
	}
}

func TestSelectStepsBuildsBeforeE2E(t *testing.T) {
	// Running e2e alone against a stale binary could report a false pass, so
	// the runner always rebuilds first.
	steps, err := selectSteps([]string{"e2e"})
	if err != nil {
		t.Fatalf("selectSteps([e2e]) returned error: %v", err)
	}

	got := strings.Join(stepNamesOf(steps), ",")
	if want := "build,e2e"; got != want {
		t.Errorf("step names = %q, want %q", got, want)
	}
}

func TestSelectStepsRejectsUnknownStep(t *testing.T) {
	_, err := selectSteps([]string{"nope"})
	if err == nil {
		t.Fatal("selectSteps([nope]) returned nil error, want an error")
	}

	// The message must list the valid names so the user can recover.
	for _, name := range []string{"format", "lint", "unit", "build", "e2e"} {
		if !strings.Contains(err.Error(), name) {
			t.Errorf("error %q does not mention valid step %q", err, name)
		}
	}
}

func TestSelectStepsRejectsMultipleArgs(t *testing.T) {
	if _, err := selectSteps([]string{"lint", "unit"}); err == nil {
		t.Fatal("selectSteps([lint unit]) returned nil error, want an error")
	}
}
