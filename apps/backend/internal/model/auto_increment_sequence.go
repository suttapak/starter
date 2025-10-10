package model

type (
	AutoIncrementSequence struct {
		CommonModel
		EntityType EntityType `db:"entity_type" json:"entity_type"` // "transaction_sale", "transaction_purchase", "product", "lot"
		TeamID     uint       `db:"team_id" json:"team_id"`         // for team-specific sequences
		EntityID   uint       `db:"entity_id" json:"entity_id"`     // additional identifier (e.g., product_id for lots)
		Sequence   uint       `db:"sequence" json:"sequence"`       // current sequence number
	}
	EntityType string
)

const (
	EntityTypeTransactionSale     EntityType = "SO"
	EntityTypeTransactionReturn   EntityType = "CN"
	EntityTypeTransactionPurchase EntityType = "PO"
	EntityTypeProduct             EntityType = "product"
	EntityTypeLot                 EntityType = "lot"
)
