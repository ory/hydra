#!/bin/bash

set -Eeuxo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."

cat > configuration.md << EOF
---
id: configuration
title: Configuration
---

\`\`\`
`cat ./docs/config.yaml`
\`\`\`

EOF
