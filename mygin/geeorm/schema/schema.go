package schema

import (
	"daydayup/mygin/geeorm/dialect"
	"go/ast"
	"reflect"
)

// Dialect 实现了一些特定的 SQL 语句的转换，接下来我们将要实现 ORM 框架中最为核心的转换——对象(object)和表(table)的转换。
// 给定一个任意的对象，转换为关系型数据库中的表结构。
// Field represents a column of database
type Field struct {
	Name string
	Type string
	Tag  string
}

// Schema represents a table of database
type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string
	fieldMap   map[string]*Field
}

// GetField returns field by name
func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

// 根据数据库中列的顺序，从对象中找到对应的值，按顺序平铺。即 u1、u2 转换为 ("Tom", 18), ("Same", 25) 这样的格式。
// 因此在实现 Insert 功能之前，还需要给 Schema 新增一个函数 RecordValues 完成上述的转换。
func (schema *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range schema.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}

type ITableName interface {
	TableName() string
}

// Parse a struct to a Schema instance
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	// typeOf() 和 ValueOf() 是 reflect 包最为基本也是最重要的 2 个方法，分别用来返回入参的类型和值。
	// 因为设计的入参是一个对象的指针，因此需要 reflect.Indirect() 获取指针指向的实例
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()

	//modelType.Name() 获取到结构体的名称作为表名。
	var tableName string
	t, ok := dest.(ITableName)
	if !ok {
		tableName = modelType.Name()
	} else {
		tableName = t.TableName()
	}
	schema := &Schema{
		Model:    dest,
		Name:     tableName,
		fieldMap: make(map[string]*Field),
	}

	//NumField() 获取实例的字段的个数，然后通过下标获取到特定字段 p := modelType.Field(i)。
	//p.Name 即字段名，p.Type 即字段类型，通过 (Dialect).DataTypeOf() 转换为数据库的字段类型，p.Tag 即额外的约束条件
	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("geeorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}
