package stages

import (
	"fmt"
	"time"

	"github.com/bootcraft-cn/leetgpu-tester/internal/helpers"
	"github.com/bootcraft-cn/tester-utils/runner"
	"github.com/bootcraft-cn/tester-utils/test_case_harness"
	"github.com/bootcraft-cn/tester-utils/tester_definition"
)

func e01MapTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "map",
		Timeout:     30 * time.Second,
		TestFunc:    testE01Map,
		CompileStep: compileStep("test_e01"),
	}
}

func testE01Map(harness *test_case_harness.TestCaseHarness) error {
	logger := harness.Logger
	workDir := harness.SubmissionDir
	lang := harness.DetectedLang

	// Run test driver
	r := runner.Run(workDir, lang.RunCmd, lang.RunArgs...).
		WithTimeout(10 * time.Second).
		WithLogger(logger).
		Execute().
		Exit(0)

	if err := r.Error(); err != nil {
		return fmt.Errorf("test driver failed: %v", err)
	}

	results := helpers.ParseStructuredOutput(string(r.Result().Stdout))

	// Define all expected results
	tests := []struct {
		name     string
		expected string
		label    string
	}{
		// Test 1: Basic correctness
		{"basic_match", "True", "map_kernel([0,1,2,3]) + 10 matches expected"},
		// Test 2: Output values
		{"output_values", "10.0,11.0,12.0,13.0", "output values are 10.0,11.0,12.0,13.0"},
		// Test 3: Larger input
		{"larger_match", "True", "map_kernel works for size=8"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All E01 tests passed!")
	return nil
}
