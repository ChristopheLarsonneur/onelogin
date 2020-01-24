package api

import (
	"strconv"
	"strings"
	"time"
)

// QueryOptions define Get Users, Get Roles, Get Events, and Get Groups resource APIs
// see https://developers.onelogin.com/api-docs/1/getting-started/using-query-parameters
type QueryOptions struct {
	fields []string
	search map[string]string
	limit  *int
	sort   string
	since  *time.Time
	until  *time.Time
}

// NewQueryOptions Create a QueryOptions to use in few Get functions
func NewQueryOptions() (ret *QueryOptions) {
	ret = new(QueryOptions)
	return
}

// SetFields initialize the list of fields to extract.
func (q *QueryOptions) SetFields(fields ...string) (ret *QueryOptions) {
	ret = q
	if ret == nil {
		return q
	}
	q.fields = fields
	return
}

// SetLimit define the query limit. Max is 50 as defined  by the API documentation
// If > 50, limit value is ignored.
func (q *QueryOptions) SetLimit(limit int) (ret *QueryOptions) {
	ret = q
	if ret == nil || limit > 50 {
		return q
	}
	q.limit = &limit
	return
}

// AddFilterOn add a field filter.
// Values support wildcards, as described in
// https://developers.onelogin.com/api-docs/1/getting-started/using-query-parameters#search
//
// Example:
// field:email, value:test@onelogin.com
// field:email, value:*@onelogin.com
// field:email, value:!@onelogin.com
func (q *QueryOptions) AddFilterOn(field, value string) (ret *QueryOptions) {
	ret = q
	if ret == nil {
		return q
	}
	if q.search == nil {
		q.search = make(map[string]string)
	}
	q.search[field] = value
	return
}

// Sort indicates which field is goig to be used to sort the result
// Do not speficy + or - as described in the documentation
// this function wil add proper sort order based on ascending parameter.
func (q *QueryOptions) Sort(field string, ascending bool) (ret *QueryOptions) {
	ret = q
	if ret == nil {
		return q
	}
	if ascending {
		q.sort = "+" + field
	} else {
		q.sort = "-" + field
	}
	return
}

// Since filter on resource 'created_at' field of get users and get events
// See https://developers.onelogin.com/api-docs/1/getting-started/using-query-parameters#sinceanduntil
func (q *QueryOptions) Since(since *time.Time) (ret *QueryOptions) {
	ret = q
	if ret == nil {
		return q
	}
	q.since = since
	return
}

// Until filter on resource 'created_at' field of get users and get events
// See https://developers.onelogin.com/api-docs/1/getting-started/using-query-parameters#sinceanduntil
func (q *QueryOptions) Until(until *time.Time) (ret *QueryOptions) {
	ret = q
	if ret == nil {
		return q
	}
	q.until = until
	return
}

func (q *QueryOptions) getQueryParameters() (ret map[string]string) {
	if q == nil {
		return nil
	}

	ret = make(map[string]string)

	if q.fields != nil {
		ret["fields"] = strings.Join(q.fields, ",")
	}

	if q.search != nil {
		for field, value := range q.search {
			ret[field] = value
		}
	}

	if q.limit != nil {
		ret["limit"] = strconv.Itoa(*q.limit)
	}

	if q.sort != "" {
		ret["sort"] = q.sort
	}

	if q.since != nil {
		ret["since"] = q.since.Format(time.RFC3339)
	}

	if q.until != nil {
		ret["until"] = q.until.Format(time.RFC3339)
	}

	return
}
