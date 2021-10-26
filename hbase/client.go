package hbase

import (
	"context"
	"io"
	"time"

	"github.com/abulo/ratel/v1/logger"
	"github.com/tsuna/gohbase"
	"github.com/tsuna/gohbase/hrpc"
	"github.com/tsuna/gohbase/pb"
)

// ZKConfig Server&Client settings.
type ZKConfig struct {
	Root    string
	Addrs   string
	Timeout time.Duration
}

// Config hbase config
type Config struct {
	Zookeeper           *ZKConfig
	RPCQueueSize        int
	FlushInterval       time.Duration
	EffectiveUser       string
	RegionLookupTimeout time.Duration
	RegionReadTimeout   time.Duration
	TestRowKey          string
	Options             []gohbase.Option
}

// HookFunc hook function call before every method and hook return function will call after finish.
type HookFunc func(ctx context.Context, call hrpc.Call, customName string) func(err error)

// Client hbase client.
type Client struct {
	hc     gohbase.Client
	addr   string
	config *Config
	hooks  []HookFunc
	admin  gohbase.AdminClient
}

//WithOptions ...
func (config *Config) WithOptions(option gohbase.Option) *Config {
	config.Options = append(config.Options, option)
	return config
}

// NewClient new a hbase client.
func (config *Config) NewClient(options ...gohbase.Option) *Client {
	zk := config.Zookeeper.Addrs
	if config.Zookeeper.Root != "" {
		options = append(options, gohbase.ZookeeperRoot(config.Zookeeper.Root))
	}
	if config.Zookeeper.Timeout != 0 {
		options = append(options, gohbase.ZookeeperTimeout(time.Duration(config.Zookeeper.Timeout)))
	}
	if config.RPCQueueSize != 0 {
		logger.Logger.Warn("RPCQueueSize configuration be ignored")
	}
	// force RpcQueueSize = 1, don't change it !!! it has reason
	options = append(options, gohbase.RpcQueueSize(1))

	if config.FlushInterval != 0 {
		options = append(options, gohbase.FlushInterval(time.Duration(config.FlushInterval)))
	}

	if config.EffectiveUser != "" {
		options = append(options, gohbase.EffectiveUser(config.EffectiveUser))
	}
	if config.RegionLookupTimeout != 0 {
		options = append(options, gohbase.RegionLookupTimeout(time.Duration(config.RegionLookupTimeout)))
	}
	if config.RegionReadTimeout != 0 {
		options = append(options, gohbase.RegionReadTimeout(time.Duration(config.RegionReadTimeout)))
	}

	options = append(options, config.Options...)
	// hc := gohbase.NewClient(zk, options...)
	// admin := gohbase.NewAdminClient(zk, options...)
	return &Client{
		hc:     gohbase.NewClient(zk, options...),
		addr:   zk,
		config: config,
		admin:  gohbase.NewAdminClient(zk, options...),
	}
}

// AddHook add hook function.
func (c *Client) AddHook(hookFn HookFunc) {
	c.hooks = append(c.hooks, hookFn)
}

func (c *Client) invokeHook(ctx context.Context, call hrpc.Call, customName string) func(error) {
	finishHooks := make([]func(error), 0, len(c.hooks))
	for _, fn := range c.hooks {
		finishHooks = append(finishHooks, fn(ctx, call, customName))
	}
	return func(err error) {
		for _, fn := range finishHooks {
			fn(err)
		}
	}
}

// ScanAll do scan command and return all result
// NOTE: if err != nil the results is safe for range operate even not result found
func (c *Client) ScanAll(ctx context.Context, table []byte, options ...func(hrpc.Call) error) (results []*hrpc.Result, err error) {
	cursor, err := c.Scan(ctx, table, options...)

	if err != nil {
		return nil, err
	}
	for {
		result, err := cursor.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

type scanTrace struct {
	hrpc.Scanner
	finishHook func(error)
}

func (s *scanTrace) Next() (*hrpc.Result, error) {
	result, err := s.Scanner.Next()
	if err != nil {
		s.finishHook(err)
	}
	return result, err
}

func (s *scanTrace) Close() error {
	err := s.Scanner.Close()
	s.finishHook(err)
	return err
}

// Scan do a scan command.
func (c *Client) Scan(ctx context.Context, table []byte, options ...func(hrpc.Call) error) (scanner hrpc.Scanner, err error) {
	var scan *hrpc.Scan
	scan, err = hrpc.NewScan(ctx, table, options...)
	if err != nil {
		return nil, err
	}
	st := &scanTrace{}
	st.finishHook = c.invokeHook(ctx, scan, "Scan")
	st.Scanner = c.hc.Scan(scan)
	return st, nil
}

// ScanStr scan string
func (c *Client) ScanStr(ctx context.Context, table string, options ...func(hrpc.Call) error) (hrpc.Scanner, error) {
	return c.Scan(ctx, []byte(table), options...)
}

// ScanStrAll scan string
// NOTE: if err != nil the results is safe for range operate even not result found
func (c *Client) ScanStrAll(ctx context.Context, table string, options ...func(hrpc.Call) error) ([]*hrpc.Result, error) {
	return c.ScanAll(ctx, []byte(table), options...)
}

// ScanRange get a scanner for the given table and key range.
// The range is half-open, i.e. [startRow; stopRow[ -- stopRow is not
// included in the range.
func (c *Client) ScanRange(ctx context.Context, table, startRow, stopRow []byte, options ...func(hrpc.Call) error) (scanner hrpc.Scanner, err error) {
	var scan *hrpc.Scan
	scan, err = hrpc.NewScanRange(ctx, table, startRow, stopRow, options...)
	if err != nil {
		return nil, err
	}
	st := &scanTrace{}
	st.finishHook = c.invokeHook(ctx, scan, "ScanRange")
	st.Scanner = c.hc.Scan(scan)
	return st, nil
}

// ScanRangeStr get a scanner for the given table and key range.
// The range is half-open, i.e. [startRow; stopRow[ -- stopRow is not
// included in the range.
func (c *Client) ScanRangeStr(ctx context.Context, table, startRow, stopRow string, options ...func(hrpc.Call) error) (hrpc.Scanner, error) {
	return c.ScanRange(ctx, []byte(table), []byte(startRow), []byte(stopRow), options...)
}

// Get get result for the given table and row key.
// NOTE: if err != nil then result != nil, if result not exists result.Cells length is 0
func (c *Client) Get(ctx context.Context, table, key []byte, options ...func(hrpc.Call) error) (result *hrpc.Result, err error) {
	var get *hrpc.Get
	get, err = hrpc.NewGet(ctx, table, key, options...)
	if err != nil {
		return nil, err
	}

	finishHook := c.invokeHook(ctx, get, "GET")
	result, err = c.hc.Get(get)
	finishHook(err)
	return
}

// GetStr do a get command.
// NOTE: if err != nil then result != nil, if result not exists result.Cells length is 0
func (c *Client) GetStr(ctx context.Context, table, key string, options ...func(hrpc.Call) error) (result *hrpc.Result, err error) {
	return c.Get(ctx, []byte(table), []byte(key), options...)
}

// PutStr insert the given family-column-values in the given row key of the given table.
func (c *Client) PutStr(ctx context.Context, table string, key string, values map[string]map[string][]byte, options ...func(hrpc.Call) error) (*hrpc.Result, error) {
	put, err := hrpc.NewPutStr(ctx, table, key, values, options...)
	if err != nil {
		return nil, err
	}

	finishHook := c.invokeHook(ctx, put, "PUT")
	result, err := c.hc.Put(put)
	finishHook(err)
	return result, err
}

// Delete is used to perform Delete operations on a single row.
// To delete entire row, values should be nil.
//
// To delete specific families, qualifiers map should be nil:
//  map[string]map[string][]byte{
//		"cf1": nil,
//		"cf2": nil,
//  }
//
// To delete specific qualifiers:
//  map[string]map[string][]byte{
//      "cf": map[string][]byte{
//			"q1": nil,
//			"q2": nil,
//		},
//  }
//
// To delete all versions before and at a timestamp, pass hrpc.Timestamp() option.
// By default all versions will be removed.
//
// To delete only a specific version at a timestamp, pass hrpc.DeleteOneVersion() option
// along with a timestamp. For delete specific qualifiers request, if timestamp is not
// passed, only the latest version will be removed. For delete specific families request,
// the timestamp should be passed or it will have no effect as it's an expensive
// operation to perform.
func (c *Client) Delete(ctx context.Context, table string, key string, values map[string]map[string][]byte, options ...func(hrpc.Call) error) (*hrpc.Result, error) {
	del, err := hrpc.NewDelStr(ctx, table, key, values, options...)
	if err != nil {
		return nil, err
	}

	finishHook := c.invokeHook(ctx, del, "Delete")
	result, err := c.hc.Delete(del)
	finishHook(err)
	return result, err
}

// Append do a append command.
func (c *Client) Append(ctx context.Context, table string, key string, values map[string]map[string][]byte, options ...func(hrpc.Call) error) (*hrpc.Result, error) {
	appd, err := hrpc.NewAppStr(ctx, table, key, values, options...)
	if err != nil {
		return nil, err
	}

	finishHook := c.invokeHook(ctx, appd, "Append")
	result, err := c.hc.Append(appd)
	finishHook(err)
	return result, err
}

// Increment the given values in HBase under the given table and key.
func (c *Client) Increment(ctx context.Context, table string, key string, values map[string]map[string][]byte, options ...func(hrpc.Call) error) (int64, error) {
	increment, err := hrpc.NewIncStr(ctx, table, key, values, options...)
	if err != nil {
		return 0, err
	}
	finishHook := c.invokeHook(ctx, increment, "Increment")
	result, err := c.hc.Increment(increment)
	finishHook(err)
	return result, err
}

// IncrementSingle increment the given value by amount in HBase under the given table, key, family and qualifier.
func (c *Client) IncrementSingle(ctx context.Context, table string, key string, family string, qualifier string, amount int64, options ...func(hrpc.Call) error) (int64, error) {
	increment, err := hrpc.NewIncStrSingle(ctx, table, key, family, qualifier, amount, options...)
	if err != nil {
		return 0, err
	}

	finishHook := c.invokeHook(ctx, increment, "IncrementSingle")
	result, err := c.hc.Increment(increment)
	finishHook(err)
	return result, err
}

// Ping ping.
func (c *Client) Ping(ctx context.Context) (err error) {
	testRowKey := "test"
	if c.config.TestRowKey != "" {
		testRowKey = c.config.TestRowKey
	}
	values := map[string]map[string][]byte{"test": map[string][]byte{"test": []byte("test")}}
	_, err = c.PutStr(ctx, "test", testRowKey, values)
	return
}

//CreateTable .
func (c *Client) CreateTable(ctx context.Context, table []byte, families map[string]map[string]string, opts ...func(*hrpc.CreateTable)) error {
	createTable := hrpc.NewCreateTable(ctx, table, families, opts...)
	return c.admin.CreateTable(createTable)
}

//DeleteTable .
func (c *Client) DeleteTable(ctx context.Context, table []byte) error {
	deleteTable := hrpc.NewDeleteTable(ctx, table)
	return c.admin.DeleteTable(deleteTable)
}

// EnableTable .
func (c *Client) EnableTable(ctx context.Context, table []byte) error {
	enableTable := hrpc.NewEnableTable(ctx, table)
	return c.admin.EnableTable(enableTable)
}

// DisableTable .
func (c *Client) DisableTable(ctx context.Context, table []byte) error {
	disableTable := hrpc.NewDisableTable(ctx, table)
	return c.admin.DisableTable(disableTable)
}

//CreateSnapshot .
func (c *Client) CreateSnapshot(ctx context.Context, name string, table string, opts ...func(hrpc.Call) error) error {
	snapshot, err := hrpc.NewSnapshot(ctx, name, table, opts...)
	if err != nil {
		return err
	}
	return c.admin.CreateSnapshot(snapshot)
}

// DeleteSnapshot .
func (c *Client) DeleteSnapshot(ctx context.Context, name string, table string, opts ...func(hrpc.Call) error) error {
	snapshot, err := hrpc.NewSnapshot(ctx, name, table, opts...)
	if err != nil {
		return err
	}
	return c.admin.DeleteSnapshot(snapshot)
}

//ListSnapshots .
func (c *Client) ListSnapshots(ctx context.Context) ([]*pb.SnapshotDescription, error) {
	return c.admin.ListSnapshots(hrpc.NewListSnapshots(ctx))
}

//RestoreSnapshot .
func (c *Client) RestoreSnapshot(ctx context.Context, name string, table string, opts ...func(hrpc.Call) error) error {
	snapshot, err := hrpc.NewSnapshot(ctx, name, table, opts...)
	if err != nil {
		return err
	}
	return c.admin.RestoreSnapshot(snapshot)
}

//ClusterStatus .
func (c *Client) ClusterStatus() (*pb.ClusterStatus, error) {
	return c.admin.ClusterStatus()
}

//ListTableNames .
func (c *Client) ListTableNames(ctx context.Context, opts ...func(hrpc.Call) error) ([]*pb.TableName, error) {
	listTableNames, err := hrpc.NewListTableNames(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return c.admin.ListTableNames(listTableNames)
}

//SetBalancer .
func (c *Client) SetBalancer(ctx context.Context, enabled bool) (bool, error) {
	setBalancer, _ := hrpc.NewSetBalancer(ctx, enabled)
	return c.admin.SetBalancer(setBalancer)
}

//MoveRegion .
func (c *Client) MoveRegion(ctx context.Context, regionName []byte, opts ...func(hrpc.Call) error) error {
	moveRegion, err := hrpc.NewMoveRegion(ctx, regionName, opts...)
	if err != nil {
		return err
	}
	return c.admin.MoveRegion(moveRegion)
}

// Close close client.
func (c *Client) Close() error {
	c.hc.Close()
	return nil
}
