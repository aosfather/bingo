package sql

import (
	"database/sql"
	//	"fmt"
	"log"
	"github.com/aosfather/bingo/utils"
)

type Page struct {
	Size  int
	Index int
	Count int
}

func (this *Page) getStart() int {
	if this.Index > 0 {
		return this.Size * (this.Index - 1)
	}
	return 0
}

func (this *Page) getEnd() int {
	if this.Index > 0 {
		return this.Size * this.Index
	}

	return this.Size - 1
}

type TxSession struct {
	tx   *sql.Tx
	db   *sql.DB
	isTx bool //是否开启了事务
}

func (this *TxSession) Begin() {
	if this.isTx {
		panic("tx has opened!")
	}
	this.tx, _ = this.db.Begin()
	this.isTx = true

}

func (this *TxSession) Commit() {
	if this.tx != nil && this.isTx {
		this.tx.Commit()
		this.isTx = false
	}
}

func (this *TxSession) Rollback() {
	if this.tx != nil && this.isTx {
		this.tx.Rollback()
		this.isTx = false
	}
}

func (this *TxSession) prepare(sql string) (*sql.Stmt, error) {
	if this.isTx {
		return this.tx.Prepare(sql)
	} else if this.db != nil {
		return this.db.Prepare(sql)
	}
	return nil, utils.CreateError(500, "no db init")
}

func (this *TxSession) Close() {
	this.Rollback()
}

func (this *TxSession) SimpleQuery(sql string, obj ...interface{}) bool {
	stmt, err := this.prepare(sql)
	if err != nil {
		log.Println(err)
		return false
	}
	defer stmt.Close()
	rs, err := stmt.Query()
	if err != nil {
		log.Println(err)
		return false
	}
	defer rs.Close()
	if rs.Next() {
		rs.Scan(obj...)
		return true
	}

	return false

}
func (this *TxSession) Find(obj interface{}, col ...string) bool {
	sql, args, err := CreateQuerySql(obj, col...)
	if err != nil {
		return false
	}
	log.Println(sql)
	return this.Query(obj, sql, args...)

}

func (this *TxSession) Insert(obj interface{}) (id int64, affect int64, err error) {
	sql, args, err := GetInsertSql(obj)
	if err != nil {
		return 0, 0, err
	}
	log.Println(sql)
	return this.ExeSql(sql, args...)
}

/**
 */
func (this *TxSession) Update(obj interface{}, col ...string) (id int64, affect int64, err error) {
	sql, args, err := CreateUpdateSql(obj, col...)
	if err != nil {
		return 0, 0, err
	}
	log.Println(sql)
	return this.ExeSql(sql, args...)
}

func (this *TxSession) Delete(obj interface{}, col ...string) (id int64, affect int64, err error) {
	sql, args, err := CreateDeleteSql(obj, col...)
	if err != nil {
		return 0, 0, err
	}
	log.Println(sql)
	return this.ExeSql(sql, args...)
}

func (this *TxSession) ExeSql(sql string, objs ...interface{}) (id int64, affect int64, err error) {
	stmt, err := this.prepare(sql)
	if err != nil {
		log.Println(err)
		return 0, 0, err
	}
	defer stmt.Close()

	rs, err := stmt.Exec(objs...)
	if err != nil {
		log.Println(err)
		return 0, 0, err
	}
	id, _ = rs.LastInsertId()
	affect, _ = rs.RowsAffected()
	return id, affect, nil

}

func (this *TxSession) QueryByPage(result interface{}, page Page, sql string, objs ...interface{}) []interface{} {
    //使用真分页的方式实现
	stmt, err := this.prepare(sql+buildMySqlLimitSql(page))
	if err != nil {
		log.Println(err)
		return nil
	}
	defer stmt.Close()
	rs, err := stmt.Query(objs...)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rs.Close()

	resultArray := []interface{}{}
	resultType := utils.GetRealType(result)
	//startIndex := page.getStart()
	//endIndex := page.getEnd()
	//var index int = 0
	cols, _ := rs.Columns()
	for {

		if rs.Next() {
			//if index < startIndex {
			//	index++
			//	continue
			//}
			//
			//if index >= endIndex {
			//	break
			//}
			columnsMap := make(map[string]interface{}, len(cols))
			refs := make([]interface{}, 0, len(cols))
			for _, col := range cols {
				var ref interface{}
				columnsMap[col] = &ref
				refs = append(refs, &ref)
			}

			rs.Scan(refs...)
			arrayItem := utils.CreateObjByType(resultType)
			//填充result
			utils.FillStruct(columnsMap, arrayItem)
			resultArray = append(resultArray, arrayItem)

			//index++

		} else {
			break
		}
	}
	return resultArray
}
func (this *TxSession) Query(result interface{}, sql string, objs ...interface{}) bool {
	stmt, err := this.prepare(sql)
	if err != nil {
		log.Println(err)
		return false
	}
	defer stmt.Close()
	rs, err := stmt.Query(objs...)
	if err != nil {
		log.Println(err)
		return false
	}
	defer rs.Close()

	cols, _ := rs.Columns()
	columnsMap := make(map[string]interface{}, len(cols))
	refs := make([]interface{}, 0, len(cols))
	for _, col := range cols {
		var ref interface{}
		columnsMap[col] = &ref
		refs = append(refs, &ref)
	}
	if rs.Next() {
		rs.Scan(refs...)
		//填充result
		utils.FillStruct(columnsMap, result)
		return true
	}

	return false
}

type SessionFactory struct {
	DBtype     string
	DBurl      string
	DBuser     string
	DBpassword string
	DBname     string
	pool       *sql.DB
}

func (this *SessionFactory) Init() {
	//如果已经初始化，不在初始化
	if this.pool != nil {
		return
	}

	if this.DBtype != "" {
		url := this.DBuser + ":" + this.DBpassword + "@" + this.DBurl + "/" + this.DBname
		var err error
		this.pool, err = sql.Open(this.DBtype, url)
		if err == nil {
			this.pool.Ping()
		} else {
			log.Printf("%v", err)
		}

	}
}

func (this *SessionFactory) GetSession() *TxSession {
	var session TxSession
	session.db = this.pool
	return &session
}
