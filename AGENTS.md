# Development workspace

A lifeharness **development** workspace for your own, trusted code.
Agent tools run on the host in this folder with the full toolset (read, write, edit, bash, search).

- Keep throwaway experiments in `scratch/`.
- This is not for untrusted third-party targets — use a web-applications, iot-firmware,
  or macos-research workspace for those; they run in an isolated guest.

## Coding conventions

- Never hardcode environment variables directly in production code. Always load them through a
  dedicated config/env layer using environment-variable handling best practices (centralized
  loading, validation, sane defaults, clear naming, and no secrets committed to source).
- Never re-implement functions already available in the Go standard library. Prefer the stdlib
  implementation over a hand-written version.
