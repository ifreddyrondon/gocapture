// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package payload

import (
	json "encoding/json"

	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson6ad23cceDecodeGithubComIfreddyrondonGocapturePayload(in *jlexer.Lexer, out *jsonPayload) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "cap":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.Cap = make(map[string]interface{})
				} else {
					out.Cap = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v1 interface{}
					if m, ok := v1.(easyjson.Unmarshaler); ok {
						m.UnmarshalEasyJSON(in)
					} else if m, ok := v1.(json.Unmarshaler); ok {
						_ = m.UnmarshalJSON(in.Raw())
					} else {
						v1 = in.Interface()
					}
					(out.Cap)[key] = v1
					in.WantComma()
				}
				in.Delim('}')
			}
		case "captures":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.Captures = make(map[string]interface{})
				} else {
					out.Captures = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v2 interface{}
					if m, ok := v2.(easyjson.Unmarshaler); ok {
						m.UnmarshalEasyJSON(in)
					} else if m, ok := v2.(json.Unmarshaler); ok {
						_ = m.UnmarshalJSON(in.Raw())
					} else {
						v2 = in.Interface()
					}
					(out.Captures)[key] = v2
					in.WantComma()
				}
				in.Delim('}')
			}
		case "data":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.Data = make(map[string]interface{})
				} else {
					out.Data = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v3 interface{}
					if m, ok := v3.(easyjson.Unmarshaler); ok {
						m.UnmarshalEasyJSON(in)
					} else if m, ok := v3.(json.Unmarshaler); ok {
						_ = m.UnmarshalJSON(in.Raw())
					} else {
						v3 = in.Interface()
					}
					(out.Data)[key] = v3
					in.WantComma()
				}
				in.Delim('}')
			}
		case "payload":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.Payload = make(map[string]interface{})
				} else {
					out.Payload = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v4 interface{}
					if m, ok := v4.(easyjson.Unmarshaler); ok {
						m.UnmarshalEasyJSON(in)
					} else if m, ok := v4.(json.Unmarshaler); ok {
						_ = m.UnmarshalJSON(in.Raw())
					} else {
						v4 = in.Interface()
					}
					(out.Payload)[key] = v4
					in.WantComma()
				}
				in.Delim('}')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}