#!/bin/bash
# 批量测试所有已实现 stage 的 solution
# 用法: ./scripts/test-all-solutions.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TESTER_DIR="$(dirname "$SCRIPT_DIR")"
SOLUTION_DIR="${TESTER_DIR}/../leetgpu-solution"
STARTER_DIR="${TESTER_DIR}/../leetgpu-starter"

# 构建 tester
cd "$TESTER_DIR"
go build -o leetgpu-tester .

# Stage 列表（按课程顺序）
STAGES=(
    # Phase 1: Thread & Memory Model
    "map"
    "zip"
    "guard"
    "map-2d"
    "broadcast"
    "blocks"
    "blocks-2d"
    # Phase 2: Shared Memory
    "shared-memory"
    "pooling"
    "dot-product"
    "conv1d"
    "prefix-sum"
    "axis-sum"
    "matmul"
    # Phase 3: ML Kernels
    "relu-kernel"
    "softmax-kernel"
    "layernorm-kernel"
    "cross-entropy-kernel"
    "attention-kernel"
    "transpose"
)

# 同步 test drivers 到 solution
mkdir -p "${SOLUTION_DIR}/tests"
cp -f "${STARTER_DIR}/tests/"*.py "${SOLUTION_DIR}/tests/" 2>/dev/null || true

PASSED=0
FAILED=0
SKIPPED=0
TOTAL_TIME=0

echo "=========================================="
echo "  LeetGPU Solution Tester"
echo "=========================================="
echo ""

for stage in "${STAGES[@]}"; do
    printf "🧪 [%-25s] Testing... " "$stage"

    start_time=$(python3 -c 'import time; print(time.time())')
    output=$(NUMBA_ENABLE_CUDASIM=1 ./leetgpu-tester -d="$SOLUTION_DIR" -s="$stage" 2>&1) && rc=0 || rc=$?
    end_time=$(python3 -c 'import time; print(time.time())')
    elapsed=$(python3 -c "print(f'{$end_time - $start_time:.2f}')")

    if echo "$output" | grep -q 'stage .* not found'; then
        echo "⏭️  SKIPPED - not yet implemented"
        ((SKIPPED++))
        continue
    fi

    if [ $rc -eq 0 ]; then
        echo "✅ PASSED (${elapsed}s)"
        ((PASSED++))
    else
        echo "❌ FAILED (${elapsed}s)"
        ((FAILED++))
    fi

    TOTAL_TIME=$(python3 -c "print(f'{$TOTAL_TIME + $elapsed:.2f}')")
done

echo ""
echo "=========================================="
echo "  Results: $PASSED passed, $FAILED failed, $SKIPPED skipped"
echo "  Total time: ${TOTAL_TIME}s"
echo "=========================================="

if [ $FAILED -gt 0 ]; then
    exit 1
fi
