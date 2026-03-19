package stages

import (
	"fmt"
	"time"

	"github.com/tensorhero/leetgpu-tester/internal/helpers"
	"github.com/tensorhero/tester-utils/runner"
	"github.com/tensorhero/tester-utils/test_case_harness"
	"github.com/tensorhero/tester-utils/tester_definition"
)

func e08SharedMemoryTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "shared-memory",
		Timeout:     30 * time.Second,
		TestFunc:    testE08SharedMemory,
		CompileStep: compileStep("test_e08"),
	}
}

func testE08SharedMemory(harness *test_case_harness.TestCaseHarness) error {
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
		{"basic_match", "True", "shared_memory_kernel output matches expected (a + 10)"},
		{"shared_used", "True", "kernel uses cuda.shared.array"},
		{"sync_present", "True", "kernel calls cuda.syncthreads()"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All E08 tests passed!")
	return nil
}
