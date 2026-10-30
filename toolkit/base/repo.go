package base

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path"
	"strings"
)

var unExpandVarPath = []string{"~", ".", ".."}

// Repo 封装一个 Git 仓库的本地缓存路径、远程地址与分支，供模板拉取与拷贝使用。
//
// Repo 由 NewRepo 创建，内部字段不对外暴露；方法不保证并发安全，同一实例不应被并发调用。
type Repo struct {
	url    string
	home   string
	branch string
}

func repoDir(url string) string {
	vcsURL, err := ParseVCSUrl(url)
	if err != nil {
		return url
	}
	host, _, err := net.SplitHostPort(vcsURL.Host)
	if err != nil {
		host = vcsURL.Host
	}
	for _, p := range unExpandVarPath {
		host = strings.TrimLeft(host, p)
	}
	dir := path.Base(path.Dir(vcsURL.Path))
	url = fmt.Sprintf("%s/%s", host, dir)
	return url
}

// NewRepo 根据远程仓库地址与分支创建 Repo，本地缓存目录由仓库地址推导得出；branch 为空时使用默认分支。
func NewRepo(url string, branch string) *Repo {
	return &Repo{
		url:    url,
		home:   ratelHomeWithDir("repo/" + repoDir(url)),
		branch: branch,
	}
}

// Path 返回仓库在本地缓存目录中的完整路径，分支为空时路径后缀为 @main。
func (r *Repo) Path() string {
	start := strings.LastIndex(r.url, "/")
	end := strings.LastIndex(r.url, ".git")
	if end == -1 {
		end = len(r.url)
	}
	var branch string
	if r.branch == "" {
		branch = "@main"
	} else {
		branch = "@" + r.branch
	}
	return path.Join(r.home, r.url[start+1:end]+branch)
}

// Pull 在本地缓存目录执行 git pull 拉取远程更新。
//
// 目录不是 Git 仓库时直接返回 nil；ctx 用于控制命令超时，git pull 失败时返回对应错误。
func (r *Repo) Pull(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "git", "symbolic-ref", "HEAD")
	cmd.Dir = r.Path()
	_, err := cmd.CombinedOutput()
	if err != nil {
		return nil
	}
	cmd = exec.CommandContext(ctx, "git", "pull")
	cmd.Dir = r.Path()
	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		return err
	}
	return err
}

// Clone 将仓库克隆到本地缓存路径；路径已存在时改为调用 Pull 拉取更新。
//
// 指定分支时通过 git clone -b 检出该分支，克隆失败时返回对应错误。
func (r *Repo) Clone(ctx context.Context) error {
	if _, err := os.Stat(r.Path()); !os.IsNotExist(err) {
		return r.Pull(ctx)
	}
	var cmd *exec.Cmd
	if r.branch == "" {
		cmd = exec.CommandContext(ctx, "git", "clone", r.url, r.Path())
	} else {
		cmd = exec.CommandContext(ctx, "git", "clone", "-b", r.branch, r.url, r.Path())
	}
	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		return err
	}
	return nil
}

// CopyTo 先确保仓库已克隆到本地，再将模板仓库拷贝到项目目录，拷贝时按模块路径与忽略清单排除文件。
func (r *Repo) CopyTo(ctx context.Context, to string, modPath string, ignores []string) error {
	if err := r.Clone(ctx); err != nil {
		return err
	}
	mod, err := ModulePath(path.Join(r.Path(), "go.mod"))
	if err != nil {
		return err
	}

	return copyDir(r.Path(), to, []string{mod, modPath}, ignores)
}

// CopyFrontTo 先确保仓库已克隆到本地，再将前端模板仓库拷贝到项目目录，不做模块路径替换。
func (r *Repo) CopyFrontTo(ctx context.Context, to string, modPath string, ignores []string) error {
	if err := r.Clone(ctx); err != nil {
		return err
	}
	return copyDir(r.Path(), to, []string{}, ignores)
}

// CopyToV2 先确保仓库已克隆到本地，再将模板仓库拷贝到项目目录，替换清单在模块路径替换之外追加。
func (r *Repo) CopyToV2(ctx context.Context, to string, modPath string, ignores, replaces []string) error {
	if err := r.Clone(ctx); err != nil {
		return err
	}
	mod, err := ModulePath(path.Join(r.Path(), "go.mod"))
	if err != nil {
		return err
	}
	replaces = append([]string{mod, modPath}, replaces...)
	return copyDir(r.Path(), to, replaces, ignores)
}
