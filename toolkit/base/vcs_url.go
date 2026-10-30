package base

import (
	"net/url"
	"regexp"
	"slices"
	"strings"

	"github.com/pkg/errors"
)

var (
	scpSyntaxRe = regexp.MustCompile(`^(\w+)@([\w.-]+):(.*)$`)
	scheme      = []string{"git", "https", "http", "git+ssh", "ssh", "file", "ftp", "ftps"}
)

// ParseVCSUrl 将给定的仓库地址解析为标准的 url.URL，支持 SCP 风格与常见 Git 协议地址。
//
// 协议不在支持列表内时返回错误；解析过程不校验地址是否可访问。
//
// 参考：Go 源码实现 https://github.com/golang/go/blob/master/src/cmd/go/internal/vcs/vcs.go
// Git 地址规范 https://git-scm.com/docs/git-clone#_git_urls
func ParseVCSUrl(repo string) (*url.URL, error) {
	var (
		repoURL *url.URL
		err     error
	)

	if m := scpSyntaxRe.FindStringSubmatch(repo); m != nil {
		// 匹配 SCP 风格地址（如 git@github.com:user/repo），转换为 ssh:// 形式的 URL。
		repoURL = &url.URL{
			Scheme: "ssh",
			User:   url.User(m[1]),
			Host:   m[2],
			Path:   m[3],
		}
	} else {
		if !strings.Contains(repo, "//") {
			repo = "//" + repo
		}
		if strings.HasPrefix(repo, "//git@") {
			repo = "ssh:" + repo
		} else if strings.HasPrefix(repo, "//") {
			repo = "https:" + repo
		}
		repoURL, err = url.Parse(repo)
		if err != nil {
			return nil, err
		}
	}

	// 白名单同时包含 http、ftp 等非加密协议：本函数只负责还原地址本身，
	// 不区分协议是否安全，否则无法得到真实的仓库 URL。
	if slices.Contains(scheme, repoURL.Scheme) {
		return repoURL, nil
	}
	return nil, errors.New("无法解析的仓库地址")
}
