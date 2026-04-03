package stages

import (
	"fmt"
	"time"

	"github.com/tensorhero-cn/leetgpu-tester/internal/helpers"
	"github.com/tensorhero-cn/tester-utils/runner"
	"github.com/tensorhero-cn/tester-utils/test_case_harness"
	"github.com/tensorhero-cn/tester-utils/tester_definition"
)

func e10DotProductTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "dot-product",
		Timeout:     30 * time.Second,
		TestFunc:    testE10DotProduct,
		CompileStep: compileStep("test_e10"),
	}
}

func testE10DotProduct(harness *test_case_harness.TestCaseHarness) error {
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

	const atol = 1e-4
	floatTests := []struct {
		name     string
		expected float64
		label    string
	}{
		{"basic_match", 32.0, "[1,2,3]·[4,5,6] = 32.0"},
		{"orthogonal", 0.0, "orthogonal vectors → 0.0"},
		{"single_element", 21.0, "single element: 7×3 = 21.0"},
	}

	for _, tc := range floatTests {
		if err := helpers.AssertFloatResultClose(results, tc.name, tc.expected, atol); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All E10 tests passed!")
	return nil
}
