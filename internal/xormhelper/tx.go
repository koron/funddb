package xormhelper

import (
	"errors"

	"xorm.io/xorm"
)

// Tx starts a transaction and manages its lifecycle with calling a function.
func Tx(engine *xorm.Engine, fn func(*xorm.Session) error) error {
	session := engine.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}
	err := fn(session)
	if err != nil {
		if err2 := session.Rollback(); err2 != nil {
			err = errors.Join(err, err2)
		}
		return err
	}
	return session.Commit()
}
