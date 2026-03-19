package stages

import (
	"fmt"
	"time"

	"github.com/tensorhero/leetgpu-tester/internal/helpers"
	"github.com/tensorhero/tester-utils/runner"
	"github.com/tensorhero/tester-utils/test_case_harness"
	"github.com/tensorhero/tester-utils/tester_definition"
)

func e03GuardTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "guard",
		Timeout:     30 * time.Second,
		TestFunc:    testE03Guard,
		CompileStep: compileStep("test_e03"),
	}
}

func testE03Guard(harness *test_case_harness.TestCaseHarness) error {
	logger := harness.Logger
	workDir := harness.SubmissionDir
	lang := harness.DetectedLang

	r := runner.Run(workDir, lang.RunCmd, lang.RunArgs...).
		WithTimeout(10 * time.Second).
		WithLogger(logger).
		Execute().
		Exit(0)

	if err := r.Error(); err != nil {
		return fmt.Errorf("test driver failed: %v", err)
	}

	results := helpers.ParseStructuredOutput(string(r.Result().Stdout))

	tests := []struct {
		name     string
		expected string
		label    string
	}{
		{"basic_match", "True", "guard_kernel([0,1,2,3], size=4, threads=8) matches expected"},
		{"no_overflow", "True", "out[size:] remains zero (no out-of-bounds writes)"},
		{"guard_works", "True", "guard_kernel works with size=6, threads=8"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All E03 tests passed!")
	return nil
}
