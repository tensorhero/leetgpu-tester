package stages

import (
	"fmt"
	"time"

	"github.com/bootcraft-cn/leetgpu-tester/internal/helpers"
	"github.com/bootcraft-cn/tester-utils/runner"
	"github.com/bootcraft-cn/tester-utils/test_case_harness"
	"github.com/bootcraft-cn/tester-utils/tester_definition"
)

func e02ZipTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "zip",
		Timeout:     30 * time.Second,
		TestFunc:    testE02Zip,
		CompileStep: compileStep("test_e02"),
	}
}

func testE02Zip(harness *test_case_harness.TestCaseHarness) error {
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
		{"basic_match", "True", "zip_kernel([0,1,2,3], [0,1,2,3]) matches expected"},
		{"output_values", "0.0,2.0,4.0,6.0", "output values are 0.0,2.0,4.0,6.0"},
		{"negative_match", "True", "zip_kernel works with negative values"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All E02 tests passed!")
	return nil
}
