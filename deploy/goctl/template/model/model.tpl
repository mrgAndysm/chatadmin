package {{.pkg}}
{{if .withCache}}
import (
    "context"
    "strings"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
    "github.com/Masterminds/squirrel"
    "xiaobien-rpc/common/globalkey"
)
{{else}}
import (
    "context"
    "strings"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
    "github.com/Masterminds/squirrel"
    "xiaobien-rpc/common/globalkey"
)
{{end}}
var _ {{.upperStartCamelObject}}Model = (*custom{{.upperStartCamelObject}}Model)(nil)

type (
	// {{.upperStartCamelObject}}Model is an interface to be customized, add more methods here,
	// and implement the added methods in custom{{.upperStartCamelObject}}Model.
	{{.upperStartCamelObject}}Model interface {
		{{.lowerStartCamelObject}}Model
        Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
        RowBuilder(row ...string) squirrel.SelectBuilder
        UpdateBuilder() squirrel.UpdateBuilder
        CountBuilder(field string) squirrel.SelectBuilder
        SumBuilder(field string) squirrel.SelectBuilder
        FindCount(ctx context.Context, countBuilder squirrel.SelectBuilder) (int64, error)
        FindSum(ctx context.Context,sumBuilder squirrel.SelectBuilder) (float64,error)
        FindAll(ctx context.Context, rowBuilder squirrel.SelectBuilder) ([]*{{.upperStartCamelObject}}, error)
        UpdateWithBuilder(ctx context.Context, session sqlx.Session, builder squirrel.UpdateBuilder) (sql.Result, error)
        DeleteSoftById(ctx context.Context, session sqlx.Session, id int64) (sql.Result, error)
        IncrColumnById(ctx context.Context, session sqlx.Session, dirId int64, column string, num int64) (sql.Result, error)
        DecrColumnById(ctx context.Context, session sqlx.Session, dirId int64, column string, num int64) (sql.Result, error)
	}

	custom{{.upperStartCamelObject}}Model struct {
		*default{{.upperStartCamelObject}}Model
	}
)

// New{{.upperStartCamelObject}}Model returns a model for the database table.
func New{{.upperStartCamelObject}}Model(conn sqlx.SqlConn{{if .withCache}}, c cache.CacheConf{{end}}) {{.upperStartCamelObject}}Model {
	return &custom{{.upperStartCamelObject}}Model{
		default{{.upperStartCamelObject}}Model: new{{.upperStartCamelObject}}Model(conn{{if .withCache}}, c{{end}}),
	}
}

func (m *default{{.upperStartCamelObject}}Model) Trans(ctx context.Context,fn func(ctx context.Context,session sqlx.Session) error) error {
	{{if .withCache}}
	return m.TransactCtx(ctx,func(ctx context.Context,session sqlx.Session) error {
		return  fn(ctx,session)
	})
	{{else}}
	return m.conn.TransactCtx(ctx,func(ctx context.Context,session sqlx.Session) error {
		return  fn(ctx,session)
	})
	{{end}}
}

func (m *default{{.upperStartCamelObject}}Model) RowBuilder(row ...string) squirrel.SelectBuilder {
	rows := {{.lowerStartCamelObject}}Rows
	if len(row) != 0 {
		rows = strings.Join(row, ",")
	}
	return squirrel.Select(rows).From(m.table)
}

func (m *default{{.upperStartCamelObject}}Model) UpdateBuilder() squirrel.UpdateBuilder {
	return squirrel.Update(m.table)
}

func (m *default{{.upperStartCamelObject}}Model) CountBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("COUNT(" + field + ")").From(m.table)
}

func (m *default{{.upperStartCamelObject}}Model) SumBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("IFNULL(SUM("+field+"),0)").From(m.table)
}

func (m *default{{.upperStartCamelObject}}Model) FindCount(ctx context.Context, countBuilder squirrel.SelectBuilder) (int64, error) {
	query, val, err := countBuilder.ToSql()
	if err != nil {
		return 0, err
	}
	var total int64
	{{if .withCache}}
	err = m.QueryRowNoCacheCtx(ctx,&resp, query, values...)
	{{else}}
	err = m.conn.QueryRowCtx(ctx, &total, query, val...)
	{{end}}
	return total, err
}

func (m *default{{.upperStartCamelObject}}Model) FindSum(ctx context.Context,sumBuilder squirrel.SelectBuilder) (float64,error) {

	query, values, err := sumBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	var resp float64
	{{if .withCache}}
	err = m.QueryRowNoCacheCtx(ctx,&resp, query, values...)
	{{else}}
	err = m.conn.QueryRowCtx(ctx,&resp, query, values...)
	{{end}}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *default{{.upperStartCamelObject}}Model) FindAll(ctx context.Context, rowBuilder squirrel.SelectBuilder) ([]*{{.upperStartCamelObject}}, error) {
	query, values, err := rowBuilder.ToSql()
	if err != nil {
		return nil, err
	}
	var resp []*{{.upperStartCamelObject}}
	{{if .withCache}}
	err = m.QueryRowsNoCacheCtx(ctx,&resp, query, values...)
	{{else}}
    err = m.conn.QueryRowsCtx(ctx,&resp, query, values...)
    {{end}}
	return resp, err
}

func (m *default{{.upperStartCamelObject}}Model) UpdateWithBuilder(ctx context.Context, session sqlx.Session, builder squirrel.UpdateBuilder) (sql.Result, error) {
	query, params, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "build sql fail", err)
	}
	if session != nil{
        return session.ExecCtx(ctx, query, params...)
    }else{
    	return m.conn.ExecCtx(ctx, query, params...)
    }
}

func (m *default{{.upperStartCamelObject}}Model) DeleteSoftById(ctx context.Context, session sqlx.Session, id int64) (sql.Result, error) {
	builder := m.UpdateBuilder().Set("is_deleted", globalkey.DeleteStatusYes).Where("id = ?", id)
	return m.UpdateWithBuilder(ctx, session, builder)
}

func (m *default{{.upperStartCamelObject}}Model) IncrColumnById(ctx context.Context, session sqlx.Session, dirId int64, column string, num int64) (sql.Result, error) {
	builder := m.UpdateBuilder().Set(column, squirrel.Expr(column+"+?", num)).Where("id=?", dirId)
	return m.UpdateWithBuilder(ctx, session, builder)
}

func (m *default{{.upperStartCamelObject}}Model) DecrColumnById(ctx context.Context, session sqlx.Session, dirId int64, column string, num int64) (sql.Result, error) {
	builder := m.UpdateBuilder().Set(column, squirrel.Expr(column+"-?", num)).Where("id=?", dirId)
	return m.UpdateWithBuilder(ctx, session, builder)
}