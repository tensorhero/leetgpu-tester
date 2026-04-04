package stages

import (
	"fmt"
	"time"

	"github.com/bootcraft-cn/leetgpu-tester/internal/helpers"
	"github.com/bootcraft-cn/tester-utils/runner"
	"github.com/bootcraft-cn/tester-utils/test_case_harness"
	"github.com/bootcraft-cn/tester-utils/tester_definition"
)

func e12PrefixSumTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "prefix-sum",
		Timeout:     30 * time.Second,
		TestFunc:    testE12PrefixSum,
		CompileStep: compileStep("test_e12"),
	}
}

func testE12PrefixSum(harness *test_case_harness.TestCaseHarness) error {
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
		{"basic_match", "True", "[1,2,3,4] → [1,3,6,10]"},
		{"all_ones", "True", "[1,1,...,1] → [1,2,...,8]"},
		{"single_element", "True", "[5] → [5]"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All E12 tests passed!")
	return nil
}
