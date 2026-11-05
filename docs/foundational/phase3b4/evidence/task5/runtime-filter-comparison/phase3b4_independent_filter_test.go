package executor

import("encoding/json";"errors";"os";"testing")

func TestPhase3B4IndependentEmptyOperationFilter(t *testing.T) {
 for _,test:=range []struct{name,path string;reject bool}{{"genuine",os.Getenv("ALAN_FILTER_PARENT"),false},{"empty-operation",os.Getenv("ALAN_FILTER_ATTACK"),true}} {t.Run(test.name,func(t *testing.T){
  b,err:=os.ReadFile(test.path);if err!=nil {t.Fatal(err)}
  var raw struct{Tables map[string][]map[string]json.RawMessage};if err=json.Unmarshal(b,&raw);err!=nil {t.Fatal(err)}
  text:=func(r map[string]json.RawMessage,k string)string{var s string;if err:=json.Unmarshal(r[k],&s);err!=nil {t.Fatal(err)};return s}
  number:=func(r map[string]json.RawMessage,k string)int64{var n int64;if err:=json.Unmarshal(r[k],&n);err!=nil {t.Fatal(err)};return n}
  var observation Observation;var sequence int64
  for _,r:=range raw.Tables["observations"] {var o Observation;if err=json.Unmarshal([]byte(text(r,"descriptor")),&o);err!=nil {t.Fatal(err)};if o.Run.Subject=="a" && o.Kind=="receipt" {observation=o;sequence=number(r,"consumed_sequence");break}}
  if sequence==0 || observation.Operation=="" {t.Fatal("missing genuine selected nonempty-operation receipt")}
  var prior Projection;pid:=""
  for _,r:=range raw.Tables["history"] {if text(r,"subject")=="a" && number(r,"sequence")==sequence-1 {if err=json.Unmarshal([]byte(text(r,"successor")),&prior);err!=nil {t.Fatal(err)}}}
  for _,r:=range raw.Tables["runs"] {if text(r,"subject")=="a" {pid=text(r,"program_id")}}
  var program Program;for _,r:=range raw.Tables["programs"] {if text(r,"id")==pid {if err=json.Unmarshal([]byte(text(r,"bytes")),&program);err!=nil {t.Fatal(err)}}}
  _,err=Reduce(program,prior,&observation,observation.Run,sequence,DefaultLimits())
  if test.reject {if !errors.Is(err,ErrIneligible) {t.Fatalf("expected exact receive ineligibility, got %v",err)}} else if err!=nil {t.Fatalf("genuine parent reduction: %v",err)}
  t.Logf("sequence=%d operation=%q error=%v",sequence,observation.Operation,err)
 })}
}
