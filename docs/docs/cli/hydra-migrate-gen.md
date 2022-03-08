---
id: hydra-migrate-gen
title: hydra migrate gen
description: hydra migrate gen Generate migration files from migration templates
---

<!--
This file is auto-generated.

To improve this file please make your change against the appropriate "./cmd/*.go" file.
-->

## hydra migrate gen

Generate migration files from migration templates

```
hydra migrate gen &lt;/source/path&gt; &lt;/target/path&gt; [flags]
```

### Options

```
      --dialects strings   Expect migrations for these dialects and no others to be either explicitly defined, or to have a generic fallback. &#34;&#34; disables dialect validation. (default [sqlite,cockroach,mysql,postgres])
  -h, --help               help for gen
```

### Options inherited from parent commands

```
  -c, --config strings   Path to one or more .json, .yaml, .yml, .toml config files. Values are loaded in the order provided, meaning that the last config file overwrites values from the previous config file.
```

### SEE ALSO

- [hydra migrate](hydra-migrate) - Various migration helpers
