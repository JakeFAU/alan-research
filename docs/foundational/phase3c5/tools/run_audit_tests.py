#!/usr/bin/env python3
"""Run the declared checker tests from their recorded source buffers."""
# This direct-script bootstrap reads the helper once and compiles those bytes.
# Accepted entrypoints then execute their own recorded source buffer as well.
if '__source_bundle__' not in globals():
    import pathlib as _pathlib
    import types as _types
    _root = _pathlib.Path(__file__).resolve().parents[4]
    _helper_path = _root / 'docs/foundational/phase3c5/tools/source_loading.py'
    _helper_bytes = _helper_path.read_bytes()
    _helper = _types.ModuleType('phase3c5_source_loading')
    _helper.__file__ = str(_helper_path)
    exec(compile(_helper_bytes, str(_helper_path), 'exec'), _helper.__dict__)
    __source_bundle__ = _helper.SourceBundle(_root, _helper_bytes)
    if __name__ == '__main__':
        __source_bundle__.load('__main__', str(_pathlib.Path(__file__).resolve().relative_to(_root)))
        raise SystemExit(0)

import json
import sys
import unittest

SOURCE_RUNNER = 'SOURCE'


def main():
    suite = unittest.TestSuite()
    for name in ('test_audit', 'test_source_loading'):
        module = __source_bundle__.load(name, 'docs/foundational/phase3c5/tools/' + name + '.py')
        suite.addTests(unittest.defaultTestLoader.loadTestsFromModule(module))
    print('SOURCE_BUFFERS: ' + json.dumps(__source_bundle__.records, sort_keys=True), flush=True)
    return 0 if unittest.TextTestRunner(verbosity=2).run(suite).wasSuccessful() else 1


if __name__ == '__main__':
    sys.exit(main())
