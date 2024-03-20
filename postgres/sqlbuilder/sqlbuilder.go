package sqlbuilder

import (
	"fmt"
	"strings"
)

type M map[string]string

type SqlBuilder struct {
	strings.Builder
	args    []interface{}
	counter int
}

func NewSqlBuilder() *SqlBuilder {
	return &SqlBuilder{}
}

func (sb *SqlBuilder) ToSql() (string, []interface{}) {
	sqlQuery := sb.String()
	args := sb.args
	return sqlQuery, args
}

func (sb *SqlBuilder) Select(fields ...string) *SqlBuilder {
	sb.WriteString("SELECT ")
	sb.WriteString(strings.Join(fields, ", "))
	sb.WriteString(" ")
	return sb
}

func (sb *SqlBuilder) From(table string) *SqlBuilder {
	sb.WriteString("FROM ")
	sb.WriteString(table)
	sb.WriteString(" ")
	return sb
}

func (sb *SqlBuilder) LeftJoin(table, condition string) *SqlBuilder {
	sb.WriteString("LEFT JOIN ")
	sb.WriteString(table)
	sb.WriteString(" ON ")
	sb.WriteString(condition)
	sb.WriteString(" ")
	return sb
}

func (sb *SqlBuilder) Where() *SqlBuilder {
	sb.WriteString("WHERE ")
	return sb
}

func (sb *SqlBuilder) Eq(field string, value interface{}) *SqlBuilder {
	sb.WriteString(fmt.Sprintf("%s = $%d ", field, sb.counter+1))
	sb.args = append(sb.args, value)
	return sb
}

func (sb *SqlBuilder) EqAny(field string, value interface{}) *SqlBuilder {
	sb.WriteString(fmt.Sprintf("%s = ANY($%d) ", field, sb.counter+1))
	sb.args = append(sb.args, value)
	return sb
}

func (sb *SqlBuilder) Lt(field string, value interface{}) *SqlBuilder {
	sb.WriteString(fmt.Sprintf("%s < $%d ", field, sb.counter+1))
	sb.args = append(sb.args, value)
	return sb
}

func (sb *SqlBuilder) Gt(field string, value interface{}) *SqlBuilder {
	sb.WriteString(fmt.Sprintf("%s > $%d ", field, sb.counter+1))
	sb.args = append(sb.args, value)
	return sb
}

func (sb *SqlBuilder) And() *SqlBuilder {
	sb.WriteString("AND ")
	return sb
}

func (sb *SqlBuilder) OrderBy(orderBy map[string]string) *SqlBuilder {
	sb.WriteString("ORDER BY ")
	i := 0
	for column, direction := range orderBy {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(column)
		sb.WriteString(" ")
		sb.WriteString(direction)
		i++
	}
	sb.WriteString(" ")
	return sb
}

func (sb *SqlBuilder) Limit(limit int) *SqlBuilder {
	sb.WriteString(fmt.Sprintf("LIMIT %d ", limit))
	return sb
}

func (sb *SqlBuilder) Offset(offset int) *SqlBuilder {
	sb.WriteString(fmt.Sprintf("OFFSET %d ", offset))
	return sb
}
