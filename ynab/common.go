package ynab

type TransactionListParams struct {
	SinceDate             *Date
	Type                  *string
	LastKnowledgeOfServer *int64
}

type ListParams struct {
	LastKnowledgeOfServer *int64
}
