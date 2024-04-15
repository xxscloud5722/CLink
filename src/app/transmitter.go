package app

import (
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"log"
	"net/url"
	"reflect"
	"strings"
)

type Transmitter interface {
	// transmitter 发送器, 推送日志消息到目的地.
	transmitter([]*LogMessage) error
}

type ClickhouseTransmitter struct {
	Server     string
	Conn       *sqlx.DB
	Table      string
	Columns    []string
	ColumnsRow string
	Fields     map[string]string
}

func NewClickhouseTransmitter(config *Config) (*ClickhouseTransmitter, error) {
	clickhouseURL := fmt.Sprintf("tcp://%s:%d?username=%s&password=%s&database=%s",
		config.Clickhouse.Server, config.Clickhouse.Port,
		config.Clickhouse.Username, url.QueryEscape(config.Clickhouse.Password),
		config.Clickhouse.Database)
	clickhouseConn, err := sqlx.Connect("clickhouse", clickhouseURL)
	if err != nil {
		return nil, err
	}
	columns := lo.Map(reflect.ValueOf(config.Clickhouse.Fields).MapKeys(), func(item reflect.Value, index int) string {
		return item.String()
	})
	columnsRow := strings.Repeat("?, ", len(columns))
	columnsRow = columnsRow[:len(columnsRow)-2]
	return &ClickhouseTransmitter{
		Server:     config.Clickhouse.Server,
		Conn:       clickhouseConn,
		Table:      config.Clickhouse.Table,
		Columns:    columns,
		ColumnsRow: columnsRow,
		Fields:     config.Clickhouse.Fields,
	}, nil
}

func (c *ClickhouseTransmitter) transmitter(message []*LogMessage) error {
	sql, params := c.toSQL(message)
	tx, err := c.Conn.Begin()
	if err != nil {
		return err
	}
	prepare, err := tx.Prepare(sql)
	if err != nil {
		return err
	}
	for _, param := range params {
		_, flag := lo.Find(param, func(item any) bool {
			return item.(string) == ""
		})
		if flag {
			continue
		}
		_, err = prepare.Exec(param...)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	log.Println(fmt.Sprintf("[Clickhouse] Data => CK(%s) %d", c.Server, len(params)))
	if err != nil {
		return err
	}
	return nil
}

// toSQL 将数据转成脚本以及参数.
func (c *ClickhouseTransmitter) toSQL(message []*LogMessage) (string, [][]any) {
	var sql = "INSERT INTO " + c.Table + " ( " + strings.Join(c.Columns, ", ") + " ) \n " +
		"VALUES \n ( " + c.ColumnsRow + " )"
	//for i := range message {
	//	if i+1 == len(message) {
	//		sql += " ( " + c.ColumnsRow + " )"
	//	} else {
	//		sql += " ( " + c.ColumnsRow + " ), \n"
	//	}
	//}
	var params [][]any
	for _, item := range message {
		var values []any
		for _, column := range c.Columns {
			dataKey := c.Fields[column]
			value, dataOk := item.Attribute[dataKey]
			if dataOk {
				values = append(values, value)
			} else {
				values = append(values, "")
			}
		}
		params = append(params, values)
	}
	return sql, params
}
