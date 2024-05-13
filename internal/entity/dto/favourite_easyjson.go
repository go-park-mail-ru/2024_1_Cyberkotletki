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

func easyjsonA4d0e38bDecodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto(in *jlexer.Lexer, out *FavouritesResponse) {
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
		case "favourites":
			if in.IsNull() {
				in.Skip()
				out.Favourites = nil
			} else {
				in.Delim('[')
				if out.Favourites == nil {
					if !in.IsDelim(']') {
						out.Favourites = make([]Favourite, 0, 2)
					} else {
						out.Favourites = []Favourite{}
					}
				} else {
					out.Favourites = (out.Favourites)[:0]
				}
				for !in.IsDelim(']') {
					var v1 Favourite
					(v1).UnmarshalEasyJSON(in)
					out.Favourites = append(out.Favourites, v1)
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
func easyjsonA4d0e38bEncodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto(out *jwriter.Writer, in FavouritesResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"favourites\":"
		out.RawString(prefix[1:])
		if in.Favourites == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Favourites {
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
func (v FavouritesResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonA4d0e38bEncodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v FavouritesResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonA4d0e38bEncodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *FavouritesResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonA4d0e38bDecodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *FavouritesResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonA4d0e38bDecodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto(l, v)
}
func easyjsonA4d0e38bDecodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto1(in *jlexer.Lexer, out *FavouriteStatusResponse) {
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
		case "status":
			out.Status = string(in.String())
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
func easyjsonA4d0e38bEncodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto1(out *jwriter.Writer, in FavouriteStatusResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix[1:])
		out.String(string(in.Status))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v FavouriteStatusResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonA4d0e38bEncodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v FavouriteStatusResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonA4d0e38bEncodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *FavouriteStatusResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonA4d0e38bDecodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *FavouriteStatusResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonA4d0e38bDecodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto1(l, v)
}
func easyjsonA4d0e38bDecodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto2(in *jlexer.Lexer, out *Favourite) {
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
		case "contentID":
			out.ContentID = int(in.Int())
		case "category":
			out.Category = string(in.String())
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
func easyjsonA4d0e38bEncodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto2(out *jwriter.Writer, in Favourite) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"contentID\":"
		out.RawString(prefix[1:])
		out.Int(int(in.ContentID))
	}
	{
		const prefix string = ",\"category\":"
		out.RawString(prefix)
		out.String(string(in.Category))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Favourite) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonA4d0e38bEncodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Favourite) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonA4d0e38bEncodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Favourite) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonA4d0e38bDecodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Favourite) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonA4d0e38bDecodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto2(l, v)
}
func easyjsonA4d0e38bDecodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto3(in *jlexer.Lexer, out *CreateFavouriteRequest) {
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
		case "contentID":
			out.ContentID = int(in.Int())
		case "category":
			out.Category = string(in.String())
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
func easyjsonA4d0e38bEncodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto3(out *jwriter.Writer, in CreateFavouriteRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"contentID\":"
		out.RawString(prefix[1:])
		out.Int(int(in.ContentID))
	}
	{
		const prefix string = ",\"category\":"
		out.RawString(prefix)
		out.String(string(in.Category))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CreateFavouriteRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonA4d0e38bEncodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CreateFavouriteRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonA4d0e38bEncodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CreateFavouriteRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonA4d0e38bDecodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CreateFavouriteRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonA4d0e38bDecodeGithubComGoParkMailRu20241CyberkotletkiInternalEntityDto3(l, v)
}
