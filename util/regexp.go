package util

import "regexp"

// 以下常量为各类文本信息的正则表达式模式字符串，可配合对应预编译正则直接使用。
const (
	// DatePattern 匹配英文或数字形式的日期字符串。
	DatePattern = `(?i)(?:[0-3]?\d(?:st|nd|rd|th)?\s+(?:of\s+)?(?:jan\.?|january|feb\.?|february|mar\.?|march|apr\.?|april|may|jun\.?|june|jul\.?|july|aug\.?|august|sep\.?|september|oct\.?|october|nov\.?|november|dec\.?|december)|(?:jan\.?|january|feb\.?|february|mar\.?|march|apr\.?|april|may|jun\.?|june|jul\.?|july|aug\.?|august|sep\.?|september|oct\.?|october|nov\.?|november|dec\.?|december)\s+[0-3]?\d(?:st|nd|rd|th)?)(?:\,)?\s*(?:\d{4})?|[0-3]?\d[-\./][0-3]?\d[-\./]\d{2,4}`
	// TimePattern 匹配带或不带上下午标记的时间字符串。
	TimePattern = `(?i)\d{1,2}:\d{2} ?(?:[ap]\.?m\.?)?|\d[ap]\.?m\.?`
	// PhonePattern 匹配常见格式的电话号码。
	PhonePattern = `(?:(?:\+?\d{1,3}[-.\s*]?)?(?:\(?\d{3}\)?[-.\s*]?)?\d{3}[-.\s*]?\d{4,6})|(?:(?:(?:\(\+?\d{2}\))|(?:\+?\d{2}))\s*\d{2}\s*\d{3}\s*\d{4})`
	// PhonesWithExtsPattern 匹配带分机号的电话号码。
	PhonesWithExtsPattern = `(?i)(?:(?:\+?1\s*(?:[.-]\s*)?)?(?:\(\s*(?:[2-9]1[02-9]|[2-9][02-8]1|[2-9][02-8][02-9])\s*\)|(?:[2-9]1[02-9]|[2-9][02-8]1|[2-9][02-8][02-9]))\s*(?:[.-]\s*)?)?(?:[2-9]1[02-9]|[2-9][02-9]1|[2-9][02-9]{2})\s*(?:[.-]\s*)?(?:[0-9]{4})(?:\s*(?:#|x\.?|ext\.?|extension)\s*(?:\d+)?)`
	// LinkPattern 匹配网址链接字符串。
	LinkPattern = `(?:(?:https?:\/\/)?(?:[a-z0-9.\-]+|www|[a-z0-9.\-])[.](?:[^\s()<>]+|\((?:[^\s()<>]+|(?:\([^\s()<>]+\)))*\))+(?:\((?:[^\s()<>]+|(?:\([^\s()<>]+\)))*\)|[^\s!()\[\]{};:\'".,<>?]))`
	// EmailPattern 匹配邮箱地址。
	EmailPattern = `(?i)([A-Za-z0-9!#$%&'*+\/=?^_{|.}~-]+@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)`
	// IPv4Pattern 匹配 IPv4 地址。
	IPv4Pattern = `(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`
	// IPv6Pattern 匹配 IPv6 地址。
	IPv6Pattern = `(?:(?:(?:[0-9A-Fa-f]{1,4}:){7}(?:[0-9A-Fa-f]{1,4}|:))|(?:(?:[0-9A-Fa-f]{1,4}:){6}(?::[0-9A-Fa-f]{1,4}|(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){5}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,2})|:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){4}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,3})|(?:(?::[0-9A-Fa-f]{1,4})?:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){3}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,4})|(?:(?::[0-9A-Fa-f]{1,4}){0,2}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){2}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,5})|(?:(?::[0-9A-Fa-f]{1,4}){0,3}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){1}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,6})|(?:(?::[0-9A-Fa-f]{1,4}){0,4}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?::(?:(?:(?::[0-9A-Fa-f]{1,4}){1,7})|(?:(?::[0-9A-Fa-f]{1,4}){0,5}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(?:%.+)?\s*`
	// IPPattern 同时匹配 IPv4 与 IPv6 地址。
	IPPattern = IPv4Pattern + `|` + IPv6Pattern
	// NotKnownPortPattern 匹配未登记为知名服务的端口号。
	NotKnownPortPattern = `6[0-5]{2}[0-3][0-5]|[1-5][\d]{4}|[2-9][\d]{3}|1[1-9][\d]{2}|10[3-9][\d]|102[4-9]`
	// PricePattern 匹配带货币符号与千分位的价格字符串。
	PricePattern = `[$]\s?[+-]?[0-9]{1,3}(?:(?:,?[0-9]{3}))*(?:\.[0-9]{1,2})?`
	// HexColorPattern 匹配三位或六位十六进制颜色值。
	HexColorPattern = `(?:#?([0-9a-fA-F]{6}|[0-9a-fA-F]{3}))`
	// CreditCardPattern 匹配通用信用卡号（15 或 16 位，允许空格或短横线分隔）。
	CreditCardPattern = `(?:(?:(?:\d{4}[- ]?){3}\d{4}|\d{15,16}))`
	// VISACreditCardPattern 匹配以 4 开头的 VISA 卡号。
	VISACreditCardPattern = `4\d{3}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}`
	// MCCreditCardPattern 匹配以 51–55 开头的 MasterCard 卡号。
	MCCreditCardPattern = `5[1-5]\d{2}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}`
	// BtcAddressPattern 匹配比特币地址。
	BtcAddressPattern = `[13][a-km-zA-HJ-NP-Z1-9]{25,34}`
	// StreetAddressPattern 匹配英文街道地址。
	StreetAddressPattern = `\d{1,4} [\w\s]{1,20}(?:street|st|avenue|ave|road|rd|highway|hwy|square|sq|trail|trl|drive|dr|court|ct|park|parkway|pkwy|circle|cir|boulevard|blvd)\W?`
	// ZipCodePattern 匹配美国邮政编码（5 位或 5+4 位）。
	ZipCodePattern = `\b\d{5}(?:[-\s]\d{4})?\b`
	// PoBoxPattern 匹配邮政信箱编号。
	PoBoxPattern = `(?i)P\.? ?O\.? Box \d+`
	// SSNPattern 匹配美国社会安全号码。
	SSNPattern = `(?:\d{3}-\d{2}-\d{4})`
	// MD5HexPattern 匹配 32 位 MD5 十六进制串。
	MD5HexPattern = `[0-9a-fA-F]{32}`
	// SHA1HexPattern 匹配 40 位 SHA1 十六进制串。
	SHA1HexPattern = `[0-9a-fA-F]{40}`
	// SHA256HexPattern 匹配 64 位 SHA256 十六进制串。
	SHA256HexPattern = `[0-9a-fA-F]{64}`
	// GUIDPattern 匹配 GUID/UUID 字符串（短横线可省略）。
	GUIDPattern = `[0-9a-fA-F]{8}-?[a-fA-F0-9]{4}-?[a-fA-F0-9]{4}-?[a-fA-F0-9]{4}-?[a-fA-F0-9]{12}`
	// ISBN13Pattern 匹配 13 位 ISBN 编号。
	ISBN13Pattern = `(?:[\d]-?){12}[\dxX]`
	// ISBN10Pattern 匹配 10 位 ISBN 编号。
	ISBN10Pattern = `(?:[\d]-?){9}[\dxX]`
	// MACAddressPattern 匹配用冒号或短横线分隔的 MAC 地址。
	MACAddressPattern = `(([a-fA-F0-9]{2}[:-]){5}([a-fA-F0-9]{2}))`
	// IBANPattern 匹配国际银行账户号码。
	IBANPattern = `[A-Z]{2}\d{2}[A-Z0-9]{4}\d{7}([A-Z\d]?){0,16}`
	// GitRepoPattern 匹配带协议前缀的 Git 仓库地址。
	GitRepoPattern = `((git|ssh|http(s)?)|(git@[\w\.]+))(:(\/\/)?)([\w\.@\:/\-~]+)(\.git)(\/)?`
)

// 以下变量为上方模式对应的预编译正则表达式，可直接用于匹配，无需再次编译。
var (
	// DateRegex 为 DatePattern 的预编译正则。
	DateRegex = regexp.MustCompile(DatePattern)
	// TimeRegex 为 TimePattern 的预编译正则。
	TimeRegex = regexp.MustCompile(TimePattern)
	// PhoneRegex 为 PhonePattern 的预编译正则。
	PhoneRegex = regexp.MustCompile(PhonePattern)
	// PhonesWithExtsRegex 为 PhonesWithExtsPattern 的预编译正则。
	PhonesWithExtsRegex = regexp.MustCompile(PhonesWithExtsPattern)
	// LinkRegex 为 LinkPattern 的预编译正则。
	LinkRegex = regexp.MustCompile(LinkPattern)
	// EmailRegex 为 EmailPattern 的预编译正则。
	EmailRegex = regexp.MustCompile(EmailPattern)
	// IPv4Regex 为 IPv4Pattern 的预编译正则。
	IPv4Regex = regexp.MustCompile(IPv4Pattern)
	// IPv6Regex 为 IPv6Pattern 的预编译正则。
	IPv6Regex = regexp.MustCompile(IPv6Pattern)
	// IPRegex 为 IPPattern 的预编译正则。
	IPRegex = regexp.MustCompile(IPPattern)
	// NotKnownPortRegex 为 NotKnownPortPattern 的预编译正则。
	NotKnownPortRegex = regexp.MustCompile(NotKnownPortPattern)
	// PriceRegex 为 PricePattern 的预编译正则。
	PriceRegex = regexp.MustCompile(PricePattern)
	// HexColorRegex 为 HexColorPattern 的预编译正则。
	HexColorRegex = regexp.MustCompile(HexColorPattern)
	// CreditCardRegex 为 CreditCardPattern 的预编译正则。
	CreditCardRegex = regexp.MustCompile(CreditCardPattern)
	// BtcAddressRegex 为 BtcAddressPattern 的预编译正则。
	BtcAddressRegex = regexp.MustCompile(BtcAddressPattern)
	// StreetAddressRegex 为 StreetAddressPattern 的预编译正则。
	StreetAddressRegex = regexp.MustCompile(StreetAddressPattern)
	// ZipCodeRegex 为 ZipCodePattern 的预编译正则。
	ZipCodeRegex = regexp.MustCompile(ZipCodePattern)
	// PoBoxRegex 为 PoBoxPattern 的预编译正则。
	PoBoxRegex = regexp.MustCompile(PoBoxPattern)
	// SSNRegex 为 SSNPattern 的预编译正则。
	SSNRegex = regexp.MustCompile(SSNPattern)
	// MD5HexRegex 为 MD5HexPattern 的预编译正则。
	MD5HexRegex = regexp.MustCompile(MD5HexPattern)
	// SHA1HexRegex 为 SHA1HexPattern 的预编译正则。
	SHA1HexRegex = regexp.MustCompile(SHA1HexPattern)
	// SHA256HexRegex 为 SHA256HexPattern 的预编译正则。
	SHA256HexRegex = regexp.MustCompile(SHA256HexPattern)
	// GUIDRegex 为 GUIDPattern 的预编译正则。
	GUIDRegex = regexp.MustCompile(GUIDPattern)
	// ISBN13Regex 为 ISBN13Pattern 的预编译正则。
	ISBN13Regex = regexp.MustCompile(ISBN13Pattern)
	// ISBN10Regex 为 ISBN10Pattern 的预编译正则。
	ISBN10Regex = regexp.MustCompile(ISBN10Pattern)
	// VISACreditCardRegex 为 VISACreditCardPattern 的预编译正则。
	VISACreditCardRegex = regexp.MustCompile(VISACreditCardPattern)
	// MCCreditCardRegex 为 MCCreditCardPattern 的预编译正则。
	MCCreditCardRegex = regexp.MustCompile(MCCreditCardPattern)
	// MACAddressRegex 为 MACAddressPattern 的预编译正则。
	MACAddressRegex = regexp.MustCompile(MACAddressPattern)
	// IBANRegex 为 IBANPattern 的预编译正则。
	IBANRegex = regexp.MustCompile(IBANPattern)
	// GitRepoRegex 为 GitRepoPattern 的预编译正则。
	GitRepoRegex = regexp.MustCompile(GitRepoPattern)
)

func regexpMatch(text string, regex *regexp.Regexp) []string {
	parsed := regex.FindAllString(text, -1)
	return parsed
}

// RegexpDate 从 text 中提取全部日期子串；无匹配时返回空切片。
func RegexpDate(text string) []string {
	return regexpMatch(text, DateRegex)
}

// RegexpTime 从 text 中提取全部时间子串；无匹配时返回空切片。
func RegexpTime(text string) []string {
	return regexpMatch(text, TimeRegex)
}

// RegexpPhones 从 text 中提取全部电话号码；无匹配时返回空切片。
func RegexpPhones(text string) []string {
	return regexpMatch(text, PhoneRegex)
}

// RegexpPhonesWithExts 从 text 中提取全部带分机号的电话号码；无匹配时返回空切片。
func RegexpPhonesWithExts(text string) []string {
	return regexpMatch(text, PhonesWithExtsRegex)
}

// RegexpLinks 从 text 中提取全部链接；无匹配时返回空切片。
func RegexpLinks(text string) []string {
	return regexpMatch(text, LinkRegex)
}

// RegexpEmails 从 text 中提取全部邮箱地址；无匹配时返回空切片。
func RegexpEmails(text string) []string {
	return regexpMatch(text, EmailRegex)
}

// RegexpIPv4s 从 text 中提取全部 IPv4 地址；无匹配时返回空切片。
func RegexpIPv4s(text string) []string {
	return regexpMatch(text, IPv4Regex)
}

// RegexpIPv6s 从 text 中提取全部 IPv6 地址；无匹配时返回空切片。
func RegexpIPv6s(text string) []string {
	return regexpMatch(text, IPv6Regex)
}

// RegexpIPs 从 text 中提取全部 IPv4 与 IPv6 地址；无匹配时返回空切片。
func RegexpIPs(text string) []string {
	return regexpMatch(text, IPRegex)
}

// RegexpNotKnownPorts 从 text 中提取全部非知名端口号；无匹配时返回空切片。
func RegexpNotKnownPorts(text string) []string {
	return regexpMatch(text, NotKnownPortRegex)
}

// RegexpPrices 从 text 中提取全部价格字符串；无匹配时返回空切片。
func RegexpPrices(text string) []string {
	return regexpMatch(text, PriceRegex)
}

// RegexpHexColors 从 text 中提取全部十六进制颜色值；无匹配时返回空切片。
func RegexpHexColors(text string) []string {
	return regexpMatch(text, HexColorRegex)
}

// RegexpCreditCards 从 text 中提取全部信用卡号；无匹配时返回空切片。
func RegexpCreditCards(text string) []string {
	return regexpMatch(text, CreditCardRegex)
}

// RegexpBtcAddresses 从 text 中提取全部比特币地址；无匹配时返回空切片。
func RegexpBtcAddresses(text string) []string {
	return regexpMatch(text, BtcAddressRegex)
}

// RegexpStreetAddresses 从 text 中提取全部街道地址；无匹配时返回空切片。
func RegexpStreetAddresses(text string) []string {
	return regexpMatch(text, StreetAddressRegex)
}

// RegexpZipCodes 从 text 中提取全部邮政编码；无匹配时返回空切片。
func RegexpZipCodes(text string) []string {
	return regexpMatch(text, ZipCodeRegex)
}

// RegexpPoBoxes 从 text 中提取全部邮政信箱编号；无匹配时返回空切片。
func RegexpPoBoxes(text string) []string {
	return regexpMatch(text, PoBoxRegex)
}

// RegexpSSNs 从 text 中提取全部社会安全号码；无匹配时返回空切片。
func RegexpSSNs(text string) []string {
	return regexpMatch(text, SSNRegex)
}

// RegexpMD5Hexes 从 text 中提取全部 MD5 十六进制串；无匹配时返回空切片。
func RegexpMD5Hexes(text string) []string {
	return regexpMatch(text, MD5HexRegex)
}

// RegexpSHA1Hexes 从 text 中提取全部 SHA1 十六进制串；无匹配时返回空切片。
func RegexpSHA1Hexes(text string) []string {
	return regexpMatch(text, SHA1HexRegex)
}

// RegexpSHA256Hexes 从 text 中提取全部 SHA256 十六进制串；无匹配时返回空切片。
func RegexpSHA256Hexes(text string) []string {
	return regexpMatch(text, SHA256HexRegex)
}

// RegexpGUIDs 从 text 中提取全部 GUID/UUID 字符串；无匹配时返回空切片。
func RegexpGUIDs(text string) []string {
	return regexpMatch(text, GUIDRegex)
}

// RegexpISBN13s 从 text 中提取全部 13 位 ISBN 编号；无匹配时返回空切片。
func RegexpISBN13s(text string) []string {
	return regexpMatch(text, ISBN13Regex)
}

// RegexpISBN10s 从 text 中提取全部 10 位 ISBN 编号；无匹配时返回空切片。
func RegexpISBN10s(text string) []string {
	return regexpMatch(text, ISBN10Regex)
}

// RegexpVISACreditCards 从 text 中提取全部 VISA 卡号；无匹配时返回空切片。
func RegexpVISACreditCards(text string) []string {
	return regexpMatch(text, VISACreditCardRegex)
}

// RegexpMCCreditCards 从 text 中提取全部 MasterCard 卡号；无匹配时返回空切片。
func RegexpMCCreditCards(text string) []string {
	return regexpMatch(text, MCCreditCardRegex)
}

// RegexpMACAddresses 从 text 中提取全部 MAC 地址；无匹配时返回空切片。
func RegexpMACAddresses(text string) []string {
	return regexpMatch(text, MACAddressRegex)
}

// RegexpIBANs 从 text 中提取全部 IBAN 号码；无匹配时返回空切片。
func RegexpIBANs(text string) []string {
	return regexpMatch(text, IBANRegex)
}

// RegexpGitRepos 从 text 中提取全部带协议前缀的 Git 仓库地址；无匹配时返回空切片。
func RegexpGitRepos(text string) []string {
	return regexpMatch(text, GitRepoRegex)
}
