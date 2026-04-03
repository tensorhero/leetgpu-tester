package stages

import (
	"fmt"
	"time"

	"github.com/tensorhero-cn/leetgpu-tester/internal/helpers"
	"github.com/tensorhero-cn/tester-utils/runner"
	"github.com/tensorhero-cn/tester-utils/test_case_harness"
	"github.com/tensorhero-cn/tester-utils/tester_definition"
)

func e09PoolingTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "pooling",
		Timeout:     30 * time.Second,
		TestFunc:    testE09Pooling,
		CompileStep: compileStep("test_e09"),
	}
}

func testE09Pooling(harness *test_case_harness.TestCaseHarness) error {
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
		{"basic_match", "True", "pooling_kernel output matches pool_spec"},
		{"window_sum", "3.0", "out[2] = a[0]+a[1]+a[2] = 3.0 (first full window)"},
		{"shared_memory_used", "True", "kernel uses cuda.shared.array"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All E09 tests passed!")
	return nil
}
