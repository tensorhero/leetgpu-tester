package stages

import (
	"fmt"
	"time"

	"github.com/bootcraft-cn/leetgpu-tester/internal/helpers"
	"github.com/bootcraft-cn/tester-utils/runner"
	"github.com/bootcraft-cn/tester-utils/test_case_harness"
	"github.com/bootcraft-cn/tester-utils/tester_definition"
)

func e17LayernormTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "layernorm-kernel",
		Timeout:     30 * time.Second,
		TestFunc:    testE17Layernorm,
		CompileStep: compileStep("test_e17"),
	}
}

func testE17Layernorm(harness *test_case_harness.TestCaseHarness) error {
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
		{"basic_match", "True", "layernorm([1,2,3,4,5]) matches NumPy"},
		{"zero_mean", "True", "normalized rows have mean ≈ 0"},
		{"unit_var", "True", "normalized rows have variance ≈ 1"},
		{"gamma_beta", "True", "gamma/beta affine transform applied"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All E17 tests passed!")
	return nil
}
