# Contributing to Asura Store

Thank you for pointing at the code. This repository is small on purpose — one
Go data plane, one Quickshell window, a folder of catalogue JSON — so the
rules here are short, and most of them are what the CI gate
(`.github/workflows/ci.yml`) already refuses to let through.

## Licensing first

- This work is **GPL-3.0** ([`LICENSE`](LICENSE)). Every contribution is under
  GPL-3.0; never paste third-party code into the tree.
- The store is a derivative of **Ryoku**'s RyoStore (GPL-3.0) — see
  [`NOTICE.md`](NOTICE.md). Files keep any copyright/SPDX/ported-from lines
  they carry. If you port something new from Ryoku, extend the ledger there
  rather than silently growing the derivative.
- Every commit carries a `Signed-off-by:` line (`git commit -s`), so the
  Developer Certificate of Origin says the patch is yours or yours to pass on.

## Environment

- **Go** — the toolchain `go.mod` pins. The module is standard library only:
  there is no `go.sum` and the go module cache is off in CI. Do not add
  dependencies without an argument for them.
- **Node.js 22+** — for the window's JavaScript tests (`node --test`, no
  npm tree).
- **Quickshell** — the window runs under the OS session: the entry file is
  written to sit at the quickshell root, beside the shell's `store/`
  directory and its `qs.Ui` kit (see `docs/store.md`). It is not a runnable
  config from this repository alone.
- **`jq`** — most CI scripts answer through it.
- **For `qmllint`** — the tool and the two modules the window imports:

  ```bash
  sudo apt-get install qt6-declarative-dev-tools qt6-declarative-dev \
    qml6-module-qtquick qml6-module-qtquick-effects
  ```

## Everyday checks

The same set CI runs, in commands a human can type:

```bash
make build                                        # bin/asura-store
cd cmd/asura-store
gofmt -l .                                        # must print nothing
go vet ./...
go test -count=1 -race ./...
node --test ui/lib/store.test.mjs                 # from the repo root
/usr/lib/qt6/bin/qmllint ../../ui/*.qml           # errors fail; Quickshell/qs.Ui warnings are expected
```

A change the CLI smoke would exercise (the refresh-dot contract, the catalogue
shape, an error path) is best proven with the smoke set in the workflow's
`backend` job, or by hand: `asura-store check`, `catalog`, `search`, and the
refusals, which must each be one sentence on stderr and exit 1.

## The rules that make review quick

- **stdout belongs to the JSON contract.** Reading commands answer exactly
  one JSON document per invocation; install/remove print `{"stage","percent"}`
  lines. Anything chatty goes to stderr, and an error is one sentence a
  window could show a person.
- **The `Item` contract is shared with Ryoku's window pages.** Field names and
  JSON tags are the compatibility surface — add fields if you must, never
  rename or repurpose one. `cmd/asura-store/model.go` is the source of truth;
  keep its comments honest.
- **Paths honour the environment** — `XDG_CONFIG_HOME`, `XDG_DATA_HOME`,
  `XDG_CACHE_HOME` — and nothing user-facing lives outside them.
- **The store never holds a privilege.** Escalation is polkit, through the
  OS's `asura-pkg-add`/`asura-pkg-remove` helpers or `pkexec xbps-*`. No
  `sudo` anywhere, ever.
- **Catalogue entries** follow `docs/store.md` ("Authoring catalogue
  entries"): the loaders silently skip a file that does not parse or whose id
  misses the grammar, so the data must stay clean — CI validates it loudly so
  a rotting manifest never disappears quietly.

## Commits

Conventional style, scope included, imperative subject, lowercase, body when
the why is not obvious:

```
feat(store): gate catalogue items by declared window manager
fix(plugins): refuse an entry with no upstream instead of succeeding empty
ci(store): give qmllint its type data
docs(store): document the categories and the revision probe
```

A commit that only reflows formatting gets none of the gate's time — gofmt
drift is caught before vet, so keep the tree formatted and the diffs honest.

## The mirrors

The distribution keeps byte-identical copies of this repository:
`desktop/cmd/asura-store/`, `shell/store.qml` + `shell/store/` and
`desktop/catalog/` in [`mohitjoer/ASURA-OS`](https://github.com/mohitjoer/ASURA-OS),
where a `store-sync` workflow diffs them against this repo's `HEAD` and fails
on any drift (plus a daily catch-up, so a forgotten mirror turns red within a
day). This repository is the canonical copy: land your change here first, then
mirror it (or ask a maintainer to). Documentation, workflows and
`ui/store.qml` as the window's entry point are standalone-only and excluded
from the mirror duty.

## Reporting problems

Bugs through the issue template, security through private reporting
([`SECURITY.md`](SECURITY.md)) — a privilege-escalation path or a shell
injection through an id is the worst this code can do, so those reports can
never be public issues first.
