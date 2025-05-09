package message

import "insider/types"

type Filter struct {
	types.PagedFilter
	Query *string
}
