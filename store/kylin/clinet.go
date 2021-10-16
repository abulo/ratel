package kylin

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/abulo/ratel/util"
	"github.com/valyala/fasthttp"
)

//Config 数据库配置
type Config struct {
	Username string //账号 root
	Password string //密码
	Host     string //地址
}

//H 是map[string]interface{} 简写
type H map[string]interface{}

//Client 客户端
type Client struct {
	config  *Config
	request *fasthttp.Request
}

//New 新连接
func New(config *Config) *Client {

	return &Client{
		config: config,
	}
}

//basicAuth 认证
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

//request 获取请求链接
func (client *Client) doRequest() *Client {
	request := fasthttp.AcquireRequest()
	request.Header.Set("Authorization", "Basic "+basicAuth(client.config.Username, client.config.Password))
	client.request = request
	return client
}

//Login 登陆
func (client *Client) Login() (code int, body []byte, err error) {
	client.doRequest()
	client.request.SetRequestURI(client.config.Host + "/kylin/api/user/authentication")
	client.request.Header.SetMethod("POST")
	return client.do()
}

//do 执行语句
func (client *Client) do() (code int, body []byte, err error) {
	response := fasthttp.AcquireResponse()
	err = fasthttp.Do(client.request, response)
	if err == nil {
		code = response.StatusCode()
		body = response.Body()
	}
	return
}

//QueryKylin 查询
func (client *Client) QueryKylin(query *Query) (code int, body []byte, err error) {
	client.doRequest()
	client.request.SetRequestURI(client.config.Host + "/kylin/api/query")
	client.request.Header.SetMethod("POST")
	client.request.Header.Set("Content-Type", "application/json")
	client.request.SetBody(query.GetBytes())
	return client.do()
}

//ListTables 列出tables project 必须参数
func (client *Client) ListTables(project string) (code int, body []byte, err error) {
	if project == "" {
		return 0, nil, fmt.Errorf("project is not empty")
	}
	client.doRequest()
	client.request.SetRequestURI(client.config.Host + "/kylin/api/tables_and_columns")
	client.request.PostArgs().Add("project", project)
	client.request.Header.SetMethod("GET")
	return client.do()
}

//ListCubes 列出cubes offset limit 必须参数 cubeName, projectName 可选参数
func (client *Client) ListCubes(offset, limit int, cubeName, projectName string) (code int, body []byte, err error) {
	client.doRequest()
	client.request.SetRequestURI(client.config.Host + "/kylin/api/cubes")
	client.request.Header.SetMethod("GET")
	client.request.PostArgs().Add("offset", util.ToString(offset))
	client.request.PostArgs().Add("limit", util.ToString(limit))
	client.request.PostArgs().Add("cubeName", cubeName)
	client.request.PostArgs().Add("projectName", projectName)
	return client.do()
}

//GetCube 获取cube //cubeName 必须参数
func (client *Client) GetCube(cubeName string) (code int, body []byte, err error) {
	if cubeName == "" {
		return 0, nil, fmt.Errorf("cubeName is not empty")
	}
	client.doRequest()
	client.request.SetRequestURI(client.config.Host + "/kylin/api/cubes/" + cubeName)
	client.request.Header.SetMethod("GET")
	return client.do()
}

//GetCubeDesc 获取cube描述 //cubeName 必须参数
func (client *Client) GetCubeDesc(cubeName string) (code int, body []byte, err error) {
	if cubeName == "" {
		return 0, nil, fmt.Errorf("cubeName is not empty")
	}
	client.doRequest()
	client.request.SetRequestURI(client.config.Host + "/kylin/api/cube_desc/" + cubeName)
	client.request.Header.SetMethod("GET")
	return client.do()
}

//GetModel 获取model modelName 必须参数
func (client *Client) GetModel(modelName string) (code int, body []byte, err error) {
	if modelName == "" {
		return 0, nil, fmt.Errorf("modelName is not empty")
	}
	client.doRequest()
	client.request.SetRequestURI(client.config.Host + "/kylin/api/model/" + modelName)
	client.request.Header.SetMethod("GET")
	return client.do()
}

//CreateCube 创建cube
func (client *Client) CreateCube(cubeName string) (code int, body []byte, err error) {
	client.doRequest()
	client.request.SetRequestURI(client.config.Host + "/kylin/api/cubes/" + cubeName + "/rebuild")
	client.request.Header.SetMethod("PUT")
	return client.do()
}

//Query 执行语句
func (client *Client) Query(sql, project string, offset, limit int) (*QueryResult, error) {
	query := &Query{
		SQL:     sql,
		Limit:   limit,
		Offset:  offset,
		Project: project,
	}

	code, body, err := client.QueryKylin(query)

	if err != nil {
		return nil, err
	}
	if code != 200 {
		return nil, fmt.Errorf("kylin server return error:" + string(body))
	}

	return handleBody(body)
}

func handleBody(body []byte) (*QueryResult, error) {
	if body == nil {
		return nil, nil
	}
	re := &QueryResult{}
	err := json.Unmarshal(body, re)
	if err != nil {
		return nil, err
	}
	if re.IsException {
		return nil, fmt.Errorf(re.ExceptionMessage)
	}
	return re, nil
}
