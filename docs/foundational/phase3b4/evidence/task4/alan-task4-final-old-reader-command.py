import pathlib,glob,shutil,subprocess,json,hashlib,os,tempfile
out=pathlib.Path('/tmp/alan-task4-final-old-readers');out.mkdir(exist_ok=True)
binary='/tmp/alan-phase3b4-task1/compatibility/old-reader.test'
sha=lambda p:hashlib.sha256(pathlib.Path(p).read_bytes()).hexdigest()
assert sha(binary)=='f08fc613f6411adb005adb9bf67409bb456eff69cd39599fa9652957e1a1db36'
records=[]
for label,pattern,want in [('13','/tmp/alan-task4-final-normal/K01-*/final/runtime.sqlite',True),('15','/tmp/alan-task4-final-normal/K02-*/final/runtime.sqlite',False),('14','/tmp/alan-task4-final-promoted-normal/K01-*/runtime.sqlite',True),('16','/tmp/alan-task4-final-promoted-normal/K02-*/runtime.sqlite',False)]:
 parents=glob.glob(pattern);assert len(parents)==1,parents
 work=pathlib.Path(tempfile.mkdtemp(prefix='old-format'+label+'-',dir=out));db=work/'copy.sqlite';dest=work/'output';dest.mkdir()
 for suffix in ['','-wal','-shm']:
  if pathlib.Path(parents[0]+suffix).exists():shutil.copyfile(parents[0]+suffix,str(db)+suffix)
 before=sha(db)
 env={k:v for k,v in os.environ.items() if not k.startswith('ALAN_')};env.update(ALAN_PHASE3B3_CHILD='1',ALAN_PHASE2_DB=str(db),ALAN_PHASE3B3_OUTPUT=str(dest))
 argv=[binary,'-test.run=^TestPhase3B3Owner$','-test.v','-test.timeout=2m']
 p=subprocess.run(argv,env=env,cwd='/private/tmp/alan-phase3b3-runtime',stdout=subprocess.PIPE,stderr=subprocess.STDOUT,timeout=130);(work/'command.log').write_bytes(p.stdout)
 result=json.loads((dest/'open.json').read_text());after=sha(db)
 assert (p.returncode==0)==want,(label,p.stdout)
 if not want: assert result['Error']=='unsupported version' and before==after,(label,result)
 records.append(dict(format=label,parent=parents[0],parent_sha256=sha(parents[0]),argv=argv,environment={k:v for k,v in env.items() if k.startswith('ALAN_')},exit=p.returncode,open=result,copy_before_sha256=before,copy_after_sha256=after,expected_accept=want))
(out/'results.json').write_text(json.dumps(dict(binary=binary,binary_sha256=sha(binary),records=records),indent=2)+'\n')
print('exact public13/15 and ancillary promoted14/16 old-reader outcomes verified')
