#!/bin/bash

set -Eeuxo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."

cat > appendix.md << EOF
---
id: hydra-appendix
title: Appendix
---

## \`hydra serve\`

\`\`\`
\$ hydra help serve

`hydra help serve`
\`\`\`

EOF
