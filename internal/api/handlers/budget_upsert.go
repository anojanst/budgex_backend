package handlers

import "gorm.io/gorm/clause"

// Updates amount on conflict (month+category_id+user_id uniqueness is implied by our query usage)
func clauseOnConflictSetAmount() clause.OnConflict {
	return clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "month"}, {Name: "category_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"amount": clause.Column{Name: "excluded.amount"}}),
	}
}
