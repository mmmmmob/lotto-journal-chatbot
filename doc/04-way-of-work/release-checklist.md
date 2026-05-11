# Release Checklist

Use this checklist for every release tag.

## 1) Decide version bump

- [ ] Choose `PATCH`, `MINOR`, or `MAJOR` using `versioning-policy.md`.
- [ ] Confirm scope of included tasks is clear.

## 2) Update version metadata

- [ ] Update root `package.json` `version`.
- [ ] Ensure release notes source docs are up to date:
  - `doc/02-task/task-board.md`
  - `doc/03-log/work-log-index.md`

## 3) Validate before tagging

- [ ] `pnpm test:api` passes.
- [ ] `pnpm build` passes.
- [ ] Production CI/CD workflow is green on `main`.

## 4) Create release commit and tag

- [ ] Commit release metadata (example convention: `chore(release): vX.Y.Z`).
- [ ] Create annotated git tag `vX.Y.Z`.

## 5) Post-release checks

- [ ] Verify deploy status on Fly.io (if release deploys automatically).
- [ ] Verify app health endpoint and key bot smoke flow.
- [ ] Record any noteworthy release notes/risks in `doc/03-log/work-log-index.md`.
