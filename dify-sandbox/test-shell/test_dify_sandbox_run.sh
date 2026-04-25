#!/bin/bash
set -eu

SCRIPT_DIR="$(dirname "$(realpath "$0")")"

source ${SCRIPT_DIR}/test_env.sh

# 健康检查接口
echo "健康检查接口:/health"
curl ${REQUEST_HOST}/health

# Python3 示例
echo ""
echo "Python3 运行代码:/v1/sandbox/run"
curl -X POST ${REQUEST_HOST}/v1/sandbox/run \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: ${X_API_KEY}" \
  -d '{
    "language": "python3",
    "code": "import json\nperson = {\"name\": \"John\", \"age\": 30, \"city\": \"New York\"}\njson_str = json.dumps(person)\nprint(json_str)",
    "preload": "",
    "enable_network": false
  }'

echo ""
echo "Node.js 运行代码:/v1/sandbox/run"
# Node.js 示例
curl -X POST ${REQUEST_HOST}/v1/sandbox/run \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: ${X_API_KEY}" \
  -d '{
    "language": "nodejs",
    "code": "const person = {name: \"John\", age: 30, city: \"New York\"};\nconst jsonString = JSON.stringify(person);\nconsole.log(jsonString);",
    "preload": "",
    "enable_network": false
  }'

echo ""
