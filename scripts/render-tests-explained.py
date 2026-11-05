#!/usr/bin/env python3
"""Regenerate TESTS_EXPLAINED.md from tests-explained.yaml (source of truth)."""

from __future__ import annotations

import sys
from collections import defaultdict
from pathlib import Path

try:
    import yaml
except ImportError:
    sys.stderr.write("PyYAML required: pip install pyyaml\n")
    sys.exit(1)

ROOT = Path(__file__).resolve().parents[1]
YAML_PATH = ROOT / "tests-explained.yaml"
MD_PATH = ROOT / "TESTS_EXPLAINED.md"


def package_from_id(case_id: str) -> str:
    # id = {import_path}::{TestName}
    left = case_id.split("::", 1)[0]
    module = "github.com/andre-ols/go-face-recognition/"
    if left.startswith(module):
        return left[len(module) :]
    return left


def main() -> int:
    data = yaml.safe_load(YAML_PATH.read_text()) or {}
    cases = data.get("cases") or []

    by_pkg_file: dict[str, dict[str, list]] = defaultdict(lambda: defaultdict(list))
    for c in cases:
        pkg = package_from_id(c["id"])
        by_pkg_file[pkg][c["test_file"]].append(c)

    lines: list[str] = []
    lines.append("# Tests explained")
    lines.append("")
    lines.append(
        "> **Generated file** — source of truth is [`tests-explained.yaml`](tests-explained.yaml). "
        "Do not edit this Markdown by hand; upsert YAML then run "
        "`python3 scripts/render-tests-explained.py`."
    )
    lines.append("")
    lines.append("## How to read this document")
    lines.append("")
    lines.append(
        "Each test case lists **why** it exists (objective), **which skill** "
        "generated or last updated it, and **what production behavior** it protects."
    )
    lines.append("")
    lines.append("## Index by package")
    lines.append("")
    for pkg in sorted(by_pkg_file):
        lines.append(f"- [`{pkg}`](#package-{pkg.replace('/', '')})")
    lines.append("")

    for pkg in sorted(by_pkg_file):
        lines.append(f"## Package `{pkg}`")
        lines.append("")
        for test_file in sorted(by_pkg_file[pkg]):
            fname = Path(test_file).name
            lines.append(f"### File `{fname}`")
            lines.append("")
            lines.append(f"Path: `{test_file}`")
            lines.append("")
            for c in sorted(by_pkg_file[pkg][test_file], key=lambda x: x["test"]):
                lines.append(f"#### `{c['test']}`")
                lines.append("")
                lines.append(f"- **Skill:** `{c['skill']}`")
                obj = " ".join(str(c.get("objective", "")).split())
                lines.append(f"- **Objective:** {obj}")
                lines.append(f"- **System under test:** {c.get('system_under_test', '')}")
                refs = c.get("production_refs") or []
                if refs:
                    lines.append(
                        "- **Production refs:** "
                        + ", ".join(f"`{r}`" for r in refs)
                    )
                if c.get("partitions"):
                    lines.append(
                        "- **Partitions:** "
                        + ", ".join(f"`{p}`" for p in c["partitions"])
                    )
                if c.get("boundaries"):
                    lines.append(
                        "- **Boundaries:** "
                        + ", ".join(f"`{b}`" for b in c["boundaries"])
                    )
                if c.get("covers"):
                    lines.append(f"- **Covers:** {c['covers']}")
                if c.get("kills_mutant"):
                    lines.append(f"- **Kills mutant:** {c['kills_mutant']}")
                if c.get("notes"):
                    lines.append(f"- **Notes:** {c['notes']}")
                lines.append(f"- **Trace:** `id={c['id']}`")
                if c.get("updated_at"):
                    lines.append(f"- **Updated:** {c['updated_at']}")
                lines.append("")

    MD_PATH.write_text("\n".join(lines).rstrip() + "\n")
    print(f"Wrote {MD_PATH.relative_to(ROOT)} ({len(cases)} cases)")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
