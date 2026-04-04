package stages

import (
	"fmt"
	"time"

	"github.com/bootcraft-cn/leetgpu-tester/internal/helpers"
	"github.com/bootcraft-cn/tester-utils/runner"
	"github.com/bootcraft-cn/tester-utils/test_case_harness"
	"github.com/bootcraft-cn/tester-utils/tester_definition"
)

func e14MatmulTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "matmul",
		Timeout:     30 * time.Second,
		TestFunc:    testE14Matmul,
		CompileStep: compileStep("test_e14"),
	}
}

func testE14Matmul(harness *test_case_harness.TestCaseHarness) error {
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
		{"basic_match", "True", "2×2 matmul correct"},
		{"identity_multiply", "True", "A @ I = A"},
		{"non_square", "True", "2×3 @ 3×4 → 2×4"},
		{"larger_match", "True", "8×8 multi-block matmul"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All E14 tests passed!")
	return nil
}
