# The store

This repository is the standalone home of ASURA's software store: the Quickshell
window (`ui/`), the `asura-store` command (`cmd/asura-store/`) and the local
catalogue (`catalog/`). It is a port of Ryoku's RyoStore — the pieces and the
ledger live in [NOTICE](NOTICE.md). In the distribution itself the trees are
byte-identical mirrors of this one (`desktop/cmd/asura-store/`,
`shell/store.qml` + `shell/store/`, `desktop/catalog/`); this copy is canonical.

| Piece | Where | What it does |
|---|---|---|
| `asura-store` | `cmd/asura-store` | The data plane: catalogue, search, check, install, remove, open, settings. |
| The window | `ui/store.qml` and the components beside it (`ui/*.qml`, `ui/Singletons/`, `ui/lib/store.js`) | RyoStore's window, ported as it was; the logic sits beside its `node --test` suite. |
| `catalog/` | `catalog/<kind>/*.json` | Local catalogue manifests, staged by `make install` to `/usr/share/asura/store/catalog/`. A checkout update changes the shelf without a rebuild. |

On disk the installed layout is `/usr/bin/asura-store` plus
`/usr/share/asura/store/{ui,catalog}`. The catalogue loader also accepts a
checkout (it walks up from the working directory and the executable looking for
`catalog/<kind>`), which is what makes `go test` work without an install.

The window's entry file is byte-identical with the distribution's
`shell/store.qml`: it is written to sit at a quickshell root, beside the
shell's `store/` directory and its `qs.Ui` kit. In this repository it is
source — mirrored into the distribution — and the runnable arrangement lives
there; `asura-store open discover` starts the window through the OS session
wrapper.

## The commands

```text
asura-store catalog [--refresh] [--category ID]   JSON: categories and items
asura-store search QUERY                          JSON: {"items": [...]} from the archive
asura-store check                                 JSON: {"revision","updateAvailable","offline"}
asura-store install CATEGORY ID                   progress lines on stdout, an error on stderr
asura-store remove CATEGORY ID
asura-store open SECTION                          discover | library | a window section
asura-store settings CATEGORY [ID]                opens Settings (asura-shell hub open) on its page
```

Every reading command answers JSON on stdout and nothing else. Every failing
invocation answers one sentence on stderr and exits 1; `main` prefixes it with
`asura-store: ` for a terminal, the window prints the sentence as it comes.
The [JSON contract](../cmd/asura-store/model.go) is the source of truth for
field names — `Catalog`, `Category`, `Item` — and the window reads them
unchanged, every page in the same shape.

`catalog` has four categories. `apps` is a shelf of about thirty seed packages
with summary, description, version and home page read from apt (`xbps-query` on
a Void host), and the installed read from dpkg (`xbps-query -l`). `barstyles`
lists the shell's three built-in bar styles — `qsbar`, `kairos`, `sumi` — and
the one `shell.json` selects. `plugins`
reads the local catalogue under `catalog/plugins/`: each entry carries an
`upstream` repo, and installed means a folder in `~/.local/share/asura/plugins`.
`themes` merges the built-ins with the `catalog/themes/` entries and any theme
directory on disk, and marks active the one the current-theme file names.
Any other category is an error — `rices`, `lockscreens`, `bundles` and the rest
exist only as window sections, with no catalogue behind them.

Flags: `--refresh` bypasses any fresh in-process cache, `--category` builds one
category alone. A full catalogue with items acknowledges its revision to disk; a
category-scoped query never acknowledges — it is a read, not a viewing.

## How an install happens

The window runs `asura-store install CATEGORY ID`. The store prints one
`{"stage","percent"}` JSON line per step on stdout — a negative percent means
the step has no number — and one sentence on stderr, which the window shows.

`apps` is native packages, and its privilege path is the host's:

1. On an Ubuntu host the `asura-pkg-add` helper exists. `asura-store` runs
   `asura-pkg-add --status vlc`, and before any prompt the helper simulates the
   change with `apt-get -s` and refuses two things: a change that would also
   remove other packages, and a package that only installs a snap (Ubuntu's
   `firefox`, `chromium-browser` and `thunderbird` are snap launchers, and ASURA
   does not use snap). The refusal is a sentence on stderr.
2. The helper re-runs itself through `pkexec`. polkit asks the session's
   authentication agent, which is `asura-shell`, for the password. Nothing in
   the store runs `sudo`.
3. As root the helper runs `apt-get install` with apt's machine status on
   stdout, and `asura-store` turns it into one `{"stage":"DOWNLOADING","percent":42}`
   line per step. The stages are `AUTHORIZING`, `DOWNLOADING`, then `INSTALLING`
   (or `REMOVING`).

On a Void host there is no helper: `pkexec xbps-install -y` (or
`xbps-remove -y`) runs directly, and the store brackets it with `AUTHORIZING`
and `INSTALLING` (or `REMOVING`), since xbps cannot number one transaction. A
declined prompt, or no agent to ask, ends `pkexec` with 126 or 127, and the
store reports the same declined-authorization sentence as the helper path.

Cancel is honoured until dpkg starts. The caller cannot signal a root process,
so `asura-store` writes `cancel` to the helper's stdin, which stops apt-get
only while it is downloading and only if no dpkg is running. After that the
transaction finishes, because stopping apt-get under dpkg leaves a
half-configured system. `asura-store` also ignores `SIGPIPE` and `SIGHUP`, so
closing the window never kills a transaction partway — only SIGTERM and SIGINT
(the window's cancel) do.

`plugins` installs nothing itself. The catalogue entry names the upstream repo,
and the store runs `asura plugin install <upstream> --yes` (remove runs
`asura plugin remove <id>`), exactly what the Add-ons page uses. Without that
link there is no payload — an empty folder is not a plugin — so the store
refuses, and its error names the command that would work.

`themes` is applied, not downloaded. Installing a theme writes its name into
`~/.config/asura/current/theme.name` (`XDG_CONFIG_HOME` honoured; the older
bare `~/.config/asura/theme.name` is still read), which the shell's theme
engine picks up live. Only a theme whose files are on disk can be activated —
a catalogue-only entry stays listed as available. A removal refuses the
active theme, then deletes the theme's directories.

`barstyles` is all shelf: the styles are built into the shell, so install and
remove refuse.

## The revision and the refresh dot

Every freshly built catalogue gets a revision: a SHA-256 over the sorted
identity of its items — category, id, version, manifest digest, the declared
window manager, a digest of declared package managers, the paused state and the
declared upstream links — everything an upstream change can alter, and nothing
volatile (local install state, timestamps). Viewing the catalogue acknowledges
the revision by writing it to `~/.cache/asura-store/seen-revision`
(`XDG_CACHE_HOME` honoured), and only when it shows items — an empty first
view, which is an offline machine, never becomes the seen baseline.

`check` builds the catalogue again and answers
`{"revision","updateAvailable","offline"}`. It never writes the seen revision:
`updateAvailable` is true only when a revision has been seen at all, the
machine is not offline, and the revision differs — which is exactly when the
refresh dot lights.

`offline` is the catalogue's honesty flag: it is raised when a source failed to
load (or serves only a cache) and nothing in the answer proves that an upstream
was reachable. A provider that fails spoils only its own category record —
its `offline`, `cachedAt` and `error` — while the others still render.

## The window-manager gate

A catalogue item may name the window manager it was authored against
(`requiredWindowManager`; Ryoku's manifest carried the same idea). The store
asks the running provider once per launch — `asura wm status`, whose
`Provider:` line is the answer — and compares names; a third provider needs no
change in the store. An item named for another manager stays listed and keeps
any installed copy removable, but install and update are refused, so a control
that cannot work is never offered as if it could. An item that names nothing
runs anywhere, and a machine whose provider cannot be read gates nothing.

## Search

`search` works on whichever archive exists. On an apt host it reads
`apt-cache search`, keeps the hits where every word is in the name or the
one-line summary, leaves out libraries, headers and debug packages, and ranks
an exact name first, then a prefix, then a substring, ties by shorter name. On
a Void host it reads `xbps-query -Rs` and applies the same ranking, folding in
what is installed. The query gets 20 seconds; the window runs it after a typed
query pauses and shows the hits with the shelf. An archive that answers nothing
answers `{"items": []}`, not an error.

## Open and settings

`open` starts the store window under its lock (`/usr/share/asura/session/asura-store`,
the session wrapper the OS ships) and calls `nav open` on it through quickshell
ipc, retrying for 3 s while the window settles; a second instance exits at the
lock, so the call simply reaches the first. It accepts `discover`, `library`
and the window's sections — `apps`, `barstyles`, `themes`, `colorschemes`, ...
(the names beyond the four categories are the window's own routes, empty until
a catalogue exists).

`settings` maps a category to the Settings page that manages it (`themes` and
`rices` to appearance, `plugins` and `bundles` to add-ons, `barstyles` to Bar
Studio, `fastfetch` to fastfetch) and asks `asura-shell hub open` for that
page; a category without one is an error.

## Authoring catalogue entries

One JSON file per item under `catalog/<kind>/`. Fictional files are skipped
silently by the loaders — a manifest that does not parse or whose id misses the
grammar vanishes from its shelf — which is why CI validates the data loudly
(the `catalog` job). A healthy entry carries at least:

```json
{
  "id": "weather",
  "category": "plugins",
  "name": "Weather",
  "summary": "Forecast island for the bar",
  "description": "...",
  "author": "...",
  "version": "0.1.0",
  "upstream": "https://github.com/USER/asura-plugin-weather"
}
```

For a plugin the `upstream` line is load-bearing: without it the entry has no
payload and install refuses, naming the command that would work. Grammars,
enforced at load:

| Kind | id grammar |
|---|---|
| `plugins` | `^[a-z0-9][a-z0-9-]*$` — what `asura plugin` enforces, so an id can never escape its directory or read as an option |
| `themes`, `barstyles` | `^[a-z0-9][a-z0-9._-]*$` |

Optional gates on any entry: `requiredWindowManager` (the gate above),
`requiredPackageManagers` (e.g. `["apt"]`), `downloadPaused` with
`downloadPauseReason` (the item stays listed and removable, marked UNDER
CONSTRUCTION). Items honor the wider `Item` contract — tags, accent, surface,
kind, licence, screenshots — when a downstream page wants them.

`barstyles` manifests are on record only: the backend itself lists the shell's
three built-in styles (`qsbar`, `kairos`, `sumi`), so a bar style entry is
metadata, not shelf content.

The catalogue ships empty of fiction today: `plugins/` and `themes/` carry no
entries at all, because a shelf entry is a promise — a plugin promises a real
upstream repo, a theme promises files that exist — and an entry that cannot
honor its promise is worse than no entry. Drop real ones in; the shelves pick
them up on the next build.

## What is not here, and why

- **No remote registry.** Ryoku's registry, extras actuator and receipt
  transactions stay behind. ASURA's catalogues are the local JSON files in
  `catalog/`, the package manager and the plugin folder, and the local files
  are also the cached content — there is no snapshot of a registry to
  reconcile.
- **No Flatpak.** A Flatpak provider is its own listing, search, remote setup
  and progress parsing, and nothing calls for it on the hosts ASURA targets.
- **No sudo.** The privilege path is polkit: the OS helpers through `pkexec`,
  or `pkexec xbps-*` directly. The store process never holds a privilege.

## Checks

```sh
cd cmd/asura-store && gofmt -l . && go vet ./... && go test -count=1 -race ./...
node --test ui/lib/store.test.mjs
/usr/lib/qt6/bin/qmllint ../../ui/*.qml     # errors fail; the Quickshell/qs.Ui warnings are expected
```

CI runs the same set plus the catalogue-data validation; see
`.github/workflows/ci.yml`.
