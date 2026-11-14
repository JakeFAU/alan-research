"""Temporary-copy regressions for the demonstrated executable-source gap."""
import hashlib
import importlib.util
import json
import os
from pathlib import Path
import py_compile
import subprocess
import sys
import tempfile
import unittest

SOURCE_CONTROL = 'SOURCE'
ROOT = Path(__file__).resolve().parents[4]
LOCAL = 'docs/foundational/phase3c5/tools/'


class SourceLoadingTests(unittest.TestCase):
    def copy_sources(self, target):
        for path, raw in __source_bundle__.buffers.items():
            p = target / path
            p.parent.mkdir(parents=True, exist_ok=True)
            p.write_bytes(raw)

    def poison(self, target, relative, old, new):
        path = target / relative
        original = path.read_bytes()
        self.assertIn(old, original)
        self.assertLessEqual(len(new), len(old))
        altered = original.replace(old, new.ljust(len(old)), 1)
        self.assertEqual(len(altered), len(original))
        stat = path.stat()
        path.write_bytes(altered)
        os.utime(path, ns=(stat.st_atime_ns, stat.st_mtime_ns))
        previous_prefix = sys.pycache_prefix
        try:
            sys.pycache_prefix = str(target / 'cache')
            cached = importlib.util.cache_from_source(str(path))
            py_compile.compile(str(path), cfile=cached, doraise=True,
                               invalidation_mode=py_compile.PycInvalidationMode.TIMESTAMP)
        finally:
            sys.pycache_prefix = previous_prefix
        path.write_bytes(original)
        os.utime(path, ns=(stat.st_atime_ns, stat.st_mtime_ns))
        self.assertEqual(path.read_bytes(), original)
        self.assertEqual(path.stat().st_mtime_ns, stat.st_mtime_ns)
        return hashlib.sha256(original).hexdigest(), cached

    def invoke(self, target, code):
        result = subprocess.run([sys.executable, '-B', '-X', 'pycache_prefix=' + str(target/'cache'), '-c', code],
                                cwd=target, capture_output=True, text=True)
        self.assertEqual(result.returncode, 0, result.stdout + result.stderr)
        self.assertEqual(result.stderr, '')
        return result.stdout.strip()

    def test_timestamp_valid_cache_at_every_declared_boundary(self):
        cases = [
            ('inherited-C4', 'docs/foundational/phase3c4/tools/audit.py', b'ALPHABET = prior.ALPHABET',
             b"ALPHABET = ('CACHE',)", 'ALPHABET', "('CACHE',)", 'audit', 'old.ALPHABET'),
            ('inherited-C3', 'docs/foundational/phase3c3/tools/audit.py',
             b"ALPHABET = tuple(k + ':' + n for n in ('AB', 'BC') for k in KINDS)",
             b"ALPHABET = ('CACHE',)", 'ALPHABET', "('CACHE',)", 'audit', 'prior.ALPHABET'),
            ('inherited-C2', 'docs/foundational/phase3c2/tools/audit.py', b"NAMES = ('AB', 'BC')",
             b"NAMES = ('XX', 'BC')", 'NAMES', "('XX', 'BC')", 'audit', 'prior.old.NAMES'),
            ('local-audit', LOCAL+'audit.py', b'ALPHABET = old.ALPHABET',
             b"ALPHABET = ('CACHE',)", 'ALPHABET', "('CACHE',)", 'audit', 'ALPHABET'),
            ('local-test', LOCAL+'test_audit.py',
             b"FIXTURE = audit.ROOT / 'docs/foundational/phase3c5/evidence/task1/capture-3/export.json.gz'",
             b"FIXTURE = 'CACHE'", 'FIXTURE', 'CACHE', 'test_audit', 'FIXTURE'),
            ('local-capture', LOCAL+'capture_audit.py', b'ROOT = audit.ROOT',
             b"ROOT = 'CACHE'", 'ROOT', 'CACHE', 'capture_audit', 'ROOT'),
            ('local-helper', LOCAL+'source_loading.py', b"LOCAL = 'docs/foundational/phase3c5/tools/'",
             b"LOCAL = 'CACHE/'", 'LOCAL', 'CACHE/', 'source_loading', 'LOCAL'),
            ('local-runner', LOCAL+'run_audit_tests.py', b"SOURCE_RUNNER = 'SOURCE'",
             b"SOURCE_RUNNER = 'CACHED'", 'SOURCE_RUNNER', 'CACHED', 'run_audit_tests', 'SOURCE_RUNNER'),
            ('local-regressions', LOCAL+'test_source_loading.py', b"SOURCE_CONTROL = 'SOURCE'",
             b"SOURCE_CONTROL = 'CACHED'", 'SOURCE_CONTROL', 'CACHED', 'test_source_loading', 'SOURCE_CONTROL'),
        ]
        for name, path, old, new, poisoned_attr, poison_value, module, safe_attr in cases:
            with self.subTest(boundary=name), tempfile.TemporaryDirectory(prefix='phase3c5-source-control-') as temp:
                target = Path(temp).resolve()
                self.copy_sources(target)
                sha, cached = self.poison(target, path, old, new)
                # Demonstrate that this precise cache is valid and would execute
                # under the normal source loader, even with -B.
                normal = ("import importlib.util; s=importlib.util.spec_from_file_location('probe', " + repr(path) + "); "
                          "m=importlib.util.module_from_spec(s); s.loader.exec_module(m); print(m." + poisoned_attr + ")")
                observed = self.invoke(target, normal)
                self.assertEqual(observed, poison_value)
                # Execute the actual declared source-only chain, including the
                # unchanged predecessors' nested spec_from_file_location calls.
                safe = ("from pathlib import Path; import types; p=Path('"+LOCAL+"source_loading.py'); "
                        "raw=p.read_bytes(); h=types.ModuleType('helper'); exec(compile(raw,str(p),'exec'),h.__dict__); "
                        "b=h.SourceBundle(Path.cwd(),raw); m=b.load("+repr(module)+","+repr(LOCAL+module+'.py')+"); "
                        "print(m."+safe_attr+")")
                source_value = self.invoke(target, safe)
                if name.startswith('inherited-C') or name == 'local-audit':
                    expected = "('AB', 'BC')" if name == 'inherited-C2' else str(tuple(k+':'+n for n in ('AB','BC') for k in ('commit','reserve-capacity','apply','consume')))
                elif name == 'local-test':
                    expected = str(target/'docs/foundational/phase3c5/evidence/task1/capture-3/export.json.gz')
                elif name == 'local-capture':
                    expected = str(target)
                elif name == 'local-helper':
                    expected = LOCAL
                else:
                    expected = 'SOURCE'
                self.assertEqual(source_value, expected)
                # Capture entrypoint itself must bypass its local stale cache;
                # using an existing directory avoids any qualification rerun.
                if name in ('local-audit','local-test','local-capture','local-helper'):
                    r = subprocess.run([sys.executable,'-B','-X','pycache_prefix='+str(target/'cache'),
                                        LOCAL+'capture_audit.py',str(target)],cwd=target,capture_output=True,text=True)
                    self.assertEqual((r.returncode,r.stderr),(1,'refusing existing capture directory\n'))
                print('CACHE_CONTROL: '+json.dumps(dict(Boundary=name,SourceSHA256=sha,
                      CacheBytes=Path(cached).stat().st_size,NormalLoader=observed,SourceOnly=source_value)),flush=True)

    def test_compiles_recorded_buffer_after_disk_change(self):
        with tempfile.TemporaryDirectory(prefix='phase3c5-buffer-control-') as temp:
            target=Path(temp).resolve()
            self.copy_sources(target)
            code = ("from pathlib import Path; import types,hashlib; p=Path('"+LOCAL+"source_loading.py'); "
                    "raw=p.read_bytes(); h=types.ModuleType('helper'); exec(compile(raw,str(p),'exec'),h.__dict__); "
                    "b=h.SourceBundle(Path.cwd(),raw); p=Path('"+LOCAL+"audit.py'); original=p.read_bytes(); "
                    "p.write_bytes(original.replace(b'ALPHABET = old.ALPHABET',b\"ALPHABET = ('CACHE',)\")); "
                    "m=b.load('audit','"+LOCAL+"audit.py'); assert m.ALPHABET[0]=='commit:AB'; "
                    "assert b.hashes['"+LOCAL+"audit.py']==hashlib.sha256(original).hexdigest(); print('recorded-buffer-executed')")
            self.assertEqual(self.invoke(target,code),'recorded-buffer-executed')

    def test_changed_pinned_source_rejects_before_execution(self):
        with tempfile.TemporaryDirectory(prefix='phase3c5-pin-control-') as temp:
            target=Path(temp).resolve()
            self.copy_sources(target)
            path=target/'docs/foundational/phase3c2/tools/audit.py'
            path.write_bytes(b"raise RuntimeError('unverified-code-executed')\n"+path.read_bytes())
            result=subprocess.run([sys.executable,'-B',LOCAL+'audit.py','irrelevant.json'],cwd=target,capture_output=True,text=True)
            self.assertNotEqual(result.returncode,0)
            self.assertIn('checker dependency hash mismatch:',result.stderr)
            self.assertNotIn('unverified-code-executed',result.stderr)
