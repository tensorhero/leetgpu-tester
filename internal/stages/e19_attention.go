package stages

import (
	"fmt"
	"time"

	"github.com/tensorhero/leetgpu-tester/internal/helpers"
	"github.com/tensorhero/tester-utils/runner"
	"github.com/tensorhero/tester-utils/test_case_harness"
	"github.com/tensorhero/tester-utils/tester_definition"
)

func e19AttentionTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "attention-kernel",
		Timeout:     30 * time.Second,
		TestFunc:    testE19Attention,
		CompileStep: compileStep("test_e19"),
	}
}

func testE19Attention(harness *test_case_harness.TestCaseHarness) error {
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
		{"basic_match", "True", "attention output matches NumPy reference"},
		{"identity_keys", "True", "K=Identity produces correct output"},
		{"attention_weights_sum", "True", "V=ones → out=ones (weights sum to 1)"},
		{"scaling_correct", "True", "scaling by 1/√d_k applied correctly"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All E19 tests passed!")
	return nil
}
