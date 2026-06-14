"""Unit tests for scripts/validate_plugin.py."""

import json
from pathlib import Path

from validate_plugin import DEFAULT_PATH, main, validate


def write(tmp_path: Path, content: str) -> Path:
    path = tmp_path / "plugin.json"
    path.write_text(content)
    return path


def valid_manifest() -> dict:
    return {"name": "engineering-principles", "version": "1.0.0", "description": "d"}


def test_valid_manifest_passes(tmp_path):
    path = write(tmp_path, json.dumps(valid_manifest()))
    assert validate(path) == []


def test_extra_keys_are_allowed(tmp_path):
    manifest = valid_manifest() | {"author": {"name": "x"}, "license": "MIT"}
    path = write(tmp_path, json.dumps(manifest))
    assert validate(path) == []


def test_missing_file_is_reported(tmp_path):
    errors = validate(tmp_path / "nope.json")
    assert len(errors) == 1
    assert "missing" in errors[0]


def test_invalid_json_is_reported(tmp_path):
    path = write(tmp_path, "{not json")
    errors = validate(path)
    assert len(errors) == 1
    assert "not valid JSON" in errors[0]


def test_non_object_top_level_is_reported(tmp_path):
    path = write(tmp_path, json.dumps(["a", "list"]))
    errors = validate(path)
    assert len(errors) == 1
    assert "JSON object" in errors[0]


def test_missing_required_keys_are_named(tmp_path):
    path = write(tmp_path, json.dumps({"name": "x"}))
    errors = validate(path)
    assert len(errors) == 1
    assert "description" in errors[0]
    assert "version" in errors[0]


def test_empty_string_value_is_reported(tmp_path):
    manifest = valid_manifest() | {"description": "   "}
    path = write(tmp_path, json.dumps(manifest))
    errors = validate(path)
    assert len(errors) == 1
    assert "'description'" in errors[0]


def test_non_string_value_is_reported(tmp_path):
    manifest = valid_manifest() | {"version": 2}
    path = write(tmp_path, json.dumps(manifest))
    errors = validate(path)
    assert len(errors) == 1
    assert "'version'" in errors[0]


def test_main_returns_zero_for_valid_file(tmp_path, capsys):
    path = write(tmp_path, json.dumps(valid_manifest()))
    assert main(["validate_plugin.py", str(path)]) == 0
    assert "ok" in capsys.readouterr().out


def test_main_returns_one_and_prints_errors_to_stderr(tmp_path, capsys):
    path = write(tmp_path, json.dumps({}))
    assert main(["validate_plugin.py", str(path)]) == 1
    assert "missing required keys" in capsys.readouterr().err


def test_repo_manifest_is_valid():
    assert validate(DEFAULT_PATH) == []
