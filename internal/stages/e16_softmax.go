package stages

import (
	"fmt"
	"time"

	"github.com/bootcraft-cn/leetgpu-tester/internal/helpers"
	"github.com/bootcraft-cn/tester-utils/runner"
	"github.com/bootcraft-cn/tester-utils/test_case_harness"
	"github.com/bootcraft-cn/tester-utils/tester_definition"
)

func e16SoftmaxTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "softmax-kernel",
		Timeout:     30 * time.Second,
		TestFunc:    testE16Softmax,
		CompileStep: compileStep("test_e16"),
	}
}

func testE16Softmax(harness *test_case_harness.TestCaseHarness) error {
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
		{"basic_match", "True", "softmax([1,2,3]) matches NumPy"},
		{"sums_to_one", "True", "each row sums to 1.0"},
		{"numerical_stability", "True", "large values ([1000,1001,1002]) no NaN/Inf"},
		{"batch_correct", "True", "multi-row batch computed independently"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All E16 tests passed!")
	return nil
}
