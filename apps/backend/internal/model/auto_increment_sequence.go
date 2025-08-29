package model

type (
	AutoIncrementSequence struct {
		CommonModel
		EntityType EntityType `json:"entity_type" gorm:"uniqueIndex:idx_sequence_entity_team"` // "transaction_sale", "transaction_purchase", "product", "lot"
		TeamID     uint       `json:"team_id" gorm:"uniqueIndex:idx_sequence_entity_team"`     // for team-specific sequences
		EntityID   uint       `json:"entity_id" gorm:"uniqueIndex:idx_sequence_entity_team"`   // additional identifier (e.g., product_id for lots)
		Sequence   uint       `json:"sequence" gorm:"default:0"`                               // current sequence number
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
