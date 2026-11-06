# Pull request

Keep it short. A reviewer should be able to tell what changed and why without
opening five other files.

## Type of change

- [ ] Backend data plane (`cmd/asura-store/`)
- [ ] Window (`ui/`)
- [ ] Catalogue (`catalog/`)
- [ ] Build, CI or tests
- [ ] Documentation only

## What changed

Facts and files, not a narrative of the session.

## Why

The problem this solves. For a catalogue entry, what it installs and who
maintains it upstream.

## Related issue

Fixes #

## Definition of done

For anything that is not documentation:

- [ ] `gofmt -l cmd/` is clean
- [ ] `go vet ./...` and `go test -count=1 -race ./...` (in `cmd/asura-store/`) pass
- [ ] `node --test ui/lib/store.test.mjs` passes
- [ ] `qmllint` is clean — errors fail; the Quickshell/qs.Ui warnings are expected
- [ ] stdout still belongs to the JSON contract; errors are one sentence on stderr
- [ ] No `sudo`, and no privilege held by this process (polkit remains the seam)

For catalogue entries:

- [ ] The file carries `id`, `category`, `name` and `summary`; the id matches
      its kind's grammar and agrees with the folder it lives in
- [ ] For a plugin: a real, maintained `upstream` repo — an empty folder is not
      a plugin

## Mirror duty

`desktop/cmd/asura-store/`, `shell/store.qml` + `shell/store/` and
`desktop/catalog/` in [`mohitjoer/ASURA-OS`](https://github.com/mohitjoer/ASURA-OS)
are kept byte-identical with this repository (standalone `NOTICE.md`,
workflows, docs and the `ui/store.qml` entry point are excluded). If this PR
changes any of those trees, the mirror commit lands with it or directly after
it — CI refuses the drift either way.

## Licensing and sign-off

- [ ] Every commit carries a `Signed-off-by:` line (`git commit -s`)
- [ ] The contribution is under GPL-3.0, per [`LICENSE`](LICENSE) and the
      porting ledger in [`NOTICE.md`](NOTICE.md)
- [ ] No third-party code was pasted into the tree

## Notes for the reviewer

What you want looked at, what you are unsure about, and what you deliberately
left alone.
