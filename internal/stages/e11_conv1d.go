package stages

import (
	"fmt"
	"time"

	"github.com/tensorhero-cn/leetgpu-tester/internal/helpers"
	"github.com/tensorhero-cn/tester-utils/runner"
	"github.com/tensorhero-cn/tester-utils/test_case_harness"
	"github.com/tensorhero-cn/tester-utils/tester_definition"
)

func e11Conv1dTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "conv1d",
		Timeout:     30 * time.Second,
		TestFunc:    testE11Conv1d,
		CompileStep: compileStep("test_e11"),
	}
}

func testE11Conv1d(harness *test_case_harness.TestCaseHarness) error {
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
		{"simple_match", "True", "simple convolution (single block)"},
		{"identity_conv", "True", "identity kernel=[1] → output equals input"},
		{"multi_block_match", "True", "multi-block convolution (size=15)"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All E11 tests passed!")
	return nil
}
