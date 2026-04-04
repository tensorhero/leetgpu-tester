package stages

import (
	"fmt"
	"time"

	"github.com/bootcraft-cn/leetgpu-tester/internal/helpers"
	"github.com/bootcraft-cn/tester-utils/runner"
	"github.com/bootcraft-cn/tester-utils/test_case_harness"
	"github.com/bootcraft-cn/tester-utils/tester_definition"
)

func e18CrossEntropyTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "cross-entropy-kernel",
		Timeout:     30 * time.Second,
		TestFunc:    testE18CrossEntropy,
		CompileStep: compileStep("test_e18"),
	}
}

func testE18CrossEntropy(harness *test_case_harness.TestCaseHarness) error {
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
		{"basic_match", "True", "cross-entropy matches NumPy reference"},
		{"perfect_prediction", "True", "perfect prediction → loss ≈ 0"},
		{"uniform_logits", "True", "uniform logits → loss = log(num_classes)"},
		{"batch_independent", "True", "batch samples computed independently"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All E18 tests passed!")
	return nil
}
