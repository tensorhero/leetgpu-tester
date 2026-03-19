# LeetGPU Tester

Automated testing tool for the LeetGPU course.

## Option 1: Build from Source

```bash
git clone https://github.com/tensorhero/leetgpu-tester
cd leetgpu-tester
go build .
NUMBA_ENABLE_CUDASIM=1 ./leetgpu-tester -s map -d ~/my-solution
```

**Dependencies:** Go 1.24+, Python 3.11+, Numba, NumPy

## Option 2: Docker Image

**Quick Start**

```bash
cd ~/my-solution  # your solution root (contains leetgpu/ and tests/)
docker pull ghcr.io/tensorhero/leetgpu-tester:latest
docker run --rm --user $(id -u):$(id -g) -v "$(pwd):/workspace" ghcr.io/tensorhero/leetgpu-tester:latest -s map -d /workspace
```

**Simplified script (recommended)**

Create `test.sh` in your solution root:

```bash
#!/bin/bash
docker run --rm --user $(id -u):$(id -g) -v "$(pwd):/workspace" ghcr.io/tensorhero/leetgpu-tester:latest \
  -s "${1:-map}" -d /workspace
```

Usage: `chmod +x test.sh && ./test.sh softmax-kernel`

**Local build (optional)**

```bash
git clone https://github.com/tensorhero/leetgpu-tester
cd leetgpu-tester
docker build -t my-tester .
# Usage: docker run --rm --user $(id -u):$(id -g) -v ~/my-solution:/workspace my-tester -s map -d /workspace
```

## License

MIT
