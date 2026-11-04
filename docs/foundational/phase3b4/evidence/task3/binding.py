from pathlib import Path
p=Path('/private/tmp/alan-phase3b3-runtime/internal/executor/request_authority.go');s=p.read_text();start=s.index('\t\tif i.Sequence <= 1 {');end=s.index('\n\t}\n\tw.Allowed = true',start);old=s[start:end];s=s[:start]+'''        if x.profile.RequestAuthority == RequestAuthorityBindingVersion {
            selected, consumed, err := selectedOriginalRequest(ctx, q, x, i, d)
            if err != nil { return nil, err }
            if selected == "" { return deny("authority_request_not_selected") }
            w.RequestObservationID, w.RequestConsumedSequence = selected, consumed
        } else {
'''+old+'\n        }'+s[end:];p.write_text(s)
