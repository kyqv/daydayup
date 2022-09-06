package session

import "daydayup/mygin/geeorm/log"

// GeeORM 之前的操作均是执行完即自动提交的，每个操作是相互独立的。
// 之前直接使用 sql.DB 对象执行 SQL 语句，如果要支持事务，
// 需要更改为 sql.Tx 执行。在 Session 结构体中新增成员变量 tx *sql.Tx，
// 当 tx 不为空时，则使用 tx 执行 SQL 语句，否则使用 db 执行 SQL 语句。
// 这样既兼容了原有的执行方式，又提供了对事务的支持

// 封装事务的 Begin、Commit 和 Rollback 三个接口
// Begin a transaction
func (s *Session) Begin() (err error) {
	log.Info("transaction begin")
	if s.tx, err = s.db.Begin(); err != nil {
		log.Error(err)
		return
	}
	return
}

// Commit a transaction
func (s *Session) Commit() (err error) {
	log.Info("transaction commit")
	if err = s.tx.Commit(); err != nil {
		log.Error(err)
	}
	return
}

// Rollback a transaction
func (s *Session) Rollback() (err error) {
	log.Info("transaction rollback")
	if err = s.tx.Rollback(); err != nil {
		log.Error(err)
	}
	return
}
