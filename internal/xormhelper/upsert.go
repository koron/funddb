package xormhelper

import (
	"fmt"

	"xorm.io/xorm"
)

// UpsertOne inserts/updates exactly a row.
func UpsertOne(session *xorm.Session, id, bean any) error {
	updated, err := session.ID(id).Update(bean)
	if err != nil {
		return err
	}
	if updated == 1 {
		return nil
	}
	if updated != 0 {
		return fmt.Errorf("expected 1 row to update, but %d rows updated", updated)
	}
	inserted, err := session.Insert(bean)
	if err != nil {
		return err
	}
	if inserted != 1 {
		return fmt.Errorf("expected 1 row to insert, but %d rows inserted", inserted)
	}
	return nil
}
