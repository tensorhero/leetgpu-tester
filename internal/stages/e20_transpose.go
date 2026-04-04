package stages

import (
	"fmt"
	"time"

	"github.com/bootcraft-cn/leetgpu-tester/internal/helpers"
	"github.com/bootcraft-cn/tester-utils/runner"
	"github.com/bootcraft-cn/tester-utils/test_case_harness"
	"github.com/bootcraft-cn/tester-utils/tester_definition"
)

func e20TransposeTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "transpose",
		Timeout:     30 * time.Second,
		TestFunc:    testE20Transpose,
		CompileStep: compileStep("test_e20"),
	}
}

func testE20Transpose(harness *test_case_harness.TestCaseHarness) error {
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
		{"basic_match", "True", "4×4 square transpose correct"},
		{"non_square", "True", "3×5 → 5×3 non-square transpose"},
		{"identity_invariant", "True", "symmetric matrix unchanged"},
		{"double_transpose", "True", "double transpose = original"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All E20 tests passed!")
	return nil
}
