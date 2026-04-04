package stages

import (
	"fmt"
	"time"

	"github.com/bootcraft-cn/leetgpu-tester/internal/helpers"
	"github.com/bootcraft-cn/tester-utils/runner"
	"github.com/bootcraft-cn/tester-utils/test_case_harness"
	"github.com/bootcraft-cn/tester-utils/tester_definition"
)

func e15ReluTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "relu-kernel",
		Timeout:     30 * time.Second,
		TestFunc:    testE15Relu,
		CompileStep: compileStep("test_e15"),
	}
}

func testE15Relu(harness *test_case_harness.TestCaseHarness) error {
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
		{"basic_match", "True", "[-1,0,1,2] → [0,0,1,2]"},
		{"all_negative", "True", "all negative → all zeros"},
		{"all_positive", "True", "all positive → unchanged"},
		{"large_input", "True", "1024-element multi-block"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All E15 tests passed!")
	return nil
}
