package util

import (
	"reflect"
	"regexp"
	"testing"
)

func Test_regexpMatch(t *testing.T) {
	type args struct {
		text  string
		regex *regexp.Regexp
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
			if got := regexpMatch(tt.args.text, tt.args.regex); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("regexpMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexpDate(t *testing.T) {
	type args struct {
		text string
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
			if got := RegexpDate(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexpDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexpTime(t *testing.T) {
	type args struct {
		text string
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
			if got := RegexpTime(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexpTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexpPhones(t *testing.T) {
	type args struct {
		text string
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
			if got := RegexpPhones(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexpPhones() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexpPhonesWithExts(t *testing.T) {
	type args struct {
		text string
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
			if got := RegexpPhonesWithExts(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexpPhonesWithExts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexpLinks(t *testing.T) {
	type args struct {
		text string
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
			if got := RegexpLinks(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexpLinks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexpEmails(t *testing.T) {
	type args struct {
		text string
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
			if got := RegexpEmails(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexpEmails() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexpIPv4s(t *testing.T) {
	type args struct {
		text string
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
			if got := RegexpIPv4s(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexpIPv4s() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexpIPv6s(t *testing.T) {
	type args struct {
		text string
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
			if got := RegexpIPv6s(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexpIPv6s() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexpIPs(t *testing.T) {
	type args struct {
		text string
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
			if got := RegexpIPs(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexpIPs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexpNotKnownPorts(t *testing.T) {
	type args struct {
		text string
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
			if got := RegexpNotKnownPorts(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexpNotKnownPorts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexpPrices(t *testing.T) {
	type args struct {
		text string
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
			if got := RegexpPrices(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexpPrices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexpHexColors(t *testing.T) {
	type args struct {
		text string
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
			if got := RegexpHexColors(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexpHexColors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexpCreditCards(t *testing.T) {
	type args struct {
		text string
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
			if got := RegexpCreditCards(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexpCreditCards() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexpBtcAddresses(t *testing.T) {
	type args struct {
		text string
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
			if got := RegexpBtcAddresses(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexpBtcAddresses() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexpStreetAddresses(t *testing.T) {
	type args struct {
		text string
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
			if got := RegexpStreetAddresses(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexpStreetAddresses() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexpZipCodes(t *testing.T) {
	type args struct {
		text string
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
			if got := RegexpZipCodes(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexpZipCodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexpPoBoxes(t *testing.T) {
	type args struct {
		text string
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
			if got := RegexpPoBoxes(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexpPoBoxes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexpSSNs(t *testing.T) {
	type args struct {
		text string
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
			if got := RegexpSSNs(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexpSSNs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexpMD5Hexes(t *testing.T) {
	type args struct {
		text string
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
			if got := RegexpMD5Hexes(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexpMD5Hexes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexpSHA1Hexes(t *testing.T) {
	type args struct {
		text string
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
			if got := RegexpSHA1Hexes(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexpSHA1Hexes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexpSHA256Hexes(t *testing.T) {
	type args struct {
		text string
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
			if got := RegexpSHA256Hexes(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexpSHA256Hexes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexpGUIDs(t *testing.T) {
	type args struct {
		text string
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
			if got := RegexpGUIDs(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexpGUIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexpISBN13s(t *testing.T) {
	type args struct {
		text string
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
			if got := RegexpISBN13s(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexpISBN13s() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexpISBN10s(t *testing.T) {
	type args struct {
		text string
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
			if got := RegexpISBN10s(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexpISBN10s() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexpVISACreditCards(t *testing.T) {
	type args struct {
		text string
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
			if got := RegexpVISACreditCards(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexpVISACreditCards() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexpMCCreditCards(t *testing.T) {
	type args struct {
		text string
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
			if got := RegexpMCCreditCards(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexpMCCreditCards() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexpMACAddresses(t *testing.T) {
	type args struct {
		text string
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
			if got := RegexpMACAddresses(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexpMACAddresses() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexpIBANs(t *testing.T) {
	type args struct {
		text string
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
			if got := RegexpIBANs(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexpIBANs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexpGitRepos(t *testing.T) {
	type args struct {
		text string
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
			if got := RegexpGitRepos(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexpGitRepos() = %v, want %v", got, tt.want)
			}
		})
	}
}
