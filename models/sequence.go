package models

type Sequence struct {
	TableSchema string `orm:"column(table_schema);size(48)" description:"스키마"`
	TableName   string `orm:"column(table_name);size(48)" description:"테이블"`
	NextSeq     int64  `orm:"column(next_seq)" description:"순번"`
}
