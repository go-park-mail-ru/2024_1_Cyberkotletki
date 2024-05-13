// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package dto

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

func easyjsonBca8643DecodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto(in *jlexer.Lexer, out *CompilationTypeResponseList) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "compilation_types":
			if in.IsNull() {
				in.Skip()
				out.CompilationTypes = nil
			} else {
				in.Delim('[')
				if out.CompilationTypes == nil {
					if !in.IsDelim(']') {
						out.CompilationTypes = make([]CompilationType, 0, 2)
					} else {
						out.CompilationTypes = []CompilationType{}
					}
				} else {
					out.CompilationTypes = (out.CompilationTypes)[:0]
				}
				for !in.IsDelim(']') {
					var v1 CompilationType
					(v1).UnmarshalEasyJSON(in)
					out.CompilationTypes = append(out.CompilationTypes, v1)
					in.WantComma()
				}
				in.Delim(']')
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
func easyjsonBca8643EncodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto(out *jwriter.Writer, in CompilationTypeResponseList) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"compilation_types\":"
		out.RawString(prefix[1:])
		if in.CompilationTypes == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.CompilationTypes {
				if v2 > 0 {
					out.RawByte(',')
				}
				(v3).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CompilationTypeResponseList) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonBca8643EncodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CompilationTypeResponseList) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonBca8643EncodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CompilationTypeResponseList) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonBca8643DecodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CompilationTypeResponseList) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonBca8643DecodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto(l, v)
}
func easyjsonBca8643DecodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto1(in *jlexer.Lexer, out *CompilationType) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int(in.Int())
		case "type":
			out.Type = string(in.String())
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
func easyjsonBca8643EncodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto1(out *jwriter.Writer, in CompilationType) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int(int(in.ID))
	}
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix)
		out.String(string(in.Type))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CompilationType) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonBca8643EncodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CompilationType) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonBca8643EncodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CompilationType) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonBca8643DecodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CompilationType) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonBca8643DecodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto1(l, v)
}
func easyjsonBca8643DecodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto2(in *jlexer.Lexer, out *CompilationResponseList) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "compilations":
			if in.IsNull() {
				in.Skip()
				out.Compilations = nil
			} else {
				in.Delim('[')
				if out.Compilations == nil {
					if !in.IsDelim(']') {
						out.Compilations = make([]Compilation, 0, 1)
					} else {
						out.Compilations = []Compilation{}
					}
				} else {
					out.Compilations = (out.Compilations)[:0]
				}
				for !in.IsDelim(']') {
					var v4 Compilation
					(v4).UnmarshalEasyJSON(in)
					out.Compilations = append(out.Compilations, v4)
					in.WantComma()
				}
				in.Delim(']')
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
func easyjsonBca8643EncodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto2(out *jwriter.Writer, in CompilationResponseList) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"compilations\":"
		out.RawString(prefix[1:])
		if in.Compilations == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.Compilations {
				if v5 > 0 {
					out.RawByte(',')
				}
				(v6).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CompilationResponseList) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonBca8643EncodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CompilationResponseList) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonBca8643EncodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CompilationResponseList) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonBca8643DecodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CompilationResponseList) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonBca8643DecodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto2(l, v)
}
func easyjsonBca8643DecodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto3(in *jlexer.Lexer, out *CompilationResponse) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "compilation":
			(out.Compilation).UnmarshalEasyJSON(in)
		case "content":
			if in.IsNull() {
				in.Skip()
				out.Content = nil
			} else {
				in.Delim('[')
				if out.Content == nil {
					if !in.IsDelim(']') {
						out.Content = make([]*PreviewContent, 0, 8)
					} else {
						out.Content = []*PreviewContent{}
					}
				} else {
					out.Content = (out.Content)[:0]
				}
				for !in.IsDelim(']') {
					var v7 *PreviewContent
					if in.IsNull() {
						in.Skip()
						v7 = nil
					} else {
						if v7 == nil {
							v7 = new(PreviewContent)
						}
						(*v7).UnmarshalEasyJSON(in)
					}
					out.Content = append(out.Content, v7)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "content_length":
			out.ContentLength = int(in.Int())
		case "page":
			out.Page = int(in.Int())
		case "per_page":
			out.PerPage = int(in.Int())
		case "total_pages":
			out.TotalPages = int(in.Int())
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
func easyjsonBca8643EncodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto3(out *jwriter.Writer, in CompilationResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"compilation\":"
		out.RawString(prefix[1:])
		(in.Compilation).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"content\":"
		out.RawString(prefix)
		if in.Content == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v8, v9 := range in.Content {
				if v8 > 0 {
					out.RawByte(',')
				}
				if v9 == nil {
					out.RawString("null")
				} else {
					(*v9).MarshalEasyJSON(out)
				}
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"content_length\":"
		out.RawString(prefix)
		out.Int(int(in.ContentLength))
	}
	{
		const prefix string = ",\"page\":"
		out.RawString(prefix)
		out.Int(int(in.Page))
	}
	{
		const prefix string = ",\"per_page\":"
		out.RawString(prefix)
		out.Int(int(in.PerPage))
	}
	{
		const prefix string = ",\"total_pages\":"
		out.RawString(prefix)
		out.Int(int(in.TotalPages))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CompilationResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonBca8643EncodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CompilationResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonBca8643EncodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CompilationResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonBca8643DecodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CompilationResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonBca8643DecodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto3(l, v)
}
func easyjsonBca8643DecodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto4(in *jlexer.Lexer, out *Compilation) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int(in.Int())
		case "title":
			out.Title = string(in.String())
		case "compilation_type_id":
			out.CompilationTypeID = int(in.Int())
		case "poster":
			out.PosterURL = string(in.String())
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
func easyjsonBca8643EncodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto4(out *jwriter.Writer, in Compilation) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int(int(in.ID))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"compilation_type_id\":"
		out.RawString(prefix)
		out.Int(int(in.CompilationTypeID))
	}
	{
		const prefix string = ",\"poster\":"
		out.RawString(prefix)
		out.String(string(in.PosterURL))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Compilation) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonBca8643EncodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Compilation) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonBca8643EncodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Compilation) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonBca8643DecodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Compilation) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonBca8643DecodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto4(l, v)
}
