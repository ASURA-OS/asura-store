# Security

## Scope

This repository is the standalone Asura Store: the `asura-store` Go data plane
(`cmd/asura-store/`), the Quickshell window (`ui/`), and the catalogue
manifests (`catalog/`). Reports about the code in this tree belong here.

Out of scope: the distribution itself
([`mohitjoer/ASURA-OS`](https://github.com/mohitjoer/ASURA-OS)) — including the
`asura-pkg-add`/`asura-pkg-remove` helpers the store calls but does not ship —
Quickshell, the plugin kit (`asura plugin`), and Ryoku upstream. Anything at
those seams is still welcome if the store's side caused or exposed it.

## The honest threat model

What an attacker could do through or around this process:

- **Privilege escalation is polkit's problem, not this code's.** The store
  never holds a privilege and never runs `sudo`; installs and removals go
  through `pkexec` (the OS helpers on apt hosts, `pkexec xbps-*` on Void
  hosts). Declined prompts surface as sentences, never silent failures.
- **A plugin install executes third-party upstream code.** `install plugins <id>`
  hands the entry's `upstream` repo to `asura plugin install`, which clones and
  installs it. The catalogue in this repository curates those links; treat an
  unexpected `upstream` in any catalogue file — in a PR, in a mirror, or on
  disk — as a supply-chain issue of the highest severity.
- **Ids are input.** Every provider matches a strict id grammar (private
  regexes in the providers) before anything reads a filesystem path or an
  argument, so an id can never escape its directory or read as an option.
- **The store writes exactly two things** in the user's own space: the
  seen-revision stamp and the active-theme pointer (`theme.name`). Themes are
  activation, not download; nothing unpacks or executes a theme payload.
- **Reading is local.** `catalog` and `search` answer from `dpkg`/`apt-cache`
  or `xbps-query`; the store makes no HTTP requests of its own.

## Reporting

**Do not open a public issue for a security problem.** Open the repository's
*Report a vulnerability* form on GitHub (Security → Reports) and describe the
outcome, the exact command sequence against it, and versions in play. Private
reports stay private while they are being fixed.

You can expect an initial reply within a few days. Fixes land on `main` first
(the distribution consumes source builds of this tree), and disclosure is
coordinated with the reporter — credit unless anonymity is requested.
