# NOTICE

The store's backend (`cmd/asura-store/`), window (`ui/`) and catalogue formats are a
port of **Ryoku**'s RyoStore — [`apps/ryostore`](https://github.com/ryoku-dev/ryoku)'s
`backend/` and `quickshell/`. Ryoku is GPL-3.0. ASURA is GPL-3.0 as well, so the
licence is compatible, and this notice is the honest ledge for where the code came
from.

- Source: Ryoku, `apps/ryostore/`, under GPL-3.0. The store in this tree is the
  backport of the Ryoku 1.0.7 era of that app; the reference checkout is Ryoku
  `f7197d72de8a21e44b0bea922bdd915f28f6b327` (2026-10-31).
- What came across: the Quickshell window — `ui/*.qml`, `ui/lib/store.js`,
  `ui/Singletons/` — from `apps/ryostore/quickshell/`, read by the pages with the same
  JSON contract; and the backend data plane in `cmd/asura-store/` — the catalogue
  compilation and its revision, the `check` refresh-dot probe, the window-manager
  gate, the crash-environment scrub, the providers' shapes and the progress-line
  contract — from `apps/ryostore/backend/`.
- What did not come across: Ryoku's remote registry, its receipt transactions, its
  extras actuator, its snapshot/cache reconciliation, its usage dashboards and its
  download plumbing for registry-hosted payloads. ASURA's store runs with no remote
  catalogue at all: software comes from the package manager, the checked-in
  `catalog/` JSON, and the plugins' own upstream repositories.
- Renames: `ryostore` → `asura-store`, Ryoku's `Ryo*` names → ASURA's, and the
  window's branding says "Asura". Library names inside the tree are the ones the
  imports resolve through.
- What is ASURA's own: the Ubuntu apt provider with its `asura-pkg-add` /
  `asura-pkg-remove` helper seam and cancel semantics, the xbps provider, the plugin
  provider's `asura plugin install` hand-off, the theme provider's activate-by-file
  model, the bar-style shelf and its `shell.json` selection, the two shell subcommands
  (`open`, `settings`) that speak to the ASURA session, the catalogue manifests in
  `catalog/`, and all of it downstream of the 1.0.7 backport.
- Files keep the copyright, SPDX and "ported from" lines they came with or were given.
  Where a port differs from Ryoku on purpose, the file says so in a comment.
- In the ASURA-OS distribution this store lives as byte-identical mirrors of
  `cmd/asura-store/`, `ui/` and `catalog/` under
  [`mohitjoer/ASURA-OS`](https://github.com/mohitjoer/ASURA-OS)
  (`desktop/cmd/asura-store/`, `shell/store.qml` + `shell/store/`, `desktop/catalog/`),
  whose `store-sync` workflow refuses any byte that drifts. This repository is the
  canonical copy.
