package stages

import (
	"fmt"
	"time"

	"github.com/tensorhero-cn/leetgpu-tester/internal/helpers"
	"github.com/tensorhero-cn/tester-utils/runner"
	"github.com/tensorhero-cn/tester-utils/test_case_harness"
	"github.com/tensorhero-cn/tester-utils/tester_definition"
)

func e06BlocksTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "blocks",
		Timeout:     30 * time.Second,
		TestFunc:    testE06Blocks,
		CompileStep: compileStep("test_e06"),
	}
}

func testE06Blocks(harness *test_case_harness.TestCaseHarness) error {
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
		{"basic_match", "True", "blocks_kernel output matches expected (a + 10)"},
		{"guard_blocks", "True", "overflow threads (9-11) did not write past size"},
		{"all_elements_processed", "True", "all 9 elements processed correctly"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All E06 tests passed!")
	return nil
}
