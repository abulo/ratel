package util

import (
	"reflect"
	"testing"
)

func TestStrTrim(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrTrim(tt.args.str); got != tt.want {
				t.Errorf("StrTrim() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrPos(t *testing.T) {
	type args struct {
		haystack string
		needle   string
		offset   int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrPos(tt.args.haystack, tt.args.needle, tt.args.offset); got != tt.want {
				t.Errorf("StrPos() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrIPos(t *testing.T) {
	type args struct {
		haystack string
		needle   string
		offset   int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrIPos(tt.args.haystack, tt.args.needle, tt.args.offset); got != tt.want {
				t.Errorf("StrIPos() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrRPos(t *testing.T) {
	type args struct {
		haystack string
		needle   string
		offset   int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrRPos(tt.args.haystack, tt.args.needle, tt.args.offset); got != tt.want {
				t.Errorf("StrRPos() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrRIPos(t *testing.T) {
	type args struct {
		haystack string
		needle   string
		offset   int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrRIPos(tt.args.haystack, tt.args.needle, tt.args.offset); got != tt.want {
				t.Errorf("StrRIPos() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrReplace(t *testing.T) {
	type args struct {
		search  string
		replace string
		subject string
		count   int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrReplace(tt.args.search, tt.args.replace, tt.args.subject, tt.args.count); got != tt.want {
				t.Errorf("StrReplace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrToUpper(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrToUpper(tt.args.str); got != tt.want {
				t.Errorf("StrToUpper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrToLower(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrToLower(tt.args.str); got != tt.want {
				t.Errorf("StrToLower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUCfirst(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UCfirst(tt.args.str); got != tt.want {
				t.Errorf("UCfirst() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLCFirst(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LCFirst(tt.args.str); got != tt.want {
				t.Errorf("LCFirst() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUCWords(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UCWords(tt.args.str); got != tt.want {
				t.Errorf("UCWords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubStr(t *testing.T) {
	type args struct {
		str    string
		start  uint
		length int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SubStr(tt.args.str, tt.args.start, tt.args.length); got != tt.want {
				t.Errorf("SubStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrRev(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrRev(tt.args.str); got != tt.want {
				t.Errorf("StrRev() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseStr(t *testing.T) {
	type args struct {
		encodedString string
		result        map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ParseStr(tt.args.encodedString, tt.args.result); (err != nil) != tt.wantErr {
				t.Errorf("ParseStr() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNumberFormat(t *testing.T) {
	type args struct {
		number       float64
		decimals     uint
		decPoint     string
		thousandsSep string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NumberFormat(tt.args.number, tt.args.decimals, tt.args.decPoint, tt.args.thousandsSep); got != tt.want {
				t.Errorf("NumberFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChunkSplit(t *testing.T) {
	type args struct {
		body     string
		chunklen uint
		end      string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ChunkSplit(tt.args.body, tt.args.chunklen, tt.args.end); got != tt.want {
				t.Errorf("ChunkSplit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrWordCount(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrWordCount(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StrWordCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWordWrap(t *testing.T) {
	type args struct {
		str   string
		width uint
		br    string
		cut   bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WordWrap(tt.args.str, tt.args.width, tt.args.br, tt.args.cut); got != tt.want {
				t.Errorf("WordWrap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrLen(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrLen(tt.args.str); got != tt.want {
				t.Errorf("StrLen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMbStrLen(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MbStrLen(tt.args.str); got != tt.want {
				t.Errorf("MbStrLen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrRepeat(t *testing.T) {
	type args struct {
		input      string
		multiplier int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrRepeat(tt.args.input, tt.args.multiplier); got != tt.want {
				t.Errorf("StrRepeat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrStr(t *testing.T) {
	type args struct {
		haystack string
		needle   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrStr(tt.args.haystack, tt.args.needle); got != tt.want {
				t.Errorf("StrStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrTr(t *testing.T) {
	type args struct {
		haystack string
		params   []interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrTr(tt.args.haystack, tt.args.params...); got != tt.want {
				t.Errorf("StrTr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrShuffle(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrShuffle(tt.args.str); got != tt.want {
				t.Errorf("StrShuffle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrim(t *testing.T) {
	type args struct {
		str           string
		characterMask []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Trim(tt.args.str, tt.args.characterMask...); got != tt.want {
				t.Errorf("Trim() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLTrim(t *testing.T) {
	type args struct {
		str           string
		characterMask []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LTrim(tt.args.str, tt.args.characterMask...); got != tt.want {
				t.Errorf("LTrim() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRTrim(t *testing.T) {
	type args struct {
		str           string
		characterMask []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RTrim(tt.args.str, tt.args.characterMask...); got != tt.want {
				t.Errorf("RTrim() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExplode(t *testing.T) {
	type args struct {
		delimiter string
		str       string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Explode(tt.args.delimiter, tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Explode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChr(t *testing.T) {
	type args struct {
		ascii int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Chr(tt.args.ascii); got != tt.want {
				t.Errorf("Chr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrd(t *testing.T) {
	type args struct {
		char string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Ord(tt.args.char); got != tt.want {
				t.Errorf("Ord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNl2br(t *testing.T) {
	type args struct {
		str     string
		isXhtml bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Nl2br(tt.args.str, tt.args.isXhtml); got != tt.want {
				t.Errorf("Nl2br() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONDecode(t *testing.T) {
	type args struct {
		data []byte
		val  interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := JSONDecode(tt.args.data, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("JSONDecode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJSONEncode(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JSONEncode(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONEncode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSONEncode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONString(t *testing.T) {
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JSONString(tt.args.obj); got != tt.want {
				t.Errorf("JSONString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddSlashes(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddSlashes(tt.args.str); got != tt.want {
				t.Errorf("AddSlashes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStripSlashes(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StripSlashes(tt.args.str); got != tt.want {
				t.Errorf("StripSlashes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuoteMeta(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := QuoteMeta(tt.args.str); got != tt.want {
				t.Errorf("QuoteMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHtmlEntities(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HtmlEntities(tt.args.str); got != tt.want {
				t.Errorf("HtmlEntities() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTMLEntityDecode(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HTMLEntityDecode(tt.args.str); got != tt.want {
				t.Errorf("HTMLEntityDecode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMd5(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Md5(tt.args.str); got != tt.want {
				t.Errorf("Md5() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMd5File(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Md5File(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Md5File() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Md5File() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSha1(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sha1(tt.args.str); got != tt.want {
				t.Errorf("Sha1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSha1File(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sha1File(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sha1File() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Sha1File() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCrc32(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Crc32(tt.args.str); got != tt.want {
				t.Errorf("Crc32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevenshtein(t *testing.T) {
	type args struct {
		str1    string
		str2    string
		costIns int
		costRep int
		costDel int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Levenshtein(tt.args.str1, tt.args.str2, tt.args.costIns, tt.args.costRep, tt.args.costDel); got != tt.want {
				t.Errorf("Levenshtein() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimilarText(t *testing.T) {
	type args struct {
		first   string
		second  string
		percent *float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SimilarText(tt.args.first, tt.args.second, tt.args.percent); got != tt.want {
				t.Errorf("SimilarText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSoundex(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Soundex(tt.args.str); got != tt.want {
				t.Errorf("Soundex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTag(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Tag(tt.args.i); got != tt.want {
				t.Errorf("Tag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringToBytes(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringToBytes(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
