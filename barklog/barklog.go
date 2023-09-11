package barklog

import (
	"encoding/json"
	"time"
)

type PostgresJsonbField map[string]interface{}

// Struct representing a log in Bark
type BarkLog struct {
	Id          int64           `db:"id"`
	LogTime     time.Time       `db:"log_time"`
	LogLevel    int             `db:"log_level"`
	ServiceName string          `db:"service_name"`
	Code        string          `db:"code"`
	Message     string          `db:"msg"`
	MoreData    json.RawMessage `db:"more_data"`
}

// func (postgresJsonbField *PostgresJsonbField) Value() (driver.Value, error) {
// 	return json.Marshal(postgresJsonbField)
// }

// func (pc *PostgresJsonbField) Scan(val interface{}) error {
// 	switch v := val.(type) {
// 	case []byte:
// 		json.Unmarshal(v, &pc)
// 		return nil
// 	case string:
// 		json.Unmarshal([]byte(v), &pc)
// 		return nil
// 	default:
// 		return errors.New(fmt.Sprintf("Unsupported type: %T", v))
// 	}
// }
