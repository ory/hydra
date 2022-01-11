#!/bin/bash

set -Eeuxo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."

# shellcheck disable=SC2006
cat > configuration.md << EOF
---
id: configuration
title: Configuration
---

\`\`\`yaml
`cat ./docs/config.yaml`
\`\`\`

EOF
