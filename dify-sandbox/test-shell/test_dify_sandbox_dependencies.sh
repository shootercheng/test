#!/bin/bash
set -eu

SCRIPT_DIR="$(dirname "$(realpath "$0")")"

source ${SCRIPT_DIR}/test_env.sh

echo "更新依赖 复制python_lib_path依赖:/v1/sandbox/dependencies/update"
curl -X POST ${REQUEST_HOST}/v1/sandbox/dependencies/update \
     -H "Content-Type: application/json" \
     -H "X-Api-Key: ${X_API_KEY}" \
     -d '{"language": "python3"}'

echo ""
echo "刷新依赖dependencies/python-requirements.txt依赖:/v1/sandbox/dependencies/refresh"
curl -X GET ${REQUEST_HOST}/v1/sandbox/dependencies/refresh \
     -H "Content-Type: application/json" \
     -H "X-Api-Key: ${X_API_KEY}" \
     -d '{"language": "python3"}'


echo ""
echo "获取已安装依赖:/v1/sandbox/dependencies"
curl -X GET ${REQUEST_HOST}/v1/sandbox/dependencies \
     -H "Content-Type: application/json" \
     -H "X-Api-Key: ${X_API_KEY}" \
     -d '{"language": "python3"}'

echo ""

