package clause

import (
	"strings"
)

// 如果想一次构造出完整的 SQL 语句是比较困难的，因此我们将构造 SQL 语句这一部分独立出来，放在子package clause 中实现。

// Clause 拼接各个独立的子句
type Clause struct {
	sql     map[Type]string
	sqlVars map[Type][]interface{}
}

// Type is the type of Clause
type Type int

// Support types for Clause
const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
	COUNT
)

// Set 方法根据 Type 调用对应的 generator，生成该子句对应的 SQL 语句
func (c *Clause) Set(name Type, vars ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	sql, vars := generators[name](vars...)
	c.sql[name] = sql
	c.sqlVars[name] = vars
}

// Build 方法根据传入的 Type 的顺序，构造出最终的 SQL 语句。
func (c *Clause) Build(orders ...Type) (string, []interface{}) {
	var sqls []string
	var vars []interface{}
	for _, order := range orders {
		if sql, ok := c.sql[order]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.sqlVars[order]...)
		}
	}
	return strings.Join(sqls, " "), vars
}
