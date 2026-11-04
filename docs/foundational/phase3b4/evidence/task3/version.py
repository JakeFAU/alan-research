from pathlib import Path
root=Path('/private/tmp/alan-phase3b3-runtime')
def edit(f,a,b):
 p=root/f;s=p.read_text();assert a in s,(f,a);p.write_text(s.replace(a,b))
f='internal/executor/request_authority.go'
edit(f,'const requestAuthoritySchema = 13','const RequestAuthorityBindingVersion = "alan/request-authority/v2"\nconst requestAuthorityBindingSchema = 15\nconst requestAuthorityBindingConflictSchema = 16\nconst requestAuthoritySchema = 13')
edit(f,'return v == requestAuthoritySchema || v == requestAuthorityConflictSchema','return requestAuthorityVersionForSchema(v) != ""')
edit(f,'func requestAuthorityEnabled(','''func requestAuthorityVersionForSchema(v int) string {
 switch v {
 case requestAuthoritySchema, requestAuthorityConflictSchema: return RequestAuthorityVersion
 case requestAuthorityBindingSchema, requestAuthorityBindingConflictSchema: return RequestAuthorityBindingVersion
 }
 return ""
}
func supportedRequestAuthority(v string) bool { return v == RequestAuthorityVersion || v == RequestAuthorityBindingVersion }
func storedRequestAuthorityVersion(ctx context.Context, q queryer) (string, error) {
 v, err := storedSchemaNumber(ctx,q)
 return requestAuthorityVersionForSchema(v), err
}
func requestAuthorityEnabled(''')
edit(f,'u.Version != RequestAuthorityVersion','!supportedRequestAuthority(u.Version)')
edit(f,'d.Version != RequestAuthorityVersion','!supportedRequestAuthority(d.Version)')
edit(f,'v.Version != RequestAuthorityVersion','!supportedRequestAuthority(v.Version)')
edit(f,'if x.profile.RequestAuthority != RequestAuthorityVersion {','version, err := storedRequestAuthorityVersion(ctx, q)\n if err != nil { return nil, err }\n if x.profile.RequestAuthority != version {')
edit(f,'Version: RequestAuthorityVersion, Namespace: x.namespace','Version: x.profile.RequestAuthority, Namespace: x.namespace')
edit(f,'root := v.Root','version, err := storedRequestAuthorityVersion(ctx, q)\n if err != nil { return LifecycleResult{}, err }\n if (r.Operation == "DelegateAuthority" && d.Version != version) || (r.Operation == "RevokeAuthority" && v.Version != version) { return LifecycleResult{}, ErrVersion }\n root := v.Root')
edit(f,'u := i.Authority','u := i.Authority\n if u.Version != x.profile.RequestAuthority { return nil, ErrVersion }')
edit(f,'if d.Root != u.Root','if d.Version != x.profile.RequestAuthority || grant.Authority.Version != x.profile.RequestAuthority { return nil, ErrVersion }\n if d.Root != u.Root')
edit(f,'func validateRequestAuthorityProgram(p Program, enabled bool) error {','func validateRequestAuthorityProgram(p Program, version string) error {')
edit(f,'if !enabled {\n\t\t\t\treturn ErrVersion','if version == "" || e.Authority.Version != version {\n\t\t\t\treturn ErrVersion')
# Literal claim version is checked without interpreting application strings.
edit(f,'\t\te := in.Effect\n\t\tif e.Authority', '''\t\te := in.Effect
 if version != "" && e.Kind == "lifecycle" && e.Payload["operation"].Literal != nil && authorityOperation(e.Payload["operation"].Literal.Text) && e.Payload["claim"].Literal != nil {
  var tag struct { Version string }
  if json.Unmarshal([]byte(e.Payload["claim"].Literal.Text), &tag) != nil || tag.Version != version { return ErrVersion }
 }
\t\tif e.Authority''')
edit(f,'func validateRequestAuthorityIntent(ctx context.Context, q queryer, i Intent, revision int64) error {\n\tenabled, err := requestAuthorityEnabled(ctx, q)','func validateRequestAuthorityIntent(ctx context.Context, q queryer, i Intent, revision int64) error {\n version, err := storedRequestAuthorityVersion(ctx, q)\n enabled := version != ""')
edit(f,'if !enabled || (i.Kind != "external"','if !enabled || i.Authority.Version != version || (i.Kind != "external"')
edit(f,'_, _, err = parseAuthorityClaim(LifecycleRequest{','d, v, err := parseAuthorityClaim(LifecycleRequest{')
edit(f,'ExpectedRevision: i.Payload["revision"].Number})\n\t\treturn err','ExpectedRevision: i.Payload["revision"].Number})\n if err == nil && ((i.Payload["operation"].Text == "DelegateAuthority" && d.Version != version) || (i.Payload["operation"].Text == "RevokeAuthority" && v.Version != version)) { return ErrVersion }\n\t\treturn err')
edit('internal/executor/request_authority_admission.go','Version: RequestAuthorityVersion, Namespace: x.namespace','Version: x.profile.RequestAuthority, Namespace: x.namespace')
edit('runtime/runtime.go','const RequestAuthorityVersion = executor.RequestAuthorityVersion','const RequestAuthorityVersion = executor.RequestAuthorityVersion\n\n// RequestAuthorityBindingVersion opts into exact historical request selection.\nconst RequestAuthorityBindingVersion = executor.RequestAuthorityBindingVersion')
edit('internal/executor/store.go','composition.RequestAuthority != RequestAuthorityVersion','!supportedRequestAuthority(composition.RequestAuthority)')
edit('internal/executor/store.go','version = requestAuthoritySchema','version = requestAuthoritySchema\n if composition.RequestAuthority == RequestAuthorityBindingVersion {\n version = requestAuthorityBindingSchema\n if _, err = tx.ExecContext(ctx, "CREATE INDEX admissions_observation_id ON admissions(observation_id)"); err != nil { return nil, err }\n }')
edit('internal/executor/store.go','enabled, e := requestAuthorityEnabled','version, e := storedRequestAuthorityVersion')
edit('internal/executor/store.go','validateRequestAuthorityProgram(p, enabled)','validateRequestAuthorityProgram(p, version)')
edit('internal/executor/composition_store.go','c.RequestAuthority != RequestAuthorityVersion','!supportedRequestAuthority(c.RequestAuthority)')
edit('internal/executor/composition_store.go','(c.RequestAuthority != "") != hasRequestAuthoritySchema(stored)','c.RequestAuthority != requestAuthorityVersionForSchema(stored)')
edit('internal/executor/inspection_deployment.go','(profile.RequestAuthority != "") != b.snap.Features.RequestAuthority','profile.RequestAuthority != requestAuthorityVersionForSchema(b.snap.StorageVersion)')
edit('internal/executor/inspection_deployment.go','profile.RequestAuthority != RequestAuthorityVersion','!supportedRequestAuthority(profile.RequestAuthority)')
edit('internal/executor/inspection_execution.go','validateRequestAuthorityProgram(p, b.snap.Features.RequestAuthority)','validateRequestAuthorityProgram(p, requestAuthorityVersionForSchema(b.snap.StorageVersion))')
edit('internal/executor/receipt_conflict.go','continuationConflictSchema, requestAuthorityConflictSchema:','continuationConflictSchema, requestAuthorityConflictSchema, requestAuthorityBindingConflictSchema:')
edit('internal/executor/receipt_conflict.go','case requestAuthoritySchema:\n','case requestAuthorityBindingSchema:\n version = requestAuthorityBindingConflictSchema\n case requestAuthoritySchema:\n')
