# Go 30 Days — Amos Makokha

A focused 30-day journey to become fast and idiomatic with Go by building tiny tools daily and a Real‑Time POS Telemetry capstone.

## Structure
```
go-30days/
  .github/workflows/ci.yml        # CI: lint, test, vuln check
  day01/ ... day30/               # daily folders
  capstone/                        # final project scaffold
  pkg/                             # reusable helpers
  scratch/                         # tiny experiments
  AUTOPLAN.md                      # daily checklist
  STREAK.md                        # streak tracker
  NOTES.md                         # concept notes
  LOG.md                           # daily reflections
  Makefile                         # convenience tasks
  go.mod                           # module
```
## Quick start
```bash
# run Day 1 example
make run DAY=01
# run all tests with race detector
make test
# lint (needs golangci-lint)
make lint
```
Update the module path in `go.mod` after you create the GitHub repo if needed.
