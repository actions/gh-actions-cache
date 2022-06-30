package types

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

var SORT_INPUT_TO_QUERY_MAP = map[string]string{
	"created-at": "created_at",
	"last-used":  "last_accessed_at",
	"size":       "size_in_bytes",
}

type BaseOptions struct {
	Repo   string
	Branch string
	Key    string
}

type ListOptions struct {
	BaseOptions
	Limit int
	Order string
	Sort  string
}

type DeleteOptions struct {
	BaseOptions
	Confirm bool
}

func (o *ListOptions) Validate() error {
	if o.Order != "" && o.Order != "asc" && o.Order != "desc" {
		return fmt.Errorf(fmt.Sprintf("%s is not a valid value for order flag. Allowed values: asc/desc", o.Order))
	}

	if o.Sort != "" && o.Sort != "last-used" && o.Sort != "size" && o.Sort != "created-at" {
		return fmt.Errorf(fmt.Sprintf("%s is not a valid value for sort flag. Allowed values: last-used/size/created-at", o.Sort))
	}

	if o.Limit < 1 || o.Limit > 100 {
		return fmt.Errorf(fmt.Sprintf("%d is not a valid value for limit flag. Allowed values: 1-100", o.Limit))
	}

	return nil
}

func (o *BaseOptions) GenerateBaseQueryParams(query url.Values) {
	if o.Branch != "" {
		if strings.HasPrefix(o.Branch, "refs/") {
			query.Add("ref", o.Branch)
		} else {
			query.Add("ref", fmt.Sprintf("refs/heads/%s", o.Branch))
		}
	}

	if o.Key != "" {
		query.Add("key", o.Key)
	}
}

func (o *ListOptions) GenerateQueryParams(query url.Values) {
	if o.Limit != 30 {
		query.Add("per_page", strconv.Itoa(o.Limit))
	}

	if o.Order != "" {
		query.Add("direction", o.Order)
	}

	if o.Sort != "" {
		query.Add("sort", SORT_INPUT_TO_QUERY_MAP[o.Sort])
	}

	o.GenerateBaseQueryParams(query)
}
