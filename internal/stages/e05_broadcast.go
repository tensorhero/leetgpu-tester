package stages

import (
	"fmt"
	"time"

	"github.com/bootcraft-cn/leetgpu-tester/internal/helpers"
	"github.com/bootcraft-cn/tester-utils/runner"
	"github.com/bootcraft-cn/tester-utils/test_case_harness"
	"github.com/bootcraft-cn/tester-utils/tester_definition"
)

func e05BroadcastTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "broadcast",
		Timeout:     30 * time.Second,
		TestFunc:    testE05Broadcast,
		CompileStep: compileStep("test_e05"),
	}
}

func testE05Broadcast(harness *test_case_harness.TestCaseHarness) error {
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
		{"basic_match", "True", "broadcast_kernel a=[0,1] + b=[0,1] matches expected"},
		{"guard_broadcast", "True", "out-of-bounds threads (3x3 grid, 2x2 data) did not write"},
		{"output_shape", "2,2", "output shape is (2, 2)"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All E05 tests passed!")
	return nil
}
