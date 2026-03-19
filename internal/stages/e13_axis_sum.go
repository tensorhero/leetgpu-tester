package stages

import (
	"fmt"
	"time"

	"github.com/tensorhero/leetgpu-tester/internal/helpers"
	"github.com/tensorhero/tester-utils/runner"
	"github.com/tensorhero/tester-utils/test_case_harness"
	"github.com/tensorhero/tester-utils/tester_definition"
)

func e13AxisSumTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "axis-sum",
		Timeout:     30 * time.Second,
		TestFunc:    testE13AxisSum,
		CompileStep: compileStep("test_e13"),
	}
}

func testE13AxisSum(harness *test_case_harness.TestCaseHarness) error {
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
		{"basic_match", "True", "[[1,2],[3,4]] → [4,6]"},
		{"single_row", "True", "[[5,3,7]] → [5,3,7]"},
		{"larger_match", "True", "4×3 column sums match numpy"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All E13 tests passed!")
	return nil
}
