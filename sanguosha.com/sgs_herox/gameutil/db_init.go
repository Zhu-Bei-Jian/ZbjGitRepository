package gameutil

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

var strProcAddColumn = `
CREATE PROCEDURE ProcAddColumn(TableName VARCHAR(50),ColumnName VARCHAR(50),SqlStr VARCHAR(4000))
BEGIN
	DECLARE Rows1 INT;
	SET Rows1=0;
	SELECT COUNT(*) INTO Rows1  FROM INFORMATION_SCHEMA.Columns
		WHERE table_schema= DATABASE() AND table_name=TableName AND column_name=ColumnName;
	-- 新增列
	IF (Rows1<=0) THEN
		SET SqlStr := CONCAT( 'ALTER TABLE ',TableName,' ADD COLUMN ',ColumnName,' ',SqlStr);
	ELSE
		SET SqlStr := '';
	END IF;
	-- 执行命令
	IF (SqlStr<>'') THEN 
		SET @SQL1 = SqlStr;
		PREPARE stmt1 FROM @SQL1;
		EXECUTE stmt1;
	END IF;
END
`

var strAddColumn = "CALL ProcAddColumn (?,?,?)"

func CheckProcAddColumn(db *sql.DB) error {
	_, err := db.Exec(strProcAddColumn)
	if err != nil {
		if err, ok := err.(*mysql.MySQLError); ok && err.Number == 1304 {
			return nil
		}
		logrus.Info("CheckProcAddColumn ", err.Error())
	}
	return nil
}

func AddColumnIfNotExist(db *sql.DB, tableName string, columnName string, columnInfo string) error {
	_, err := db.Exec(strAddColumn, tableName, columnName, columnInfo)
	return err
}
