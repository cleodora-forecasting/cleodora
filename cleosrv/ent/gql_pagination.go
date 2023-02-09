// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/errcode"
	"github.com/cleodora-forecasting/cleodora/cleosrv/ent/estimate"
	"github.com/cleodora-forecasting/cleodora/cleosrv/ent/forecast"
	"github.com/cleodora-forecasting/cleodora/cleosrv/ent/outcome"
	"github.com/cleodora-forecasting/cleodora/cleosrv/ent/probability"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/vmihailenco/msgpack/v5"
)

// OrderDirection defines the directions in which to order a list of items.
type OrderDirection string

const (
	// OrderDirectionAsc specifies an ascending order.
	OrderDirectionAsc OrderDirection = "ASC"
	// OrderDirectionDesc specifies a descending order.
	OrderDirectionDesc OrderDirection = "DESC"
)

// Validate the order direction value.
func (o OrderDirection) Validate() error {
	if o != OrderDirectionAsc && o != OrderDirectionDesc {
		return fmt.Errorf("%s is not a valid OrderDirection", o)
	}
	return nil
}

// String implements fmt.Stringer interface.
func (o OrderDirection) String() string {
	return string(o)
}

// MarshalGQL implements graphql.Marshaler interface.
func (o OrderDirection) MarshalGQL(w io.Writer) {
	io.WriteString(w, strconv.Quote(o.String()))
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (o *OrderDirection) UnmarshalGQL(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return fmt.Errorf("order direction %T must be a string", val)
	}
	*o = OrderDirection(str)
	return o.Validate()
}

func (o OrderDirection) reverse() OrderDirection {
	if o == OrderDirectionDesc {
		return OrderDirectionAsc
	}
	return OrderDirectionDesc
}

func (o OrderDirection) orderFunc(field string) OrderFunc {
	if o == OrderDirectionDesc {
		return Desc(field)
	}
	return Asc(field)
}

func cursorsToPredicates(direction OrderDirection, after, before *Cursor, field, idField string) []func(s *sql.Selector) {
	var predicates []func(s *sql.Selector)
	if after != nil {
		if after.Value != nil {
			var predicate func([]string, ...interface{}) *sql.Predicate
			if direction == OrderDirectionAsc {
				predicate = sql.CompositeGT
			} else {
				predicate = sql.CompositeLT
			}
			predicates = append(predicates, func(s *sql.Selector) {
				s.Where(predicate(
					s.Columns(field, idField),
					after.Value, after.ID,
				))
			})
		} else {
			var predicate func(string, interface{}) *sql.Predicate
			if direction == OrderDirectionAsc {
				predicate = sql.GT
			} else {
				predicate = sql.LT
			}
			predicates = append(predicates, func(s *sql.Selector) {
				s.Where(predicate(
					s.C(idField),
					after.ID,
				))
			})
		}
	}
	if before != nil {
		if before.Value != nil {
			var predicate func([]string, ...interface{}) *sql.Predicate
			if direction == OrderDirectionAsc {
				predicate = sql.CompositeLT
			} else {
				predicate = sql.CompositeGT
			}
			predicates = append(predicates, func(s *sql.Selector) {
				s.Where(predicate(
					s.Columns(field, idField),
					before.Value, before.ID,
				))
			})
		} else {
			var predicate func(string, interface{}) *sql.Predicate
			if direction == OrderDirectionAsc {
				predicate = sql.LT
			} else {
				predicate = sql.GT
			}
			predicates = append(predicates, func(s *sql.Selector) {
				s.Where(predicate(
					s.C(idField),
					before.ID,
				))
			})
		}
	}
	return predicates
}

// PageInfo of a connection type.
type PageInfo struct {
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
	StartCursor     *Cursor `json:"startCursor"`
	EndCursor       *Cursor `json:"endCursor"`
}

// Cursor of an edge type.
type Cursor struct {
	ID    int   `msgpack:"i"`
	Value Value `msgpack:"v,omitempty"`
}

// MarshalGQL implements graphql.Marshaler interface.
func (c Cursor) MarshalGQL(w io.Writer) {
	quote := []byte{'"'}
	w.Write(quote)
	defer w.Write(quote)
	wc := base64.NewEncoder(base64.RawStdEncoding, w)
	defer wc.Close()
	_ = msgpack.NewEncoder(wc).Encode(c)
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (c *Cursor) UnmarshalGQL(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("%T is not a string", v)
	}
	if err := msgpack.NewDecoder(
		base64.NewDecoder(
			base64.RawStdEncoding,
			strings.NewReader(s),
		),
	).Decode(c); err != nil {
		return fmt.Errorf("cannot decode cursor: %w", err)
	}
	return nil
}

const errInvalidPagination = "INVALID_PAGINATION"

func validateFirstLast(first, last *int) (err *gqlerror.Error) {
	switch {
	case first != nil && last != nil:
		err = &gqlerror.Error{
			Message: "Passing both `first` and `last` to paginate a connection is not supported.",
		}
	case first != nil && *first < 0:
		err = &gqlerror.Error{
			Message: "`first` on a connection cannot be less than zero.",
		}
		errcode.Set(err, errInvalidPagination)
	case last != nil && *last < 0:
		err = &gqlerror.Error{
			Message: "`last` on a connection cannot be less than zero.",
		}
		errcode.Set(err, errInvalidPagination)
	}
	return err
}

func collectedField(ctx context.Context, path ...string) *graphql.CollectedField {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return nil
	}
	field := fc.Field
	oc := graphql.GetOperationContext(ctx)
walk:
	for _, name := range path {
		for _, f := range graphql.CollectFields(oc, field.Selections, nil) {
			if f.Alias == name {
				field = f
				continue walk
			}
		}
		return nil
	}
	return &field
}

func hasCollectedField(ctx context.Context, path ...string) bool {
	if graphql.GetFieldContext(ctx) == nil {
		return true
	}
	return collectedField(ctx, path...) != nil
}

const (
	edgesField      = "edges"
	nodeField       = "node"
	pageInfoField   = "pageInfo"
	totalCountField = "totalCount"
)

func paginateLimit(first, last *int) int {
	var limit int
	if first != nil {
		limit = *first + 1
	} else if last != nil {
		limit = *last + 1
	}
	return limit
}

// EstimateEdge is the edge representation of Estimate.
type EstimateEdge struct {
	Node   *Estimate `json:"node"`
	Cursor Cursor    `json:"cursor"`
}

// EstimateConnection is the connection containing edges to Estimate.
type EstimateConnection struct {
	Edges      []*EstimateEdge `json:"edges"`
	PageInfo   PageInfo        `json:"pageInfo"`
	TotalCount int             `json:"totalCount"`
}

func (c *EstimateConnection) build(nodes []*Estimate, pager *estimatePager, after *Cursor, first *int, before *Cursor, last *int) {
	c.PageInfo.HasNextPage = before != nil
	c.PageInfo.HasPreviousPage = after != nil
	if first != nil && *first+1 == len(nodes) {
		c.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && *last+1 == len(nodes) {
		c.PageInfo.HasPreviousPage = true
		nodes = nodes[:len(nodes)-1]
	}
	var nodeAt func(int) *Estimate
	if last != nil {
		n := len(nodes) - 1
		nodeAt = func(i int) *Estimate {
			return nodes[n-i]
		}
	} else {
		nodeAt = func(i int) *Estimate {
			return nodes[i]
		}
	}
	c.Edges = make([]*EstimateEdge, len(nodes))
	for i := range nodes {
		node := nodeAt(i)
		c.Edges[i] = &EstimateEdge{
			Node:   node,
			Cursor: pager.toCursor(node),
		}
	}
	if l := len(c.Edges); l > 0 {
		c.PageInfo.StartCursor = &c.Edges[0].Cursor
		c.PageInfo.EndCursor = &c.Edges[l-1].Cursor
	}
	if c.TotalCount == 0 {
		c.TotalCount = len(nodes)
	}
}

// EstimatePaginateOption enables pagination customization.
type EstimatePaginateOption func(*estimatePager) error

// WithEstimateOrder configures pagination ordering.
func WithEstimateOrder(order *EstimateOrder) EstimatePaginateOption {
	if order == nil {
		order = DefaultEstimateOrder
	}
	o := *order
	return func(pager *estimatePager) error {
		if err := o.Direction.Validate(); err != nil {
			return err
		}
		if o.Field == nil {
			o.Field = DefaultEstimateOrder.Field
		}
		pager.order = &o
		return nil
	}
}

// WithEstimateFilter configures pagination filter.
func WithEstimateFilter(filter func(*EstimateQuery) (*EstimateQuery, error)) EstimatePaginateOption {
	return func(pager *estimatePager) error {
		if filter == nil {
			return errors.New("EstimateQuery filter cannot be nil")
		}
		pager.filter = filter
		return nil
	}
}

type estimatePager struct {
	order  *EstimateOrder
	filter func(*EstimateQuery) (*EstimateQuery, error)
}

func newEstimatePager(opts []EstimatePaginateOption) (*estimatePager, error) {
	pager := &estimatePager{}
	for _, opt := range opts {
		if err := opt(pager); err != nil {
			return nil, err
		}
	}
	if pager.order == nil {
		pager.order = DefaultEstimateOrder
	}
	return pager, nil
}

func (p *estimatePager) applyFilter(query *EstimateQuery) (*EstimateQuery, error) {
	if p.filter != nil {
		return p.filter(query)
	}
	return query, nil
}

func (p *estimatePager) toCursor(e *Estimate) Cursor {
	return p.order.Field.toCursor(e)
}

func (p *estimatePager) applyCursors(query *EstimateQuery, after, before *Cursor) *EstimateQuery {
	for _, predicate := range cursorsToPredicates(
		p.order.Direction, after, before,
		p.order.Field.field, DefaultEstimateOrder.Field.field,
	) {
		query = query.Where(predicate)
	}
	return query
}

func (p *estimatePager) applyOrder(query *EstimateQuery, reverse bool) *EstimateQuery {
	direction := p.order.Direction
	if reverse {
		direction = direction.reverse()
	}
	query = query.Order(direction.orderFunc(p.order.Field.field))
	if p.order.Field != DefaultEstimateOrder.Field {
		query = query.Order(direction.orderFunc(DefaultEstimateOrder.Field.field))
	}
	return query
}

func (p *estimatePager) orderExpr(reverse bool) sql.Querier {
	direction := p.order.Direction
	if reverse {
		direction = direction.reverse()
	}
	return sql.ExprFunc(func(b *sql.Builder) {
		b.Ident(p.order.Field.field).Pad().WriteString(string(direction))
		if p.order.Field != DefaultEstimateOrder.Field {
			b.Comma().Ident(DefaultEstimateOrder.Field.field).Pad().WriteString(string(direction))
		}
	})
}

// Paginate executes the query and returns a relay based cursor connection to Estimate.
func (e *EstimateQuery) Paginate(
	ctx context.Context, after *Cursor, first *int,
	before *Cursor, last *int, opts ...EstimatePaginateOption,
) (*EstimateConnection, error) {
	if err := validateFirstLast(first, last); err != nil {
		return nil, err
	}
	pager, err := newEstimatePager(opts)
	if err != nil {
		return nil, err
	}
	if e, err = pager.applyFilter(e); err != nil {
		return nil, err
	}
	conn := &EstimateConnection{Edges: []*EstimateEdge{}}
	ignoredEdges := !hasCollectedField(ctx, edgesField)
	if hasCollectedField(ctx, totalCountField) || hasCollectedField(ctx, pageInfoField) {
		hasPagination := after != nil || first != nil || before != nil || last != nil
		if hasPagination || ignoredEdges {
			if conn.TotalCount, err = e.Clone().Count(ctx); err != nil {
				return nil, err
			}
			conn.PageInfo.HasNextPage = first != nil && conn.TotalCount > 0
			conn.PageInfo.HasPreviousPage = last != nil && conn.TotalCount > 0
		}
	}
	if ignoredEdges || (first != nil && *first == 0) || (last != nil && *last == 0) {
		return conn, nil
	}

	e = pager.applyCursors(e, after, before)
	e = pager.applyOrder(e, last != nil)
	if limit := paginateLimit(first, last); limit != 0 {
		e.Limit(limit)
	}
	if field := collectedField(ctx, edgesField, nodeField); field != nil {
		if err := e.collectField(ctx, graphql.GetOperationContext(ctx), *field, []string{edgesField, nodeField}); err != nil {
			return nil, err
		}
	}

	nodes, err := e.All(ctx)
	if err != nil {
		return nil, err
	}
	conn.build(nodes, pager, after, first, before, last)
	return conn, nil
}

// EstimateOrderField defines the ordering field of Estimate.
type EstimateOrderField struct {
	field    string
	toCursor func(*Estimate) Cursor
}

// EstimateOrder defines the ordering of Estimate.
type EstimateOrder struct {
	Direction OrderDirection      `json:"direction"`
	Field     *EstimateOrderField `json:"field"`
}

// DefaultEstimateOrder is the default ordering of Estimate.
var DefaultEstimateOrder = &EstimateOrder{
	Direction: OrderDirectionAsc,
	Field: &EstimateOrderField{
		field: estimate.FieldID,
		toCursor: func(e *Estimate) Cursor {
			return Cursor{ID: e.ID}
		},
	},
}

// ToEdge converts Estimate into EstimateEdge.
func (e *Estimate) ToEdge(order *EstimateOrder) *EstimateEdge {
	if order == nil {
		order = DefaultEstimateOrder
	}
	return &EstimateEdge{
		Node:   e,
		Cursor: order.Field.toCursor(e),
	}
}

// ForecastEdge is the edge representation of Forecast.
type ForecastEdge struct {
	Node   *Forecast `json:"node"`
	Cursor Cursor    `json:"cursor"`
}

// ForecastConnection is the connection containing edges to Forecast.
type ForecastConnection struct {
	Edges      []*ForecastEdge `json:"edges"`
	PageInfo   PageInfo        `json:"pageInfo"`
	TotalCount int             `json:"totalCount"`
}

func (c *ForecastConnection) build(nodes []*Forecast, pager *forecastPager, after *Cursor, first *int, before *Cursor, last *int) {
	c.PageInfo.HasNextPage = before != nil
	c.PageInfo.HasPreviousPage = after != nil
	if first != nil && *first+1 == len(nodes) {
		c.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && *last+1 == len(nodes) {
		c.PageInfo.HasPreviousPage = true
		nodes = nodes[:len(nodes)-1]
	}
	var nodeAt func(int) *Forecast
	if last != nil {
		n := len(nodes) - 1
		nodeAt = func(i int) *Forecast {
			return nodes[n-i]
		}
	} else {
		nodeAt = func(i int) *Forecast {
			return nodes[i]
		}
	}
	c.Edges = make([]*ForecastEdge, len(nodes))
	for i := range nodes {
		node := nodeAt(i)
		c.Edges[i] = &ForecastEdge{
			Node:   node,
			Cursor: pager.toCursor(node),
		}
	}
	if l := len(c.Edges); l > 0 {
		c.PageInfo.StartCursor = &c.Edges[0].Cursor
		c.PageInfo.EndCursor = &c.Edges[l-1].Cursor
	}
	if c.TotalCount == 0 {
		c.TotalCount = len(nodes)
	}
}

// ForecastPaginateOption enables pagination customization.
type ForecastPaginateOption func(*forecastPager) error

// WithForecastOrder configures pagination ordering.
func WithForecastOrder(order *ForecastOrder) ForecastPaginateOption {
	if order == nil {
		order = DefaultForecastOrder
	}
	o := *order
	return func(pager *forecastPager) error {
		if err := o.Direction.Validate(); err != nil {
			return err
		}
		if o.Field == nil {
			o.Field = DefaultForecastOrder.Field
		}
		pager.order = &o
		return nil
	}
}

// WithForecastFilter configures pagination filter.
func WithForecastFilter(filter func(*ForecastQuery) (*ForecastQuery, error)) ForecastPaginateOption {
	return func(pager *forecastPager) error {
		if filter == nil {
			return errors.New("ForecastQuery filter cannot be nil")
		}
		pager.filter = filter
		return nil
	}
}

type forecastPager struct {
	order  *ForecastOrder
	filter func(*ForecastQuery) (*ForecastQuery, error)
}

func newForecastPager(opts []ForecastPaginateOption) (*forecastPager, error) {
	pager := &forecastPager{}
	for _, opt := range opts {
		if err := opt(pager); err != nil {
			return nil, err
		}
	}
	if pager.order == nil {
		pager.order = DefaultForecastOrder
	}
	return pager, nil
}

func (p *forecastPager) applyFilter(query *ForecastQuery) (*ForecastQuery, error) {
	if p.filter != nil {
		return p.filter(query)
	}
	return query, nil
}

func (p *forecastPager) toCursor(f *Forecast) Cursor {
	return p.order.Field.toCursor(f)
}

func (p *forecastPager) applyCursors(query *ForecastQuery, after, before *Cursor) *ForecastQuery {
	for _, predicate := range cursorsToPredicates(
		p.order.Direction, after, before,
		p.order.Field.field, DefaultForecastOrder.Field.field,
	) {
		query = query.Where(predicate)
	}
	return query
}

func (p *forecastPager) applyOrder(query *ForecastQuery, reverse bool) *ForecastQuery {
	direction := p.order.Direction
	if reverse {
		direction = direction.reverse()
	}
	query = query.Order(direction.orderFunc(p.order.Field.field))
	if p.order.Field != DefaultForecastOrder.Field {
		query = query.Order(direction.orderFunc(DefaultForecastOrder.Field.field))
	}
	return query
}

func (p *forecastPager) orderExpr(reverse bool) sql.Querier {
	direction := p.order.Direction
	if reverse {
		direction = direction.reverse()
	}
	return sql.ExprFunc(func(b *sql.Builder) {
		b.Ident(p.order.Field.field).Pad().WriteString(string(direction))
		if p.order.Field != DefaultForecastOrder.Field {
			b.Comma().Ident(DefaultForecastOrder.Field.field).Pad().WriteString(string(direction))
		}
	})
}

// Paginate executes the query and returns a relay based cursor connection to Forecast.
func (f *ForecastQuery) Paginate(
	ctx context.Context, after *Cursor, first *int,
	before *Cursor, last *int, opts ...ForecastPaginateOption,
) (*ForecastConnection, error) {
	if err := validateFirstLast(first, last); err != nil {
		return nil, err
	}
	pager, err := newForecastPager(opts)
	if err != nil {
		return nil, err
	}
	if f, err = pager.applyFilter(f); err != nil {
		return nil, err
	}
	conn := &ForecastConnection{Edges: []*ForecastEdge{}}
	ignoredEdges := !hasCollectedField(ctx, edgesField)
	if hasCollectedField(ctx, totalCountField) || hasCollectedField(ctx, pageInfoField) {
		hasPagination := after != nil || first != nil || before != nil || last != nil
		if hasPagination || ignoredEdges {
			if conn.TotalCount, err = f.Clone().Count(ctx); err != nil {
				return nil, err
			}
			conn.PageInfo.HasNextPage = first != nil && conn.TotalCount > 0
			conn.PageInfo.HasPreviousPage = last != nil && conn.TotalCount > 0
		}
	}
	if ignoredEdges || (first != nil && *first == 0) || (last != nil && *last == 0) {
		return conn, nil
	}

	f = pager.applyCursors(f, after, before)
	f = pager.applyOrder(f, last != nil)
	if limit := paginateLimit(first, last); limit != 0 {
		f.Limit(limit)
	}
	if field := collectedField(ctx, edgesField, nodeField); field != nil {
		if err := f.collectField(ctx, graphql.GetOperationContext(ctx), *field, []string{edgesField, nodeField}); err != nil {
			return nil, err
		}
	}

	nodes, err := f.All(ctx)
	if err != nil {
		return nil, err
	}
	conn.build(nodes, pager, after, first, before, last)
	return conn, nil
}

// ForecastOrderField defines the ordering field of Forecast.
type ForecastOrderField struct {
	field    string
	toCursor func(*Forecast) Cursor
}

// ForecastOrder defines the ordering of Forecast.
type ForecastOrder struct {
	Direction OrderDirection      `json:"direction"`
	Field     *ForecastOrderField `json:"field"`
}

// DefaultForecastOrder is the default ordering of Forecast.
var DefaultForecastOrder = &ForecastOrder{
	Direction: OrderDirectionAsc,
	Field: &ForecastOrderField{
		field: forecast.FieldID,
		toCursor: func(f *Forecast) Cursor {
			return Cursor{ID: f.ID}
		},
	},
}

// ToEdge converts Forecast into ForecastEdge.
func (f *Forecast) ToEdge(order *ForecastOrder) *ForecastEdge {
	if order == nil {
		order = DefaultForecastOrder
	}
	return &ForecastEdge{
		Node:   f,
		Cursor: order.Field.toCursor(f),
	}
}

// OutcomeEdge is the edge representation of Outcome.
type OutcomeEdge struct {
	Node   *Outcome `json:"node"`
	Cursor Cursor   `json:"cursor"`
}

// OutcomeConnection is the connection containing edges to Outcome.
type OutcomeConnection struct {
	Edges      []*OutcomeEdge `json:"edges"`
	PageInfo   PageInfo       `json:"pageInfo"`
	TotalCount int            `json:"totalCount"`
}

func (c *OutcomeConnection) build(nodes []*Outcome, pager *outcomePager, after *Cursor, first *int, before *Cursor, last *int) {
	c.PageInfo.HasNextPage = before != nil
	c.PageInfo.HasPreviousPage = after != nil
	if first != nil && *first+1 == len(nodes) {
		c.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && *last+1 == len(nodes) {
		c.PageInfo.HasPreviousPage = true
		nodes = nodes[:len(nodes)-1]
	}
	var nodeAt func(int) *Outcome
	if last != nil {
		n := len(nodes) - 1
		nodeAt = func(i int) *Outcome {
			return nodes[n-i]
		}
	} else {
		nodeAt = func(i int) *Outcome {
			return nodes[i]
		}
	}
	c.Edges = make([]*OutcomeEdge, len(nodes))
	for i := range nodes {
		node := nodeAt(i)
		c.Edges[i] = &OutcomeEdge{
			Node:   node,
			Cursor: pager.toCursor(node),
		}
	}
	if l := len(c.Edges); l > 0 {
		c.PageInfo.StartCursor = &c.Edges[0].Cursor
		c.PageInfo.EndCursor = &c.Edges[l-1].Cursor
	}
	if c.TotalCount == 0 {
		c.TotalCount = len(nodes)
	}
}

// OutcomePaginateOption enables pagination customization.
type OutcomePaginateOption func(*outcomePager) error

// WithOutcomeOrder configures pagination ordering.
func WithOutcomeOrder(order *OutcomeOrder) OutcomePaginateOption {
	if order == nil {
		order = DefaultOutcomeOrder
	}
	o := *order
	return func(pager *outcomePager) error {
		if err := o.Direction.Validate(); err != nil {
			return err
		}
		if o.Field == nil {
			o.Field = DefaultOutcomeOrder.Field
		}
		pager.order = &o
		return nil
	}
}

// WithOutcomeFilter configures pagination filter.
func WithOutcomeFilter(filter func(*OutcomeQuery) (*OutcomeQuery, error)) OutcomePaginateOption {
	return func(pager *outcomePager) error {
		if filter == nil {
			return errors.New("OutcomeQuery filter cannot be nil")
		}
		pager.filter = filter
		return nil
	}
}

type outcomePager struct {
	order  *OutcomeOrder
	filter func(*OutcomeQuery) (*OutcomeQuery, error)
}

func newOutcomePager(opts []OutcomePaginateOption) (*outcomePager, error) {
	pager := &outcomePager{}
	for _, opt := range opts {
		if err := opt(pager); err != nil {
			return nil, err
		}
	}
	if pager.order == nil {
		pager.order = DefaultOutcomeOrder
	}
	return pager, nil
}

func (p *outcomePager) applyFilter(query *OutcomeQuery) (*OutcomeQuery, error) {
	if p.filter != nil {
		return p.filter(query)
	}
	return query, nil
}

func (p *outcomePager) toCursor(o *Outcome) Cursor {
	return p.order.Field.toCursor(o)
}

func (p *outcomePager) applyCursors(query *OutcomeQuery, after, before *Cursor) *OutcomeQuery {
	for _, predicate := range cursorsToPredicates(
		p.order.Direction, after, before,
		p.order.Field.field, DefaultOutcomeOrder.Field.field,
	) {
		query = query.Where(predicate)
	}
	return query
}

func (p *outcomePager) applyOrder(query *OutcomeQuery, reverse bool) *OutcomeQuery {
	direction := p.order.Direction
	if reverse {
		direction = direction.reverse()
	}
	query = query.Order(direction.orderFunc(p.order.Field.field))
	if p.order.Field != DefaultOutcomeOrder.Field {
		query = query.Order(direction.orderFunc(DefaultOutcomeOrder.Field.field))
	}
	return query
}

func (p *outcomePager) orderExpr(reverse bool) sql.Querier {
	direction := p.order.Direction
	if reverse {
		direction = direction.reverse()
	}
	return sql.ExprFunc(func(b *sql.Builder) {
		b.Ident(p.order.Field.field).Pad().WriteString(string(direction))
		if p.order.Field != DefaultOutcomeOrder.Field {
			b.Comma().Ident(DefaultOutcomeOrder.Field.field).Pad().WriteString(string(direction))
		}
	})
}

// Paginate executes the query and returns a relay based cursor connection to Outcome.
func (o *OutcomeQuery) Paginate(
	ctx context.Context, after *Cursor, first *int,
	before *Cursor, last *int, opts ...OutcomePaginateOption,
) (*OutcomeConnection, error) {
	if err := validateFirstLast(first, last); err != nil {
		return nil, err
	}
	pager, err := newOutcomePager(opts)
	if err != nil {
		return nil, err
	}
	if o, err = pager.applyFilter(o); err != nil {
		return nil, err
	}
	conn := &OutcomeConnection{Edges: []*OutcomeEdge{}}
	ignoredEdges := !hasCollectedField(ctx, edgesField)
	if hasCollectedField(ctx, totalCountField) || hasCollectedField(ctx, pageInfoField) {
		hasPagination := after != nil || first != nil || before != nil || last != nil
		if hasPagination || ignoredEdges {
			if conn.TotalCount, err = o.Clone().Count(ctx); err != nil {
				return nil, err
			}
			conn.PageInfo.HasNextPage = first != nil && conn.TotalCount > 0
			conn.PageInfo.HasPreviousPage = last != nil && conn.TotalCount > 0
		}
	}
	if ignoredEdges || (first != nil && *first == 0) || (last != nil && *last == 0) {
		return conn, nil
	}

	o = pager.applyCursors(o, after, before)
	o = pager.applyOrder(o, last != nil)
	if limit := paginateLimit(first, last); limit != 0 {
		o.Limit(limit)
	}
	if field := collectedField(ctx, edgesField, nodeField); field != nil {
		if err := o.collectField(ctx, graphql.GetOperationContext(ctx), *field, []string{edgesField, nodeField}); err != nil {
			return nil, err
		}
	}

	nodes, err := o.All(ctx)
	if err != nil {
		return nil, err
	}
	conn.build(nodes, pager, after, first, before, last)
	return conn, nil
}

// OutcomeOrderField defines the ordering field of Outcome.
type OutcomeOrderField struct {
	field    string
	toCursor func(*Outcome) Cursor
}

// OutcomeOrder defines the ordering of Outcome.
type OutcomeOrder struct {
	Direction OrderDirection     `json:"direction"`
	Field     *OutcomeOrderField `json:"field"`
}

// DefaultOutcomeOrder is the default ordering of Outcome.
var DefaultOutcomeOrder = &OutcomeOrder{
	Direction: OrderDirectionAsc,
	Field: &OutcomeOrderField{
		field: outcome.FieldID,
		toCursor: func(o *Outcome) Cursor {
			return Cursor{ID: o.ID}
		},
	},
}

// ToEdge converts Outcome into OutcomeEdge.
func (o *Outcome) ToEdge(order *OutcomeOrder) *OutcomeEdge {
	if order == nil {
		order = DefaultOutcomeOrder
	}
	return &OutcomeEdge{
		Node:   o,
		Cursor: order.Field.toCursor(o),
	}
}

// ProbabilityEdge is the edge representation of Probability.
type ProbabilityEdge struct {
	Node   *Probability `json:"node"`
	Cursor Cursor       `json:"cursor"`
}

// ProbabilityConnection is the connection containing edges to Probability.
type ProbabilityConnection struct {
	Edges      []*ProbabilityEdge `json:"edges"`
	PageInfo   PageInfo           `json:"pageInfo"`
	TotalCount int                `json:"totalCount"`
}

func (c *ProbabilityConnection) build(nodes []*Probability, pager *probabilityPager, after *Cursor, first *int, before *Cursor, last *int) {
	c.PageInfo.HasNextPage = before != nil
	c.PageInfo.HasPreviousPage = after != nil
	if first != nil && *first+1 == len(nodes) {
		c.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && *last+1 == len(nodes) {
		c.PageInfo.HasPreviousPage = true
		nodes = nodes[:len(nodes)-1]
	}
	var nodeAt func(int) *Probability
	if last != nil {
		n := len(nodes) - 1
		nodeAt = func(i int) *Probability {
			return nodes[n-i]
		}
	} else {
		nodeAt = func(i int) *Probability {
			return nodes[i]
		}
	}
	c.Edges = make([]*ProbabilityEdge, len(nodes))
	for i := range nodes {
		node := nodeAt(i)
		c.Edges[i] = &ProbabilityEdge{
			Node:   node,
			Cursor: pager.toCursor(node),
		}
	}
	if l := len(c.Edges); l > 0 {
		c.PageInfo.StartCursor = &c.Edges[0].Cursor
		c.PageInfo.EndCursor = &c.Edges[l-1].Cursor
	}
	if c.TotalCount == 0 {
		c.TotalCount = len(nodes)
	}
}

// ProbabilityPaginateOption enables pagination customization.
type ProbabilityPaginateOption func(*probabilityPager) error

// WithProbabilityOrder configures pagination ordering.
func WithProbabilityOrder(order *ProbabilityOrder) ProbabilityPaginateOption {
	if order == nil {
		order = DefaultProbabilityOrder
	}
	o := *order
	return func(pager *probabilityPager) error {
		if err := o.Direction.Validate(); err != nil {
			return err
		}
		if o.Field == nil {
			o.Field = DefaultProbabilityOrder.Field
		}
		pager.order = &o
		return nil
	}
}

// WithProbabilityFilter configures pagination filter.
func WithProbabilityFilter(filter func(*ProbabilityQuery) (*ProbabilityQuery, error)) ProbabilityPaginateOption {
	return func(pager *probabilityPager) error {
		if filter == nil {
			return errors.New("ProbabilityQuery filter cannot be nil")
		}
		pager.filter = filter
		return nil
	}
}

type probabilityPager struct {
	order  *ProbabilityOrder
	filter func(*ProbabilityQuery) (*ProbabilityQuery, error)
}

func newProbabilityPager(opts []ProbabilityPaginateOption) (*probabilityPager, error) {
	pager := &probabilityPager{}
	for _, opt := range opts {
		if err := opt(pager); err != nil {
			return nil, err
		}
	}
	if pager.order == nil {
		pager.order = DefaultProbabilityOrder
	}
	return pager, nil
}

func (p *probabilityPager) applyFilter(query *ProbabilityQuery) (*ProbabilityQuery, error) {
	if p.filter != nil {
		return p.filter(query)
	}
	return query, nil
}

func (p *probabilityPager) toCursor(pr *Probability) Cursor {
	return p.order.Field.toCursor(pr)
}

func (p *probabilityPager) applyCursors(query *ProbabilityQuery, after, before *Cursor) *ProbabilityQuery {
	for _, predicate := range cursorsToPredicates(
		p.order.Direction, after, before,
		p.order.Field.field, DefaultProbabilityOrder.Field.field,
	) {
		query = query.Where(predicate)
	}
	return query
}

func (p *probabilityPager) applyOrder(query *ProbabilityQuery, reverse bool) *ProbabilityQuery {
	direction := p.order.Direction
	if reverse {
		direction = direction.reverse()
	}
	query = query.Order(direction.orderFunc(p.order.Field.field))
	if p.order.Field != DefaultProbabilityOrder.Field {
		query = query.Order(direction.orderFunc(DefaultProbabilityOrder.Field.field))
	}
	return query
}

func (p *probabilityPager) orderExpr(reverse bool) sql.Querier {
	direction := p.order.Direction
	if reverse {
		direction = direction.reverse()
	}
	return sql.ExprFunc(func(b *sql.Builder) {
		b.Ident(p.order.Field.field).Pad().WriteString(string(direction))
		if p.order.Field != DefaultProbabilityOrder.Field {
			b.Comma().Ident(DefaultProbabilityOrder.Field.field).Pad().WriteString(string(direction))
		}
	})
}

// Paginate executes the query and returns a relay based cursor connection to Probability.
func (pr *ProbabilityQuery) Paginate(
	ctx context.Context, after *Cursor, first *int,
	before *Cursor, last *int, opts ...ProbabilityPaginateOption,
) (*ProbabilityConnection, error) {
	if err := validateFirstLast(first, last); err != nil {
		return nil, err
	}
	pager, err := newProbabilityPager(opts)
	if err != nil {
		return nil, err
	}
	if pr, err = pager.applyFilter(pr); err != nil {
		return nil, err
	}
	conn := &ProbabilityConnection{Edges: []*ProbabilityEdge{}}
	ignoredEdges := !hasCollectedField(ctx, edgesField)
	if hasCollectedField(ctx, totalCountField) || hasCollectedField(ctx, pageInfoField) {
		hasPagination := after != nil || first != nil || before != nil || last != nil
		if hasPagination || ignoredEdges {
			if conn.TotalCount, err = pr.Clone().Count(ctx); err != nil {
				return nil, err
			}
			conn.PageInfo.HasNextPage = first != nil && conn.TotalCount > 0
			conn.PageInfo.HasPreviousPage = last != nil && conn.TotalCount > 0
		}
	}
	if ignoredEdges || (first != nil && *first == 0) || (last != nil && *last == 0) {
		return conn, nil
	}

	pr = pager.applyCursors(pr, after, before)
	pr = pager.applyOrder(pr, last != nil)
	if limit := paginateLimit(first, last); limit != 0 {
		pr.Limit(limit)
	}
	if field := collectedField(ctx, edgesField, nodeField); field != nil {
		if err := pr.collectField(ctx, graphql.GetOperationContext(ctx), *field, []string{edgesField, nodeField}); err != nil {
			return nil, err
		}
	}

	nodes, err := pr.All(ctx)
	if err != nil {
		return nil, err
	}
	conn.build(nodes, pager, after, first, before, last)
	return conn, nil
}

// ProbabilityOrderField defines the ordering field of Probability.
type ProbabilityOrderField struct {
	field    string
	toCursor func(*Probability) Cursor
}

// ProbabilityOrder defines the ordering of Probability.
type ProbabilityOrder struct {
	Direction OrderDirection         `json:"direction"`
	Field     *ProbabilityOrderField `json:"field"`
}

// DefaultProbabilityOrder is the default ordering of Probability.
var DefaultProbabilityOrder = &ProbabilityOrder{
	Direction: OrderDirectionAsc,
	Field: &ProbabilityOrderField{
		field: probability.FieldID,
		toCursor: func(pr *Probability) Cursor {
			return Cursor{ID: pr.ID}
		},
	},
}

// ToEdge converts Probability into ProbabilityEdge.
func (pr *Probability) ToEdge(order *ProbabilityOrder) *ProbabilityEdge {
	if order == nil {
		order = DefaultProbabilityOrder
	}
	return &ProbabilityEdge{
		Node:   pr,
		Cursor: order.Field.toCursor(pr),
	}
}