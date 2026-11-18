# NOTICE

The store's backend (`cmd/asura-store/`), window (`ui/`) and catalogue formats are
derived under GPL-3.0. ASURA is GPL-3.0 as well, so the
licence is compatible, and this notice is the honest record for where the code came
from.

- Source: upstream store components under GPL-3.0.
- What came across: the Quickshell window — `ui/*.qml`, `ui/lib/store.js`,
  `ui/Singletons/` — read by the pages with the same
  JSON contract; and the backend data plane in `cmd/asura-store/` — the catalogue
  compilation and its revision, the `check` refresh-dot probe, the window-manager
  gate, the crash-environment scrub, the providers' shapes and the progress-line
  contract.
- What did not come across: remote registry, receipt transactions,
  extras actuator, snapshot/cache reconciliation, usage dashboards and
  download plumbing for registry-hosted payloads. ASURA's store runs with no remote
  catalogue at all: software comes from the package manager, the checked-in
  `catalog/` JSON, and the plugins' own upstream repositories.
- Renames: CLI binary is `asura-store`, and the
  window's branding says "Asura". Library names inside the tree are the ones the
  imports resolve through.
- What is ASURA's own: the xbps provider (an Ubuntu apt provider with an
  `asura-pkg-add` / `asura-pkg-remove` helper seam was removed when ASURA moved to a
  Void Linux base), the plugin
  provider's `asura plugin install` hand-off, the theme provider's activate-by-file
  model, the bar-style shelf and its `shell.json` selection, the two shell subcommands
  (`open`, `settings`) that speak to the ASURA session, the catalogue manifests in
  `catalog/`, and all of it downstream of the 1.0.7 backport.
- Files keep the copyright and SPDX lines they came with or were given.
  Where a port differs on purpose, the file says so in a comment.
- In the ASURA-OS distribution this store lives as byte-identical mirrors of
  `cmd/asura-store/`, `ui/` and `catalog/` under
  [`mohitjoer/ASURA-OS`](https://github.com/mohitjoer/ASURA-OS)
  (`desktop/cmd/asura-store/`, `shell/store.qml` + `shell/store/`, `desktop/catalog/`),
  whose `store-sync` workflow refuses any byte that drifts. This repository is the
  canonical copy.
