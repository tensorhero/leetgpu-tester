package stages

import (
	"fmt"
	"time"

	"github.com/bootcraft-cn/leetgpu-tester/internal/helpers"
	"github.com/bootcraft-cn/tester-utils/runner"
	"github.com/bootcraft-cn/tester-utils/test_case_harness"
	"github.com/bootcraft-cn/tester-utils/tester_definition"
)

func e07Blocks2dTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "blocks-2d",
		Timeout:     30 * time.Second,
		TestFunc:    testE07Blocks2d,
		CompileStep: compileStep("test_e07"),
	}
}

func testE07Blocks2d(harness *test_case_harness.TestCaseHarness) error {
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
		{"basic_match", "True", "blocks_2d_kernel output matches expected (a + 10)"},
		{"guard_2d_blocks", "True", "overflow threads in 6x6 grid did not write past 5x5"},
		{"all_elements_processed", "True", "all 25 elements processed correctly"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All E07 tests passed!")
	return nil
}
