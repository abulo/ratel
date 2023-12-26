package redis

import (
	"context"

	"github.com/spf13/cast"
)

type FTCreateOptions struct {
	OnHash          bool
	OnJSON          bool
	Prefix          []any
	Filter          string
	DefaultLanguage string
	LanguageField   string
	Score           float64
	ScoreField      string
	PayloadField    string
	MaxTextFields   int
	NoOffsets       bool
	Temporary       int
	NoHL            bool
	NoFields        bool
	NoFreqs         bool
	StopWords       []any
	SkipInitalScan  bool
}

type FieldSchema struct {
	FieldName         string
	As                string
	FieldType         SearchFieldType
	Sortable          bool
	UNF               bool
	NoStem            bool
	NoIndex           bool
	PhoneticMatcher   string
	Weight            float64
	Seperator         string
	CaseSensitive     bool
	WithSuffixtrie    bool
	VectorArgs        *FTVectorArgs
	GeoShapeFieldType string
}

type FTVectorArgs struct {
	FlatOptions *FTFlatOptions
	HNSWOptions *FTHNSWOptions
}

type FTFlatOptions struct {
	Type            string
	Dim             int
	DistanceMetric  string
	InitialCapacity int
	BlockSize       int
}

type FTHNSWOptions struct {
	Type                   string
	Dim                    int
	DistanceMetric         string
	InitialCapacity        int
	MaxEdgesPerNode        int
	MaxAllowedEdgesPerNode int
	EFRunTime              int
	Epsilon                float64
}

type FTDropIndexOptions struct {
	DeleteDocs bool
}

type SpellCheckTerms struct {
	Include    bool
	Exclude    bool
	Dictionary string
}

type FTSpellCheckOptions struct {
	Distance int
	Terms    SpellCheckTerms
	Dialect  int
}

type FTExplainOptions struct {
	Dialect string
}

type FTSynUpdateOptions struct {
	SkipInitialScan bool
}

type SearchAggregator int

const (
	SearchInvalid = SearchAggregator(iota)
	SearchAvg
	SearchSum
	SearchMin
	SearchMax
	SearchCount
	SearchCountDistinct
	SearchCountDistinctish
	SearchStdDev
	SearchQuantile
	SearchToList
	SearchFirstValue
	SearchRandomSample
)

func (a SearchAggregator) String() string {
	switch a {
	case SearchInvalid:
		return ""
	case SearchAvg:
		return "AVG"
	case SearchSum:
		return "SUM"
	case SearchMin:
		return "MIN"
	case SearchMax:
		return "MAX"
	case SearchCount:
		return "COUNT"
	case SearchCountDistinct:
		return "COUNT_DISTINCT"
	case SearchCountDistinctish:
		return "COUNT_DISTINCTISH"
	case SearchStdDev:
		return "STDDEV"
	case SearchQuantile:
		return "QUANTILE"
	case SearchToList:
		return "TOLIST"
	case SearchFirstValue:
		return "FIRST_VALUE"
	case SearchRandomSample:
		return "RANDOM_SAMPLE"
	default:
		return ""
	}
}

type SearchFieldType int

const (
	SearchFieldTypeInvalid = SearchFieldType(iota)
	SearchFieldTypeNumeric
	SearchFieldTypeTag
	SearchFieldTypeText
	SearchFieldTypeGeo
	SearchFieldTypeVector
	SearchFieldTypeGeoShape
)

func (t SearchFieldType) String() string {
	switch t {
	case SearchFieldTypeInvalid:
		return ""
	case SearchFieldTypeNumeric:
		return "NUMERIC"
	case SearchFieldTypeTag:
		return "TAG"
	case SearchFieldTypeText:
		return "TEXT"
	case SearchFieldTypeGeo:
		return "GEO"
	case SearchFieldTypeVector:
		return "VECTOR"
	case SearchFieldTypeGeoShape:
		return "GEOSHAPE"
	default:
		return "TEXT"
	}
}

// Each AggregateReducer have different args.
// Please follow https://redis.io/docs/interact/search-and-query/search/aggregations/#supported-groupby-reducers for more information.
type FTAggregateReducer struct {
	Reducer SearchAggregator
	Args    []any
	As      string
}

type FTAggregateGroupBy struct {
	Fields []any
	Reduce []FTAggregateReducer
}

type FTAggregateSortBy struct {
	FieldName string
	Asc       bool
	Desc      bool
}

type FTAggregateApply struct {
	Field string
	As    string
}

type FTAggregateLoad struct {
	Field string
	As    string
}

type FTAggregateWithCursor struct {
	Count   int
	MaxIdle int
}

type FTAggregateOptions struct {
	Verbatim          bool
	LoadAll           bool
	Load              []FTAggregateLoad
	Timeout           int
	GroupBy           []FTAggregateGroupBy
	SortBy            []FTAggregateSortBy
	SortByMax         int
	Apply             []FTAggregateApply
	LimitOffset       int
	Limit             int
	Filter            string
	WithCursor        bool
	WithCursorOptions *FTAggregateWithCursor
	Params            map[string]any
	DialectVersion    int
}

type FTSearchFilter struct {
	FieldName any
	Min       any
	Max       any
}

type FTSearchGeoFilter struct {
	FieldName string
	Longitude float64
	Latitude  float64
	Radius    float64
	Unit      string
}

type FTSearchReturn struct {
	FieldName string
	As        string
}

type FTSearchSortBy struct {
	FieldName string
	Asc       bool
	Desc      bool
}

type FTSearchOptions struct {
	NoContent       bool
	Verbatim        bool
	NoStopWrods     bool
	WithScores      bool
	WithPayloads    bool
	WithSortKeys    bool
	Filters         []FTSearchFilter
	GeoFilter       []FTSearchGeoFilter
	InKeys          []any
	InFields        []any
	Return          []FTSearchReturn
	Slop            int
	Timeout         int
	InOrder         bool
	Language        string
	Expander        string
	Scorer          string
	ExplainScore    bool
	Payload         string
	SortBy          []FTSearchSortBy
	SortByWithCount bool
	LimitOffset     int
	Limit           int
	Params          map[string]any
	DialectVersion  int
}

// FT_List - Lists all the existing indexes in the database.
// For more information, please refer to the Redis documentation:
// [FT._LIST]: (https://redis.io/commands/ft._list/)
func (r *Client) FTList(ctx context.Context) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Do(ctx, "FT._LIST").StringSlice()
		return err
	}, acceptable)
	return
}

// FTAggregate - Performs a search query on an index and applies a series of aggregate transformations to the result.
// The 'index' parameter specifies the index to search, and the 'query' parameter specifies the search query.
// For more information, please refer to the Redis documentation:
// [FT.AGGREGATE]: (https://redis.io/commands/ft.aggregate/)
func (r *Client) FTAggregate(ctx context.Context, index, query string) (val map[string]any, err error) {
	args := []any{"FT.AGGREGATE", index, query}
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		var res any
		res, err = conn.Do(ctx, args...).Result()
		val = cast.ToStringMap(res)
		return err
	}, acceptable)
	return
}

type AggregateQuery []any

func FTAggregateQuery(query string, options *FTAggregateOptions) AggregateQuery {
	queryArgs := []any{query}
	if options != nil {
		if options.Verbatim {
			queryArgs = append(queryArgs, "VERBATIM")
		}
		if options.LoadAll && options.Load != nil {
			panic("FT.AGGREGATE: LOADALL and LOAD are mutually exclusive")
		}
		if options.LoadAll {
			queryArgs = append(queryArgs, "LOAD", "*")
		}
		if options.Load != nil {
			queryArgs = append(queryArgs, "LOAD", len(options.Load))
			for _, load := range options.Load {
				queryArgs = append(queryArgs, load.Field)
				if load.As != "" {
					queryArgs = append(queryArgs, "AS", load.As)
				}
			}
		}
		if options.Timeout > 0 {
			queryArgs = append(queryArgs, "TIMEOUT", options.Timeout)
		}
		if options.GroupBy != nil {
			for _, groupBy := range options.GroupBy {
				queryArgs = append(queryArgs, "GROUPBY", len(groupBy.Fields))
				queryArgs = append(queryArgs, groupBy.Fields...)

				for _, reducer := range groupBy.Reduce {
					queryArgs = append(queryArgs, "REDUCE")
					queryArgs = append(queryArgs, reducer.Reducer.String())
					if reducer.Args != nil {
						queryArgs = append(queryArgs, len(reducer.Args))
						queryArgs = append(queryArgs, reducer.Args...)
					} else {
						queryArgs = append(queryArgs, 0)
					}
					if reducer.As != "" {
						queryArgs = append(queryArgs, "AS", reducer.As)
					}
				}
			}
		}
		if options.SortBy != nil {
			queryArgs = append(queryArgs, "SORTBY")
			sortByOptions := []any{}
			for _, sortBy := range options.SortBy {
				sortByOptions = append(sortByOptions, sortBy.FieldName)
				if sortBy.Asc && sortBy.Desc {
					panic("FT.AGGREGATE: ASC and DESC are mutually exclusive")
				}
				if sortBy.Asc {
					sortByOptions = append(sortByOptions, "ASC")
				}
				if sortBy.Desc {
					sortByOptions = append(sortByOptions, "DESC")
				}
			}
			queryArgs = append(queryArgs, len(sortByOptions))
			queryArgs = append(queryArgs, sortByOptions...)
		}
		if options.SortByMax > 0 {
			queryArgs = append(queryArgs, "MAX", options.SortByMax)
		}
		for _, apply := range options.Apply {
			queryArgs = append(queryArgs, "APPLY", apply.Field)
			if apply.As != "" {
				queryArgs = append(queryArgs, "AS", apply.As)
			}
		}
		if options.LimitOffset > 0 {
			queryArgs = append(queryArgs, "LIMIT", options.LimitOffset)
		}
		if options.Limit > 0 {
			queryArgs = append(queryArgs, options.Limit)
		}
		if options.Filter != "" {
			queryArgs = append(queryArgs, "FILTER", options.Filter)
		}
		if options.WithCursor {
			queryArgs = append(queryArgs, "WITHCURSOR")
			if options.WithCursorOptions != nil {
				if options.WithCursorOptions.Count > 0 {
					queryArgs = append(queryArgs, "COUNT", options.WithCursorOptions.Count)
				}
				if options.WithCursorOptions.MaxIdle > 0 {
					queryArgs = append(queryArgs, "MAXIDLE", options.WithCursorOptions.MaxIdle)
				}
			}
		}
		if options.Params != nil {
			queryArgs = append(queryArgs, "PARAMS", len(options.Params)*2)
			for key, value := range options.Params {
				queryArgs = append(queryArgs, key, value)
			}
		}
		if options.DialectVersion > 0 {
			queryArgs = append(queryArgs, "DIALECT", options.DialectVersion)
		}
	}
	return queryArgs
}

// FTAggregateWithArgs - Performs a search query on an index and applies a series of aggregate transformations to the result.
// The 'index' parameter specifies the index to search, and the 'query' parameter specifies the search query.
// This function also allows for specifying additional options such as: Verbatim, LoadAll, Load, Timeout, GroupBy, SortBy, SortByMax, Apply, LimitOffset, Limit, Filter, WithCursor, Params, and DialectVersion.
// For more information, please refer to the Redis documentation:
// [FT.AGGREGATE]: (https://redis.io/commands/ft.aggregate/)
func (r *Client) FTAggregateWithArgs(ctx context.Context, index string, query string, options *FTAggregateOptions) (val map[string]any, err error) {
	args := []any{"FT.AGGREGATE", index, query}
	if options != nil {
		if options.Verbatim {
			args = append(args, "VERBATIM")
		}
		if options.LoadAll && options.Load != nil {
			panic("FT.AGGREGATE: LOADALL and LOAD are mutually exclusive")
		}
		if options.LoadAll {
			args = append(args, "LOAD", "*")
		}
		if options.Load != nil {
			args = append(args, "LOAD", len(options.Load))
			for _, load := range options.Load {
				args = append(args, load.Field)
				if load.As != "" {
					args = append(args, "AS", load.As)
				}
			}
		}
		if options.Timeout > 0 {
			args = append(args, "TIMEOUT", options.Timeout)
		}
		if options.GroupBy != nil {
			for _, groupBy := range options.GroupBy {
				args = append(args, "GROUPBY", len(groupBy.Fields))
				args = append(args, groupBy.Fields...)

				for _, reducer := range groupBy.Reduce {
					args = append(args, "REDUCE")
					args = append(args, reducer.Reducer.String())
					if reducer.Args != nil {
						args = append(args, len(reducer.Args))
						args = append(args, reducer.Args...)
					} else {
						args = append(args, 0)
					}
					if reducer.As != "" {
						args = append(args, "AS", reducer.As)
					}
				}
			}
		}
		if options.SortBy != nil {
			args = append(args, "SORTBY")
			sortByOptions := []any{}
			for _, sortBy := range options.SortBy {
				sortByOptions = append(sortByOptions, sortBy.FieldName)
				if sortBy.Asc && sortBy.Desc {
					panic("FT.AGGREGATE: ASC and DESC are mutually exclusive")
				}
				if sortBy.Asc {
					sortByOptions = append(sortByOptions, "ASC")
				}
				if sortBy.Desc {
					sortByOptions = append(sortByOptions, "DESC")
				}
			}
			args = append(args, len(sortByOptions))
			args = append(args, sortByOptions...)
		}
		if options.SortByMax > 0 {
			args = append(args, "MAX", options.SortByMax)
		}
		for _, apply := range options.Apply {
			args = append(args, "APPLY", apply.Field)
			if apply.As != "" {
				args = append(args, "AS", apply.As)
			}
		}
		if options.LimitOffset > 0 {
			args = append(args, "LIMIT", options.LimitOffset)
		}
		if options.Limit > 0 {
			args = append(args, options.Limit)
		}
		if options.Filter != "" {
			args = append(args, "FILTER", options.Filter)
		}
		if options.WithCursor {
			args = append(args, "WITHCURSOR")
			if options.WithCursorOptions != nil {
				if options.WithCursorOptions.Count > 0 {
					args = append(args, "COUNT", options.WithCursorOptions.Count)
				}
				if options.WithCursorOptions.MaxIdle > 0 {
					args = append(args, "MAXIDLE", options.WithCursorOptions.MaxIdle)
				}
			}
		}
		if options.Params != nil {
			args = append(args, "PARAMS", len(options.Params)*2)
			for key, value := range options.Params {
				args = append(args, key, value)
			}
		}
		if options.DialectVersion > 0 {
			args = append(args, "DIALECT", options.DialectVersion)
		}
	}

	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		var res any
		res, err = conn.Do(ctx, args...).Result()
		val = cast.ToStringMap(res)
		return err
	}, acceptable)
	return
}

// FTAliasAdd - Adds an alias to an index.
// The 'index' parameter specifies the index to which the alias is added, and the 'alias' parameter specifies the alias.
// For more information, please refer to the Redis documentation:
// [FT.ALIASADD]: (https://redis.io/commands/ft.aliasadd/)
func (r *Client) FTAliasAdd(ctx context.Context, index, alias string) (val string, err error) {
	args := []any{"FT.ALIASADD", alias, index}
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Do(ctx, args...).Text()
		return err
	}, acceptable)
	return
}

// FTAliasDel - Removes an alias from an index.
// The 'alias' parameter specifies the alias to be removed.
// For more information, please refer to the Redis documentation:
// [FT.ALIASDEL]: (https://redis.io/commands/ft.aliasdel/)
func (r *Client) FTAliasDel(ctx context.Context, alias string) (val string, err error) {
	args := []any{"FT.ALIASDEL", alias}
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Do(ctx, args...).Text()
		return err
	}, acceptable)
	return
}

// FTAliasUpdate - Updates an alias to an index.
// The 'index' parameter specifies the index to which the alias is updated, and the 'alias' parameter specifies the alias.
// If the alias already exists for a different index, it updates the alias to point to the specified index instead.
// For more information, please refer to the Redis documentation:
// [FT.ALIASUPDATE]: (https://redis.io/commands/ft.aliasupdate/)
func (r *Client) FTAliasUpdate(ctx context.Context, index, alias string) (val string, err error) {
	args := []any{"FT.ALIASUPDATE", alias, index}
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Do(ctx, args...).Text()
		return err
	}, acceptable)
	return
}

// FTAlter - Alters the definition of an existing index.
// The 'index' parameter specifies the index to alter, and the 'skipInitalScan' parameter specifies whether to skip the initial scan.
// The 'definition' parameter specifies the new definition for the index.
// For more information, please refer to the Redis documentation:
// [FT.ALTER]: (https://redis.io/commands/ft.alter/)
func (r *Client) FTAlter(ctx context.Context, index string, skipInitalScan bool, definition []any) (val string, err error) {
	args := []any{"FT.ALTER", index}
	if skipInitalScan {
		args = append(args, "SKIPINITIALSCAN")
	}
	args = append(args, "SCHEMA", "ADD")
	args = append(args, definition...)
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Do(ctx, args...).Text()
		return err
	}, acceptable)
	return
}

// FTConfigGet - Retrieves the value of a RediSearch configuration parameter.
// The 'option' parameter specifies the configuration parameter to retrieve.
// For more information, please refer to the Redis documentation:
// [FT.CONFIG GET]: (https://redis.io/commands/ft.config-get/)
func (r *Client) FTConfigGet(ctx context.Context, option string) (val map[string]any, err error) {
	args := []any{"FT.CONFIG", "GET", option}
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		var res any
		res, err = conn.Do(ctx, args...).Result()
		val = cast.ToStringMap(res)
		return err
	}, acceptable)
	return
}

// FTConfigSet - Sets the value of a RediSearch configuration parameter.
// The 'option' parameter specifies the configuration parameter to set, and the 'value' parameter specifies the new value.
// For more information, please refer to the Redis documentation:
// [FT.CONFIG SET]: (https://redis.io/commands/ft.config-set/)
func (r *Client) FTConfigSet(ctx context.Context, option string, value any) (val string, err error) {
	args := []any{"FT.CONFIG", "SET", option, value}
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Do(ctx, args...).Text()
		return err
	}, acceptable)
	return
}

// FTCreate - Creates a new index with the given options and schema.
// The 'index' parameter specifies the name of the index to create.
// The 'options' parameter specifies various options for the index, such as:
// whether to index hashes or JSONs, prefixes, filters, default language, score, score field, payload field, etc.
// The 'schema' parameter specifies the schema for the index, which includes the field name, field type, etc.
// For more information, please refer to the Redis documentation:
// [FT.CREATE]: (https://redis.io/commands/ft.create/)
func (r *Client) FTCreate(ctx context.Context, index string, options *FTCreateOptions, schema ...*FieldSchema) (val string, err error) {
	args := []any{"FT.CREATE", index}
	if options != nil {
		if options.OnHash && !options.OnJSON {
			args = append(args, "ON", "HASH")
		}
		if options.OnJSON && !options.OnHash {
			args = append(args, "ON", "JSON")
		}
		if options.OnHash && options.OnJSON {
			panic("FT.CREATE: ON HASH and ON JSON are mutually exclusive")
		}
		if options.Prefix != nil {
			args = append(args, "PREFIX", len(options.Prefix))
			args = append(args, options.Prefix...)
		}
		if options.Filter != "" {
			args = append(args, "FILTER", options.Filter)
		}
		if options.DefaultLanguage != "" {
			args = append(args, "LANGUAGE", options.DefaultLanguage)
		}
		if options.LanguageField != "" {
			args = append(args, "LANGUAGE_FIELD", options.LanguageField)
		}
		if options.Score > 0 {
			args = append(args, "SCORE", options.Score)
		}
		if options.ScoreField != "" {
			args = append(args, "SCORE_FIELD", options.ScoreField)
		}
		if options.PayloadField != "" {
			args = append(args, "PAYLOAD_FIELD", options.PayloadField)
		}
		if options.MaxTextFields > 0 {
			args = append(args, "MAXTEXTFIELDS", options.MaxTextFields)
		}
		if options.NoOffsets {
			args = append(args, "NOOFFSETS")
		}
		if options.Temporary > 0 {
			args = append(args, "TEMPORARY", options.Temporary)
		}
		if options.NoHL {
			args = append(args, "NOHL")
		}
		if options.NoFields {
			args = append(args, "NOFIELDS")
		}
		if options.NoFreqs {
			args = append(args, "NOFREQS")
		}
		if options.StopWords != nil {
			args = append(args, "STOPWORDS", len(options.StopWords))
			args = append(args, options.StopWords...)
		}
		if options.SkipInitalScan {
			args = append(args, "SKIPINITIALSCAN")
		}
	}
	if schema == nil {
		panic("FT.CREATE: SCHEMA is required")
	}
	args = append(args, "SCHEMA")
	for _, schema := range schema {
		if schema.FieldName == "" || schema.FieldType == SearchFieldTypeInvalid {
			panic("FT.CREATE: SCHEMA FieldName and FieldType are required")
		}
		args = append(args, schema.FieldName)
		if schema.As != "" {
			args = append(args, "AS", schema.As)
		}
		args = append(args, schema.FieldType.String())
		if schema.VectorArgs != nil {
			if schema.FieldType != SearchFieldTypeVector {
				panic("FT.CREATE: SCHEMA FieldType VECTOR is required for VectorArgs")
			}
			if schema.VectorArgs.FlatOptions != nil && schema.VectorArgs.HNSWOptions != nil {
				panic("FT.CREATE: SCHEMA VectorArgs FlatOptions and HNSWOptions are mutually exclusive")
			}
			if schema.VectorArgs.FlatOptions != nil {
				args = append(args, "FLAT")
				if schema.VectorArgs.FlatOptions.Type == "" || schema.VectorArgs.FlatOptions.Dim == 0 || schema.VectorArgs.FlatOptions.DistanceMetric == "" {
					panic("FT.CREATE: Type, Dim and DistanceMetric are required for VECTOR FLAT")
				}
				flatArgs := []any{
					"TYPE", schema.VectorArgs.FlatOptions.Type,
					"DIM", schema.VectorArgs.FlatOptions.Dim,
					"DISTANCE_METRIC", schema.VectorArgs.FlatOptions.DistanceMetric,
				}
				if schema.VectorArgs.FlatOptions.InitialCapacity > 0 {
					flatArgs = append(flatArgs, "INITIAL_CAP", schema.VectorArgs.FlatOptions.InitialCapacity)
				}
				if schema.VectorArgs.FlatOptions.BlockSize > 0 {
					flatArgs = append(flatArgs, "BLOCK_SIZE", schema.VectorArgs.FlatOptions.BlockSize)
				}
				args = append(args, len(flatArgs))
				args = append(args, flatArgs...)
			}
			if schema.VectorArgs.HNSWOptions != nil {
				args = append(args, "HNSW")
				if schema.VectorArgs.HNSWOptions.Type == "" || schema.VectorArgs.HNSWOptions.Dim == 0 || schema.VectorArgs.HNSWOptions.DistanceMetric == "" {
					panic("FT.CREATE: Type, Dim and DistanceMetric are required for VECTOR HNSW")
				}
				hnswArgs := []any{
					"TYPE", schema.VectorArgs.HNSWOptions.Type,
					"DIM", schema.VectorArgs.HNSWOptions.Dim,
					"DISTANCE_METRIC", schema.VectorArgs.HNSWOptions.DistanceMetric,
				}
				if schema.VectorArgs.HNSWOptions.InitialCapacity > 0 {
					hnswArgs = append(hnswArgs, "INITIAL_CAP", schema.VectorArgs.HNSWOptions.InitialCapacity)
				}
				if schema.VectorArgs.HNSWOptions.MaxEdgesPerNode > 0 {
					hnswArgs = append(hnswArgs, "M", schema.VectorArgs.HNSWOptions.MaxEdgesPerNode)
				}
				if schema.VectorArgs.HNSWOptions.MaxAllowedEdgesPerNode > 0 {
					hnswArgs = append(hnswArgs, "EF_CONSTRUCTION", schema.VectorArgs.HNSWOptions.MaxAllowedEdgesPerNode)
				}
				if schema.VectorArgs.HNSWOptions.EFRunTime > 0 {
					hnswArgs = append(hnswArgs, "EF_RUNTIME", schema.VectorArgs.HNSWOptions.EFRunTime)
				}
				if schema.VectorArgs.HNSWOptions.Epsilon > 0 {
					hnswArgs = append(hnswArgs, "EPSILON", schema.VectorArgs.HNSWOptions.Epsilon)
				}
				args = append(args, len(hnswArgs))
				args = append(args, hnswArgs...)
			}
		}
		if schema.GeoShapeFieldType != "" {
			if schema.FieldType != SearchFieldTypeGeoShape {
				panic("FT.CREATE: SCHEMA FieldType GEOSHAPE is required for GeoShapeFieldType")
			}
			args = append(args, schema.GeoShapeFieldType)
		}
		if schema.NoStem {
			args = append(args, "NOSTEM")
		}
		if schema.Sortable {
			args = append(args, "SORTABLE")
		}
		if schema.UNF {
			args = append(args, "UNF")
		}
		if schema.NoIndex {
			args = append(args, "NOINDEX")
		}
		if schema.PhoneticMatcher != "" {
			args = append(args, "PHONETIC", schema.PhoneticMatcher)
		}
		if schema.Weight > 0 {
			args = append(args, "WEIGHT", schema.Weight)
		}
		if schema.Seperator != "" {
			args = append(args, "SEPERATOR", schema.Seperator)
		}
		if schema.CaseSensitive {
			args = append(args, "CASESENSITIVE")
		}
		if schema.WithSuffixtrie {
			args = append(args, "WITHSUFFIXTRIE")
		}
	}
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Do(ctx, args...).Text()
		return err
	}, acceptable)
	return
}

// FTCursorDel - Deletes a cursor from an existing index.
// The 'index' parameter specifies the index from which to delete the cursor, and the 'cursorId' parameter specifies the ID of the cursor to delete.
// For more information, please refer to the Redis documentation:
// [FT.CURSOR DEL]: (https://redis.io/commands/ft.cursor-del/)
func (r *Client) FTCursorDel(ctx context.Context, index string, cursorId int) (val string, err error) {
	args := []any{"FT.CURSOR", "DEL", index, cursorId}
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Do(ctx, args...).Text()
		return err
	}, acceptable)
	return
}

// FTCursorRead - Reads the next results from an existing cursor.
// The 'index' parameter specifies the index from which to read the cursor, the 'cursorId' parameter specifies the ID of the cursor to read, and the 'count' parameter specifies the number of results to read.
// For more information, please refer to the Redis documentation:
// [FT.CURSOR READ]: (https://redis.io/commands/ft.cursor-read/)
func (r *Client) FTCursorRead(ctx context.Context, index string, cursorId int, count int) (val map[string]any, err error) {
	args := []any{"FT.CURSOR", "READ", index, cursorId}
	if count > 0 {
		args = append(args, "COUNT", count)
	}
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		var res any
		res, err = conn.Do(ctx, args...).Result()
		val = cast.ToStringMap(res)
		return err
	}, acceptable)
	return
}

// FTDictAdd - Adds terms to a dictionary.
// The 'dict' parameter specifies the dictionary to which to add the terms, and the 'term' parameter specifies the terms to add.
// For more information, please refer to the Redis documentation:
// [FT.DICTADD]: (https://redis.io/commands/ft.dictadd/)
func (r *Client) FTDictAdd(ctx context.Context, dict string, term []any) (val int64, err error) {
	args := []any{"FT.DICTADD", dict}
	args = append(args, term...)
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Do(ctx, args...).Int64()
		return err
	}, acceptable)
	return
}

// FTDictDel - Deletes terms from a dictionary.
// The 'dict' parameter specifies the dictionary from which to delete the terms, and the 'term' parameter specifies the terms to delete.
// For more information, please refer to the Redis documentation:
// [FT.DICTDEL]: (https://redis.io/commands/ft.dictdel/)
func (r *Client) FTDictDel(ctx context.Context, dict string, term []any) (val int64, err error) {
	args := []any{"FT.DICTDEL", dict}
	args = append(args, term...)
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Do(ctx, args...).Int64()
		return err
	}, acceptable)
	return
}

// FTDictDump - Returns all terms in the specified dictionary.
// The 'dict' parameter specifies the dictionary from which to return the terms.
// For more information, please refer to the Redis documentation:
// [FT.DICTDUMP]: (https://redis.io/commands/ft.dictdump/)
func (r *Client) FTDictDump(ctx context.Context, dict string) (val map[string]any, err error) {
	args := []any{"FT.DICTDUMP", dict}
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		var res any
		res, err = conn.Do(ctx, args...).Result()
		val = cast.ToStringMap(res)
		return err
	}, acceptable)
	return
}

// FTDropIndex - Deletes an index.
// The 'index' parameter specifies the index to delete.
// For more information, please refer to the Redis documentation:
// [FT.DROPINDEX]: (https://redis.io/commands/ft.dropindex/)
func (r *Client) FTDropIndex(ctx context.Context, index string) (val string, err error) {
	args := []any{"FT.DROPINDEX", index}
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Do(ctx, args...).Text()
		return err
	}, acceptable)
	return
}

// FTDropIndexWithArgs - Deletes an index with options.
// The 'index' parameter specifies the index to delete, and the 'options' parameter specifies the DeleteDocs option for docs deletion.
// For more information, please refer to the Redis documentation:
// [FT.DROPINDEX]: (https://redis.io/commands/ft.dropindex/)
func (r *Client) FTDropIndexWithArgs(ctx context.Context, index string, options *FTDropIndexOptions) (val string, err error) {
	args := []any{"FT.DROPINDEX", index}
	if options != nil {
		if options.DeleteDocs {
			args = append(args, "DD")
		}
	}
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Do(ctx, args...).Text()
		return err
	}, acceptable)
	return
}

// FTExplain - Returns the execution plan for a complex query.
// The 'index' parameter specifies the index to query, and the 'query' parameter specifies the query string.
// For more information, please refer to the Redis documentation:
// [FT.EXPLAIN]: (https://redis.io/commands/ft.explain/)
func (r *Client) FTExplain(ctx context.Context, index string, query string) (val int64, err error) {
	args := []any{"FT.EXPLAIN", index, query}
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Do(ctx, args...).Int64()
		return err
	}, acceptable)
	return
}

// FTExplainWithArgs - Returns the execution plan for a complex query with options.
// The 'index' parameter specifies the index to query, the 'query' parameter specifies the query string, and the 'options' parameter specifies the Dialect for the query.
// For more information, please refer to the Redis documentation:
// [FT.EXPLAIN]: (https://redis.io/commands/ft.explain/)
func (r *Client) FTExplainWithArgs(ctx context.Context, index string, query string, options *FTExplainOptions) (val int64, err error) {
	args := []any{"FT.EXPLAIN", index, query}
	if options.Dialect != "" {
		args = append(args, "DIALECT", options.Dialect)
	}
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Do(ctx, args...).Int64()
		return err
	}, acceptable)
	return
}

// FTExplainCli - Returns the execution plan for a complex query. [Not Implemented]
// For more information, see https://redis.io/commands/ft.explaincli/
func (r *Client) FTExplainCli(ctx context.Context, key, path string) error {
	panic("not implemented")
}

// FTInfo - Retrieves information about an index.
// The 'index' parameter specifies the index to retrieve information about.
// For more information, please refer to the Redis documentation:
// [FT.INFO]: (https://redis.io/commands/ft.info/)
func (r *Client) FTInfo(ctx context.Context, index string) (val map[string]any, err error) {
	args := []any{"FT.INFO", index}
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		var res any
		res, err = conn.Do(ctx, args...).Result()
		val = cast.ToStringMap(res)
		return err
	}, acceptable)
	return
}

// FTProfileSearch - Executes a search query and returns a profile of how the query was processed.
// The 'index' parameter specifies the index to search, the 'limited' parameter specifies whether to limit the results,
// and the 'query' parameter specifies the search / aggreagte query. Please notice that you must either pass a SearchQuery or an AggregateQuery.
// For more information, please refer to the Redis documentation:
// [FT.PROFILE SEARCH]: (https://redis.io/commands/ft.profile/)
func (r *Client) FTProfile(ctx context.Context, index string, limited bool, query any) (val map[string]any, err error) {
	queryType := ""
	var argsQuery []any

	switch v := query.(type) {
	case AggregateQuery:
		queryType = "AGGREGATE"
		argsQuery = v
	case SearchQuery:
		queryType = "SEARCH"
		argsQuery = v
	default:
		panic("FT.PROFILE: query must be either AggregateQuery or SearchQuery")
	}

	args := []any{"FT.PROFILE", index, queryType}

	if limited {
		args = append(args, "LIMITED")
	}
	args = append(args, "QUERY")
	args = append(args, argsQuery...)

	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		var res any
		res, err = conn.Do(ctx, args...).Result()
		val = cast.ToStringMap(res)
		return err
	}, acceptable)
	return
}

// FTSpellCheck - Checks a query string for spelling errors.
// For more details about spellcheck query please follow:
// https://redis.io/docs/interact/search-and-query/advanced-concepts/spellcheck/
// For more information, please refer to the Redis documentation:
// [FT.SPELLCHECK]: (https://redis.io/commands/ft.spellcheck/)
func (r *Client) FTSpellCheck(ctx context.Context, index string, query string) (val map[string]any, err error) {
	args := []any{"FT.SPELLCHECK", index, query}
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		var res any
		res, err = conn.Do(ctx, args...).Result()
		val = cast.ToStringMap(res)
		return err
	}, acceptable)
	return
}

// FTSpellCheckWithArgs - Checks a query string for spelling errors with additional options.
// For more details about spellcheck query please follow:
// https://redis.io/docs/interact/search-and-query/advanced-concepts/spellcheck/
// For more information, please refer to the Redis documentation:
// [FT.SPELLCHECK]: (https://redis.io/commands/ft.spellcheck/)
func (r *Client) FTSpellCheckWithArgs(ctx context.Context, index string, query string, options *FTSpellCheckOptions) (val map[string]any, err error) {
	args := []any{"FT.SPELLCHECK", index, query}
	if options != nil {
		if options.Distance > 4 {
			panic("FT.SPELLCHECK: DISTANCE must be between 0 and 4")
		}
		if options.Distance > 0 {
			args = append(args, "DISTANCE", options.Distance)
		}
		if options.Terms.Include && options.Terms.Exclude {
			panic("FT.SPELLCHECK: INCLUDE and EXCLUDE are mutually exclusive")
		}
		if options.Terms.Include {
			args = append(args, "TERMS", "INCLUDE")
		}
		if options.Terms.Exclude {
			args = append(args, "TERMS", "EXCLUDE")
		}
		if options.Terms.Dictionary != "" {
			args = append(args, options.Terms.Dictionary)
		}
		if options.Dialect > 0 {
			args = append(args, "DIALECT", options.Dialect)
		}
	}
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		var res any
		res, err = conn.Do(ctx, args...).Result()
		val = cast.ToStringMap(res)
		return err
	}, acceptable)
	return
}

// FTSearch - Executes a search query on an index.
// The 'index' parameter specifies the index to search, and the 'query' parameter specifies the search query.
// For more information, please refer to the Redis documentation:
// [FT.SEARCH]: (https://redis.io/commands/ft.search/)
func (r *Client) FTSearch(ctx context.Context, index string, query string) (val any, err error) {
	args := []any{"FT.SEARCH", index, query}
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Do(ctx, args...).Result()
		return err
	}, acceptable)
	return
}

type SearchQuery []any

func FTSearchQuery(query string, options *FTSearchOptions) SearchQuery {
	queryArgs := []any{query}
	if options != nil {
		if options.NoContent {
			queryArgs = append(queryArgs, "NOCONTENT")
		}
		if options.Verbatim {
			queryArgs = append(queryArgs, "VERBATIM")
		}
		if options.NoStopWrods {
			queryArgs = append(queryArgs, "NOSTOPWORDS")
		}
		if options.WithScores {
			queryArgs = append(queryArgs, "WITHSCORES")
		}
		if options.WithPayloads {
			queryArgs = append(queryArgs, "WITHPAYLOADS")
		}
		if options.WithSortKeys {
			queryArgs = append(queryArgs, "WITHSORTKEYS")
		}
		if options.Filters != nil {
			for _, filter := range options.Filters {
				queryArgs = append(queryArgs, "FILTER", filter.FieldName, filter.Min, filter.Max)
			}
		}
		if options.GeoFilter != nil {
			for _, geoFilter := range options.GeoFilter {
				queryArgs = append(queryArgs, "GEOFILTER", geoFilter.FieldName, geoFilter.Longitude, geoFilter.Latitude, geoFilter.Radius, geoFilter.Unit)
			}
		}
		if options.InKeys != nil {
			queryArgs = append(queryArgs, "INKEYS", len(options.InKeys))
			queryArgs = append(queryArgs, options.InKeys...)
		}
		if options.InFields != nil {
			queryArgs = append(queryArgs, "INFIELDS", len(options.InFields))
			queryArgs = append(queryArgs, options.InFields...)
		}
		if options.Return != nil {
			queryArgs = append(queryArgs, "RETURN")
			queryArgsReturn := []any{}
			for _, ret := range options.Return {
				queryArgsReturn = append(queryArgsReturn, ret.FieldName)
				if ret.As != "" {
					queryArgsReturn = append(queryArgsReturn, "AS", ret.As)
				}
			}
			queryArgs = append(queryArgs, len(queryArgsReturn))
			queryArgs = append(queryArgs, queryArgsReturn...)
		}
		if options.Slop > 0 {
			queryArgs = append(queryArgs, "SLOP", options.Slop)
		}
		if options.Timeout > 0 {
			queryArgs = append(queryArgs, "TIMEOUT", options.Timeout)
		}
		if options.InOrder {
			queryArgs = append(queryArgs, "INORDER")
		}
		if options.Language != "" {
			queryArgs = append(queryArgs, "LANGUAGE", options.Language)
		}
		if options.Expander != "" {
			queryArgs = append(queryArgs, "EXPANDER", options.Expander)
		}
		if options.Scorer != "" {
			queryArgs = append(queryArgs, "SCORER", options.Scorer)
		}
		if options.ExplainScore {
			queryArgs = append(queryArgs, "EXPLAINSCORE")
		}
		if options.Payload != "" {
			queryArgs = append(queryArgs, "PAYLOAD", options.Payload)
		}
		if options.SortBy != nil {
			queryArgs = append(queryArgs, "SORTBY")
			for _, sortBy := range options.SortBy {
				queryArgs = append(queryArgs, sortBy.FieldName)
				if sortBy.Asc && sortBy.Desc {
					panic("FT.SEARCH: ASC and DESC are mutually exclusive")
				}
				if sortBy.Asc {
					queryArgs = append(queryArgs, "ASC")
				}
				if sortBy.Desc {
					queryArgs = append(queryArgs, "DESC")
				}
			}
			if options.SortByWithCount {
				queryArgs = append(queryArgs, "WITHCOUT")
			}
		}
		if options.LimitOffset >= 0 && options.Limit > 0 {
			queryArgs = append(queryArgs, "LIMIT", options.LimitOffset, options.Limit)
		}
		if options.Params != nil {
			queryArgs = append(queryArgs, "PARAMS", len(options.Params)*2)
			for key, value := range options.Params {
				queryArgs = append(queryArgs, key, value)
			}
		}
		if options.DialectVersion > 0 {
			queryArgs = append(queryArgs, "DIALECT", options.DialectVersion)
		}
	}
	return queryArgs
}

// FTSearchWithArgs - Executes a search query on an index with additional options.
// The 'index' parameter specifies the index to search, the 'query' parameter specifies the search query,
// and the 'options' parameter specifies additional options for the search.
// For more information, please refer to the Redis documentation:
// [FT.SEARCH]: (https://redis.io/commands/ft.search/)
func (r *Client) FTSearchWithArgs(ctx context.Context, index string, query string, options *FTSearchOptions) (val any, err error) {
	args := []any{"FT.SEARCH", index, query}
	if options != nil {
		if options.NoContent {
			args = append(args, "NOCONTENT")
		}
		if options.Verbatim {
			args = append(args, "VERBATIM")
		}
		if options.NoStopWrods {
			args = append(args, "NOSTOPWORDS")
		}
		if options.WithScores {
			args = append(args, "WITHSCORES")
		}
		if options.WithPayloads {
			args = append(args, "WITHPAYLOADS")
		}
		if options.WithSortKeys {
			args = append(args, "WITHSORTKEYS")
		}
		if options.Filters != nil {
			for _, filter := range options.Filters {
				args = append(args, "FILTER", filter.FieldName, filter.Min, filter.Max)
			}
		}
		if options.GeoFilter != nil {
			for _, geoFilter := range options.GeoFilter {
				args = append(args, "GEOFILTER", geoFilter.FieldName, geoFilter.Longitude, geoFilter.Latitude, geoFilter.Radius, geoFilter.Unit)
			}
		}
		if options.InKeys != nil {
			args = append(args, "INKEYS", len(options.InKeys))
			args = append(args, options.InKeys...)
		}
		if options.InFields != nil {
			args = append(args, "INFIELDS", len(options.InFields))
			args = append(args, options.InFields...)
		}
		if options.Return != nil {
			args = append(args, "RETURN")
			argsReturn := []any{}
			for _, ret := range options.Return {
				argsReturn = append(argsReturn, ret.FieldName)
				if ret.As != "" {
					argsReturn = append(argsReturn, "AS", ret.As)
				}
			}
			args = append(args, len(argsReturn))
			args = append(args, argsReturn...)
		}
		if options.Slop > 0 {
			args = append(args, "SLOP", options.Slop)
		}
		if options.Timeout > 0 {
			args = append(args, "TIMEOUT", options.Timeout)
		}
		if options.InOrder {
			args = append(args, "INORDER")
		}
		if options.Language != "" {
			args = append(args, "LANGUAGE", options.Language)
		}
		if options.Expander != "" {
			args = append(args, "EXPANDER", options.Expander)
		}
		if options.Scorer != "" {
			args = append(args, "SCORER", options.Scorer)
		}
		if options.ExplainScore {
			args = append(args, "EXPLAINSCORE")
		}
		if options.Payload != "" {
			args = append(args, "PAYLOAD", options.Payload)
		}
		if options.SortBy != nil {
			args = append(args, "SORTBY")
			for _, sortBy := range options.SortBy {
				args = append(args, sortBy.FieldName)
				if sortBy.Asc && sortBy.Desc {
					panic("FT.SEARCH: ASC and DESC are mutually exclusive")
				}
				if sortBy.Asc {
					args = append(args, "ASC")
				}
				if sortBy.Desc {
					args = append(args, "DESC")
				}
			}
			if options.SortByWithCount {
				args = append(args, "WITHCOUT")
			}
		}
		if options.LimitOffset >= 0 && options.Limit > 0 {
			args = append(args, "LIMIT", options.LimitOffset, options.Limit)
		}
		if options.Params != nil {
			args = append(args, "PARAMS", len(options.Params)*2)
			for key, value := range options.Params {
				args = append(args, key, value)
			}
		}
		if options.DialectVersion > 0 {
			args = append(args, "DIALECT", options.DialectVersion)
		}
	}
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Do(ctx, args...).Result()
		return err
	}, acceptable)
	return
}

// FTSynDump - Dumps the synonyms data structure.
// The 'index' parameter specifies the index to dump.
// For more information, please refer to the Redis documentation:
// [FT.SYNDUMP]: (https://redis.io/commands/ft.syndump/)
func (r *Client) FTSynDump(ctx context.Context, index string) (val map[string]any, err error) {
	args := []any{"FT.SYNDUMP", index}
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		var res any
		res, err = conn.Do(ctx, args...).Result()
		val = cast.ToStringMap(res)
		return err
	}, acceptable)
	return
}

// FTSynUpdate - Creates or updates a synonym group with additional terms.
// The 'index' parameter specifies the index to update, the 'synGroupId' parameter specifies the synonym group id, and the 'terms' parameter specifies the additional terms.
// For more information, please refer to the Redis documentation:
// [FT.SYNUPDATE]: (https://redis.io/commands/ft.synupdate/)
func (r *Client) FTSynUpdate(ctx context.Context, index string, synGroupId any, terms []any) (val string, err error) {
	args := []any{"FT.SYNUPDATE", index, synGroupId}
	args = append(args, terms...)
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Do(ctx, args...).Text()
		return err
	}, acceptable)
	return
}

// FTSynUpdateWithArgs - Creates or updates a synonym group with additional terms and options.
// The 'index' parameter specifies the index to update, the 'synGroupId' parameter specifies the synonym group id, the 'options' parameter specifies additional options for the update, and the 'terms' parameter specifies the additional terms.
// For more information, please refer to the Redis documentation:
// [FT.SYNUPDATE]: (https://redis.io/commands/ft.synupdate/)
func (r *Client) FTSynUpdateWithArgs(ctx context.Context, index string, synGroupId any, options *FTSynUpdateOptions, terms []any) (val string, err error) {
	args := []any{"FT.SYNUPDATE", index, synGroupId}
	if options.SkipInitialScan {
		args = append(args, "SKIPINITIALSCAN")
	}
	args = append(args, terms...)
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Do(ctx, args...).Text()
		return err
	}, acceptable)
	return
}

// FTTagVals - Returns all distinct values indexed in a tag field.
// The 'index' parameter specifies the index to check, and the 'field' parameter specifies the tag field to retrieve values from.
// For more information, please refer to the Redis documentation:
// [FT.TAGVALS]: (https://redis.io/commands/ft.tagvals/)
func (r *Client) FTTagVals(ctx context.Context, index string, field string) (val map[string]any, err error) {
	args := []any{"FT.TAGVALS", index, field}
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		var res any
		res, err = conn.Do(ctx, args...).Result()
		val = cast.ToStringMap(res)
		return err
	}, acceptable)
	return
}
