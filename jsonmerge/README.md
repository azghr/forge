# jsonmerge

Recursively merge and diff JSON-like data (maps with `interface{}` values).
Designed for combining configurations or comparing state snapshots.

## Problem

Applications often need to layer configurations from multiple sources
(defaults + overrides, or base + environment-specific settings). The standard
library provides `encoding/json` for marshalling but no way to deep-merge maps
or compute structural diffs. This package fills that gap with zero dependencies.

## Quick start

```go
import "github.com/azghr/forge/jsonmerge"

a := map[string]interface{}{"x": 1, "y": map[string]interface{}{"v": 2}}
b := map[string]interface{}{"y": map[string]interface{}{"v": 3}, "z": 4}

jsonmerge.Merge(a, b)
// a == map[x:1 y:map[v:3] z:4]

diff := jsonmerge.Diff(a, b)
// diff == ["y.v", "z"]
```

## API

### Functions

- **`Merge(dst, src map[string]interface{}, opts ...Option)`** —
  recursively merge `src` into `dst`. Map values are merged recursively;
  non-map values from `src` override `dst`. Slices are replaced or appended
  depending on the option.
- **`Diff(a, b map[string]interface{}) []string`** — return dot-notation paths
  for every key in `a` whose value differs from the corresponding key in `b`.

### Options

- **`WithSliceMode(m SliceMode) Option`** — controls slice merge behaviour:
  - `SliceReplace` (default) — replace destination slice with source slice.
  - `SliceAppend` — append source slice elements to destination slice.

## Performance

| Function | Time | Space | Notes |
|----------|------|-------|-------|
| Merge    | O(n) | O(1)  | In-place on dst; no extra allocations |
| Diff     | O(n) | O(n)  | Allocates result slice only |

Both functions walk maps recursively and are concurrency-safe (no shared state).
