package ynab

import (
	"fmt"
	"net/url"
)

// TransactionListParams holds optional filter parameters for transaction list endpoints.
type TransactionListParams struct {
	SinceDate             *Date   // only return transactions on or after this date
	Type                  *string // filter by "uncategorized" or "unapproved"
	LastKnowledgeOfServer *int64  // for delta requests; pass the value returned by a prior call
}

// ListParams holds optional filter parameters for list endpoints that support delta requests.
type ListParams struct {
	LastKnowledgeOfServer *int64 // for delta requests; pass the value returned by a prior call
}

func buildListParams(params *ListParams) url.Values {
	q := url.Values{}
	if params != nil && params.LastKnowledgeOfServer != nil {
		q.Set("last_knowledge_of_server", fmt.Sprintf("%d", *params.LastKnowledgeOfServer))
	}
	return q
}
