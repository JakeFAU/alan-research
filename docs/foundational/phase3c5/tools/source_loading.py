"""Source-buffer loading for this finite checker, not a Python sandbox.

Python itself and its standard library are trusted. Only the declared local and
predecessor modules are intercepted. Every intercepted module is compiled from
its recorded byte buffer; no SourceFileLoader or bytecode cache is consulted.
"""
import builtins
import hashlib
import importlib.util
from pathlib import Path
import types

PREDECESSORS = {
    'docs/foundational/phase3c4/tools/audit.py': '163fb14ea65af59a507e52b88c188c0b4fe2df65b40e8d2d5c7aef7f076fc88e',
    'docs/foundational/phase3c3/tools/audit.py': '287be61a4e632b37510d401a658be7fa3fea465697d5586b43f07a58b2b7ad75',
    'docs/foundational/phase3c2/tools/audit.py': '1dc485fb2e3f0594b6122daa3a63ec05a12c8f5a6b9988afe5d2a775a768f193',
}
LOCAL = 'docs/foundational/phase3c5/tools/'
LOCAL_NAMES = ('audit', 'test_audit', 'test_source_loading', 'capture_audit', 'run_audit_tests', 'source_loading')


class BufferLoader:
    def __init__(self, bundle, path):
        self.bundle, self.path = bundle, path

    def create_module(self, spec):
        return None

    def exec_module(self, module):
        module.__file__ = str(self.bundle.root / self.path)
        module.__source_bundle__ = self.bundle
        module.__builtins__ = dict(vars(builtins), __import__=self.bundle.import_module)
        # The same immutable bytes supplied the manifest hash in __init__.
        exec(compile(self.bundle.buffers[self.path], module.__file__, 'exec'), module.__dict__)


class SourceBundle:
    def __init__(self, root, helper_bytes):
        self.root = Path(root).resolve()
        paths = dict(PREDECESSORS, **{LOCAL + n + '.py': None for n in LOCAL_NAMES})
        self.buffers, self.hashes, self.modules = {}, {}, {}
        for path, expected in paths.items():
            raw = helper_bytes if path == LOCAL + 'source_loading.py' else (self.root / path).read_bytes()
            sha = hashlib.sha256(raw).hexdigest()
            if expected is not None and sha != expected:
                raise RuntimeError('checker dependency hash mismatch: ' + path)
            self.buffers[path], self.hashes[path] = raw, sha
        self.records = [dict(Path=p, SHA256=self.hashes[p]) for p in sorted(self.buffers)]
        # Give unchanged predecessors a private importlib facade. Their explicit
        # spec_from_file_location calls cannot reach the cache-capable loader.
        util = types.ModuleType('importlib.util')
        util.__dict__.update(vars(importlib.util))
        util.spec_from_file_location = self.spec_from_file_location
        self.importlib = types.ModuleType('importlib')
        self.importlib.__dict__.update(vars(importlib))
        self.importlib.util = util

    def relative(self, location):
        path = str(Path(location).resolve().relative_to(self.root))
        if path not in self.buffers:
            raise ImportError('undeclared checker source: ' + path)
        return path

    def spec_from_file_location(self, name, location, *args, **kwargs):
        if args or kwargs:
            raise ImportError('unsupported checker loader arguments')
        path = self.relative(location)
        return importlib.util.spec_from_loader(name, BufferLoader(self, path), origin=str(self.root / path))

    def load(self, name, path):
        key = (name, path)
        if key not in self.modules:
            spec = self.spec_from_file_location(name, self.root / path)
            module = importlib.util.module_from_spec(spec)
            self.modules[key] = module
            try:
                spec.loader.exec_module(module)
            except BaseException:
                del self.modules[key]
                raise
        return self.modules[key]

    def import_module(self, name, globals=None, locals=None, fromlist=(), level=0):
        if level == 0 and name in LOCAL_NAMES:
            return self.load(name, LOCAL + name + '.py')
        if level == 0 and name in ('importlib', 'importlib.util'):
            return self.importlib.util if fromlist and name == 'importlib.util' else self.importlib
        return builtins.__import__(name, globals, locals, fromlist, level)
