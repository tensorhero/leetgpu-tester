#!/bin/bash
# 批量测试所有 stage 的 solution（numba）
# 用法: ./scripts/test-all-solutions.sh
#
# 分支模型：solution 仓库每种语言一个分支（numba），
# 脚本通过 git worktree 将各分支 checkout 到临时目录中测试。

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

# 语言列表
LANGUAGES=("numba")

PASSED=0
FAILED=0
SKIPPED=0
TOTAL_TIME=0

echo "=========================================="
echo "  LeetGPU Solution Tester"
echo "=========================================="
echo ""

for lang in "${LANGUAGES[@]}"; do
    echo "--- Language: ${lang} ---"
    echo ""

    # 使用 git worktree 将语言分支 checkout 到临时目录
    worktree_dir="${SOLUTION_DIR}/.worktree-${lang}"
    if [ -d "$worktree_dir" ]; then
        git -C "$SOLUTION_DIR" worktree remove --force "$worktree_dir" 2>/dev/null || rm -rf "$worktree_dir"
    fi
    git -C "$SOLUTION_DIR" worktree add "$worktree_dir" "$lang" 2>/dev/null

    sol_dir="$worktree_dir"

    if [ ! -d "$sol_dir" ]; then
        echo "⏭️  [${lang}] SKIPPED - branch not found"
        ((SKIPPED += ${#STAGES[@]}))
        echo ""
        continue
    fi

    # Ensure test drivers are present in solution dir (copy from starter branch)
    starter_worktree="${STARTER_DIR}/.worktree-${lang}"
    if [ -d "$starter_worktree" ]; then
        git -C "$STARTER_DIR" worktree remove --force "$starter_worktree" 2>/dev/null || rm -rf "$starter_worktree"
    fi
    git -C "$STARTER_DIR" worktree add "$starter_worktree" "$lang" 2>/dev/null

    mkdir -p "${sol_dir}/tests"
    cp -f "${starter_worktree}/tests/"*.py "${sol_dir}/tests/" 2>/dev/null || true

    for stage in "${STAGES[@]}"; do
        printf "🧪 [%-25s %8s] Testing... " "$stage" "$lang"

        start_time=$(python3 -c 'import time; print(time.time())')
        output=$(NUMBA_ENABLE_CUDASIM=1 ./leetgpu-tester -d="$sol_dir" -s="$stage" 2>&1) && rc=0 || rc=$?
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

    # 清理 worktree
    git -C "$SOLUTION_DIR" worktree remove --force "$worktree_dir" 2>/dev/null || rm -rf "$worktree_dir"
    git -C "$STARTER_DIR" worktree remove --force "$starter_worktree" 2>/dev/null || rm -rf "$starter_worktree"

    echo ""
done

echo "=========================================="
echo "  Results: $PASSED passed, $FAILED failed, $SKIPPED skipped"
echo "  Total time: ${TOTAL_TIME}s"
echo "=========================================="

if [ $FAILED -gt 0 ]; then
    exit 1
fi
