words := []string{"林茹", "林如", "临蓐", "空子", "霸王龙", "我是个SB", "是我", "TMD", "他妈的", "他妈"}
	handel := filter.Strings(words)
	str := []byte("我空ss子sss我是霸**王*龙,我是我我是个(S)(B)真的,TMD，他妈的")
	fmt.Println(handel.Find(str))
	fmt.Println(string(handel.Replace(str, '*')))
	fmt.Println(string(handel.ReplaceRune(str, '*')))