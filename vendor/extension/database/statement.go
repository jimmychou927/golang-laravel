package database

import (
	"errors"
	"strconv"
	"strings"
	"sync"
)

type Where struct {
	operation string
	field     string
	qmark     string
}

type Join struct {
	method    string
	table     string
	alias     string
	fieldA    string
	operation string
	fieldB    string
	args      []interface{}
}

type RawUpdate struct {
	expression string
	args       []interface{}
}

type Sql struct {
	fields    []string
	fieldsRaw string
	table     string
	wheres    []Where
	joins     []Join
	argsJoin  []interface{}
	args      []interface{}
	order     []string
	group     string
	offset    string
	limit     string
	whereRaw  string
	updateRaw []RawUpdate
	statement string
}

var SqlPool = sync.Pool{
	New: func() interface{} {
		return &Sql{
			fields:    make([]string, 0),
			fieldsRaw: "",
			table:     "",
			wheres:    make([]Where, 0),
			joins:     make([]Join, 0),
			argsJoin:  make([]interface{}, 0),
			args:      make([]interface{}, 0),
			order:     make([]string, 0),
			group:     "",
			offset:    "",
			limit:     "",
			whereRaw:  "",
			updateRaw: make([]RawUpdate, 0),
		}
	},
}

type H map[string]interface{}

func newSql() *Sql {
	return SqlPool.Get().(*Sql)
}

// *******************************
// process method
// *******************************

func Table(table string) *Sql {
	sql := newSql()
	sql.table = table
	return sql
}

func (sql *Sql) Select(fields ...string) *Sql {
	sql.fields = fields
	return sql
}

func (sql *Sql) SelectRaw(fields string) *Sql {
	if fields != "" {
		sql.fieldsRaw = ", " + fields
	}
	return sql
}

func (sql *Sql) OrderBy(filed string, order string) *Sql {
	if filed == "" {
		return sql
	}
	sql.order = append(sql.order, filed+" "+order)
	return sql
}

func (sql *Sql) GroupBy(fileds ...string) *Sql {
	sql.group = strings.Join(fileds, ",")
	return sql
}

func (sql *Sql) Skip(offset int) *Sql {
	sql.offset = strconv.Itoa(offset)
	return sql
}

func (sql *Sql) Take(take int) *Sql {
	sql.limit = strconv.Itoa(take)
	return sql
}

func (sql *Sql) convertWhereArgument(arg interface{}) ([]interface{}, string) {
	conv := make([]interface{}, 0)
	qmark := ""

	switch arg.(type) {
	case []int:
		if len(arg.([]int)) < 1 {
			return nil, ""
		}
		qmark = "(" + strings.Repeat("?, ", len(arg.([]int))-1) + "?)"
		for _, arrVar := range arg.([]int) {
			conv = append(conv, arrVar)
		}
	case []string:
		if len(arg.([]string)) < 1 {
			return nil, ""
		}
		qmark = "(" + strings.Repeat("?, ", len(arg.([]string))-1) + "?)"
		for _, arrVar := range arg.([]string) {
			conv = append(conv, arrVar)
		}
	case *Sql:
		qmark = "(" + arg.(*Sql).ToSQL() + ")"
		conv = arg.(*Sql).args
	default:
		return nil, ""
	}

	return conv, qmark
}

func (sql *Sql) Where(field string, args ...interface{}) *Sql {
	if len(args) == 1 {
		return sql.addWhereConditions(field, "=", []interface{}{args[0]}, "?")
	}

	return sql.addWhereConditions(field, args[0].(string), []interface{}{args[1]}, "?")
}

func (sql *Sql) WhereIn(field string, arg interface{}) *Sql {
	args, qmark := sql.convertWhereArgument(arg)
	return sql.addWhereConditions(field, "in", args, qmark)
}

func (sql *Sql) WhereNotIn(field string, arg interface{}) *Sql {
	args, qmark := sql.convertWhereArgument(arg)
	return sql.addWhereConditions(field, "not in", args, qmark)
}

func (sql *Sql) addWhereConditions(field string, operation string, args []interface{}, qmark string) *Sql {
	if field == "" || args == nil {
		return sql
	}
	sql.wheres = append(sql.wheres, Where{
		field:     field,
		operation: operation,
		qmark:     qmark,
	})
	if len(args) > 0 {
		sql.args = append(sql.args, args...)
	}
	return sql
}

func (sql *Sql) Find(arg interface{}) (map[string]interface{}, error) {
	return sql.Where("id", "=", arg).First()
}

func (sql *Sql) Count() (int64, error) {
	var (
		res map[string]interface{}
		err error
	)
	if res, err = sql.SelectRaw("count(*) as total").First(); err != nil {
		return 0, err
	}
	return res["total"].(int64), nil
}

func (sql *Sql) WhereRaw(raw string, args ...interface{}) *Sql {
	if raw == "" {
		return sql
	}
	sql.whereRaw = raw
	sql.args = append(sql.args, args...)
	return sql
}

func (sql *Sql) UpdateRaw(raw string, args ...interface{}) *Sql {
	sql.updateRaw = append(sql.updateRaw, RawUpdate{
		expression: raw,
		args:       args,
	})
	return sql
}

func (sql *Sql) convertJoinArgument(tbl interface{}, args ...string) (string, string, string, string, string, []interface{}) {
	table := ""
	alias := ""
	fieldA := ""
	operation := ""
	fieldB := ""

	if len(args) == 3 {
		fieldA = args[0]
		operation = args[1]
		fieldB = args[2]
	} else if len(args) == 4 {
		alias = args[0]
		fieldA = args[1]
		operation = args[2]
		fieldB = args[3]
	} else {
		return "", "", "", "", "", nil
	}

	switch tbl.(type) {
	case string:
		table = tbl.(string)
		return table, alias, fieldA, operation, fieldB, []interface{}{}
	case *Sql:
		table = "(" + tbl.(*Sql).ToSQL() + ")"
		return table, alias, fieldA, operation, fieldB, tbl.(*Sql).args
	}

	return "", "", "", "", "", nil
}

func (sql *Sql) Join(tbl interface{}, args ...string) *Sql {
	jointTable, JoinAlias, fieldA, operation, fieldB, joinArgs := sql.convertJoinArgument(tbl, args...)
	return sql.addJoinConditions("inner", jointTable, JoinAlias, fieldA, operation, fieldB, joinArgs)
}

func (sql *Sql) LeftJoin(tbl interface{}, args ...string) *Sql {
	jointTable, JoinAlias, fieldA, operation, fieldB, joinArgs := sql.convertJoinArgument(tbl, args...)
	return sql.addJoinConditions("left", jointTable, JoinAlias, fieldA, operation, fieldB, joinArgs)
}

func (sql *Sql) RightJoin(tbl interface{}, args ...string) *Sql {
	jointTable, JoinAlias, fieldA, operation, fieldB, joinArgs := sql.convertJoinArgument(tbl, args...)
	return sql.addJoinConditions("right", jointTable, JoinAlias, fieldA, operation, fieldB, joinArgs)
}

func (sql *Sql) addJoinConditions(method string, table string, alias string, fieldA string, operation string, fieldB string, args []interface{}) *Sql {
	if table == "" {
		return sql
	}
	sql.joins = append(sql.joins, Join{
		method:    method,
		fieldA:    fieldA,
		fieldB:    fieldB,
		table:     table,
		alias:     alias,
		operation: operation,
		args:      args,
	})
	return sql
}

// *******************************
// terminal method
// -------------------------------
// sql args order:
// update ... => where ...
// *******************************

func (sql *Sql) First() (map[string]interface{}, error) {
	defer RecycleSql(sql)

	sql.Take(1)
	res, _ := sql.All()

	if len(res) < 1 {
		return nil, errors.New("out of index")
	}
	return res[0], nil
}

func (sql *Sql) All() ([]map[string]interface{}, error) {
	defer RecycleSql(sql)

	sql.statement = sql.ToSQL()
	args := sql.argsJoin
	args = append(args, sql.args...)

	res, _ := GetConnection().Query(sql.statement, args...)

	return res, nil
}

func (sql *Sql) ToSQL() string {
	return "select " + sql.getFields() + " from " + sql.table +
		sql.getJoins() +
		sql.getWheres() +
		sql.getGroupBy() +
		sql.getOrderBy() +
		sql.getLimit() +
		sql.getOffset()
}

func (sql *Sql) Update(values H) (int64, error) {
	defer RecycleSql(sql)

	sql.prepareUpdate(values)

	res := GetConnection().Exec(sql.statement, sql.args...)

	if affectRow, _ := res.RowsAffected(); affectRow < 1 {
		return 0, errors.New("no affect row")
	}

	return res.LastInsertId()
}

func (sql *Sql) Delete() error {
	defer RecycleSql(sql)

	sql.statement = "delete from " + sql.table + sql.getWheres()

	res := GetConnection().Exec(sql.statement, sql.args...)

	if affectRow, _ := res.RowsAffected(); affectRow < 1 {
		return errors.New("no affect row")
	}

	return nil
}

func (sql *Sql) Exec() (int64, error) {
	defer RecycleSql(sql)

	sql.prepareUpdate(H{})

	res := GetConnection().Exec(sql.statement, sql.args...)

	if affectRow, _ := res.RowsAffected(); affectRow < 1 {
		return 0, errors.New("no affect row")
	}

	return res.LastInsertId()
}

func (sql *Sql) Insert(values H) (int64, error) {
	defer RecycleSql(sql)

	sql.prepareInsert(values)

	res := GetConnection().Exec(sql.statement, sql.args...)

	if affectRow, _ := res.RowsAffected(); affectRow < 1 {
		return 0, errors.New("no affect row")
	}

	return res.LastInsertId()
}

// *******************************
// internal help function
// *******************************

func (sql *Sql) getLimit() string {
	if sql.limit == "" {
		return ""
	}
	return "\nlimit " + sql.limit + " "
}

func (sql *Sql) getOffset() string {
	if sql.offset == "" {
		return ""
	}
	return "\noffset " + sql.offset + " "
}

func (sql *Sql) getOrderBy() string {
	if len(sql.order) < 1 {
		return ""
	}
	return "\norder by " + strings.Join(sql.order, ", ")
}

func (sql *Sql) getGroupBy() string {
	if sql.group == "" {
		return ""
	}
	return "\ngroup by " + sql.group + " "
}

func (sql *Sql) getJoins() string {
	joins := ""
	if len(sql.joins) > 0 {
		for _, join := range sql.joins {
			alias := join.alias
			if alias != "" {
				alias = " as " + alias + " "
			}

			joins += "\n" + join.method + " join " + join.table + alias + " on " + join.fieldA + " " + join.operation + " " + join.fieldB

			if len(join.args) > 0 {
				sql.argsJoin = append(sql.argsJoin, join.args...)
			}
		}
	}
	return joins
}

func (sql *Sql) getFields() string {
	if len(sql.fields) == 0 {
		if sql.fieldsRaw == "" {
			return "*"
		}

		return sql.fieldsRaw[1:]
	}
	if sql.fields[0] == "count(*)" {
		return "count(*)" + sql.fieldsRaw
	}
	fields := ""
	for _, field := range sql.fields {
		fields += sql.quoteField(field) + ","
	}

	return fields[:len(fields)-1] + sql.fieldsRaw
}

func (sql *Sql) getWheres() string {
	if len(sql.wheres) == 0 {
		if sql.whereRaw != "" {
			return "\nwhere " + sql.whereRaw
		}
		return ""
	}
	wheres := "\nwhere "
	for _, where := range sql.wheres {
		wheres += sql.quoteField(where.field) + " " + where.operation + " " + where.qmark + " and "
	}

	if sql.whereRaw != "" {
		return wheres + sql.whereRaw
	} else {
		return wheres[:len(wheres)-5]
	}
}

func (sql *Sql) addQuoteChar(fieldName string) string {
	quoteStr := fieldName
	if !strings.HasPrefix(quoteStr, "`") && !strings.HasPrefix(quoteStr, " ") {
		quoteStr = "`" + quoteStr
	}
	if !strings.HasSuffix(quoteStr, "`") && !strings.HasSuffix(quoteStr, " ") {
		quoteStr = quoteStr + "`"
	}
	return quoteStr
}

func (sql *Sql) quoteField(fieldStr string) string {
	quoteField := ""
	arr := strings.Split(strings.TrimSpace(fieldStr), ".")
	if len(arr) > 1 {
		quoteField += sql.addQuoteChar(arr[0]) + "." + sql.addQuoteChar(arr[1])
	} else {
		quoteField += sql.addQuoteChar(fieldStr)
	}
	return quoteField
}

func (sql *Sql) prepareUpdate(values H) {
	fields := ""
	args := make([]interface{}, 0)

	if len(values) != 0 {

		for key, value := range values {
			fields += sql.quoteField(key) + " =?, "
			args = append(args, value)
		}

		if len(sql.updateRaw) == 0 {
			fields = fields[:len(fields)-2]
		} else {
			for i := 0; i < len(sql.updateRaw); i++ {
				if i == len(sql.updateRaw)-1 {
					fields += sql.updateRaw[i].expression + " "
				} else {
					fields += sql.updateRaw[i].expression + ","
				}
				args = append(args, sql.updateRaw[i].args...)
			}
		}

		sql.args = append(args, sql.args...)
	} else {
		if len(sql.updateRaw) == 0 {
			panic("prepareUpdate: wrong parameter")
		} else {
			for i := 0; i < len(sql.updateRaw); i++ {
				if i == len(sql.updateRaw)-1 {
					fields += sql.updateRaw[i].expression + " "
				} else {
					fields += sql.updateRaw[i].expression + ","
				}
				args = append(args, sql.updateRaw[i].args...)
			}
		}
		sql.args = append(args, sql.args...)
	}

	sql.statement = "update " + sql.table + " set " + fields + sql.getWheres()
}

func (sql *Sql) prepareInsert(values H) {
	fields := make([]string, 0)
	quesMark := make([]string, 0)

	for key, value := range values {
		fields = append(fields, sql.quoteField(key))
		quesMark = append(quesMark, "?")
		sql.args = append(sql.args, value)
	}

	sql.statement = "insert into " + sql.table +
		"(" + strings.Join(fields, ", ") + ") values " +
		"(" + strings.Join(quesMark, ", ") + ")"
}

func (sql *Sql) empty() *Sql {
	sql.fields = make([]string, 0)
	sql.fieldsRaw = ""
	sql.wheres = make([]Where, 0)
	sql.joins = make([]Join, 0)
	sql.argsJoin = make([]interface{}, 0)
	sql.args = make([]interface{}, 0)
	sql.order = make([]string, 0)
	sql.group = ""
	sql.offset = ""
	sql.limit = ""
	sql.whereRaw = ""
	sql.updateRaw = make([]RawUpdate, 0)

	return sql
}

func RecycleSql(sql *Sql) {
	sql.fields = make([]string, 0)
	sql.fieldsRaw = ""
	sql.table = ""
	sql.wheres = make([]Where, 0)
	sql.joins = make([]Join, 0)
	sql.argsJoin = make([]interface{}, 0)
	sql.args = make([]interface{}, 0)
	sql.order = make([]string, 0)
	sql.group = ""
	sql.offset = ""
	sql.limit = ""
	sql.whereRaw = ""
	sql.updateRaw = make([]RawUpdate, 0)
	sql.statement = ""

	SqlPool.Put(sql)
}
