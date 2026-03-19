package stages

import (
	"github.com/tensorhero/tester-utils/tester_definition"
)

// GetDefinition returns the TesterDefinition for the leetgpu course.
func GetDefinition() tester_definition.TesterDefinition {
	return tester_definition.TesterDefinition{
		TestCases: []tester_definition.TestCase{
			// Phase 1: Thread & Memory Model
			e01MapTestCase(),
			e02ZipTestCase(),
			e03GuardTestCase(),
			e04Map2dTestCase(),
			e05BroadcastTestCase(),
			e06BlocksTestCase(),
			e07Blocks2dTestCase(),
			// Phase 2: Shared Memory
			e08SharedMemoryTestCase(),
			e09PoolingTestCase(),
			e10DotProductTestCase(),
			e11Conv1dTestCase(),
			e12PrefixSumTestCase(),
			e13AxisSumTestCase(),
			e14MatmulTestCase(),
			// Phase 3: ML Kernels
			e15ReluTestCase(),
			e16SoftmaxTestCase(),
			e17LayernormTestCase(),
			e18CrossEntropyTestCase(),
			e19AttentionTestCase(),
			e20TransposeTestCase(),
		},
	}
}

// pythonRule creates a LanguageRule for Python detection.
// testDriver is the module name without extension (e.g. "test_e01").
func pythonRule(testDriver string) tester_definition.LanguageRule {
	return tester_definition.LanguageRule{
		DetectFile: "leetgpu/__init__.py",
		Language:   "python",
		Source:     "leetgpu/__init__.py",
		RunCmd:     "python3",
		RunArgs:    []string{"tests/" + testDriver + ".py"},
	}
}

// compileStep returns a CompileStep that detects the leetgpu Python package.
func compileStep(testDriver string) *tester_definition.CompileStep {
	return &tester_definition.CompileStep{
		Language: "auto",
		AutoDetect: []tester_definition.LanguageRule{
			pythonRule(testDriver),
		},
	}
}
