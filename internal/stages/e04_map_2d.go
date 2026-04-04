package stages

import (
	"fmt"
	"time"

	"github.com/bootcraft-cn/leetgpu-tester/internal/helpers"
	"github.com/bootcraft-cn/tester-utils/runner"
	"github.com/bootcraft-cn/tester-utils/test_case_harness"
	"github.com/bootcraft-cn/tester-utils/tester_definition"
)

func e04Map2dTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "map-2d",
		Timeout:     30 * time.Second,
		TestFunc:    testE04Map2d,
		CompileStep: compileStep("test_e04"),
	}
}

func testE04Map2d(harness *test_case_harness.TestCaseHarness) error {
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
		{"basic_match", "True", "map_2d_kernel on 2x2 matrix matches expected"},
		{"guard_2d", "True", "out-of-bounds threads (3x3 grid, 2x2 data) did not write"},
		{"output_shape", "2,2", "output shape is (2, 2)"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All E04 tests passed!")
	return nil
}
