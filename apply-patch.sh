#!/bin/bash

set -euo pipefail

go mod vendor

patch -p0 --fuzz=3 < openai-client.patch