package db

import (
	_ "github.com/lib/pq" // needed for postgres driver

)

//
// // Fetches multiple logs, the number (default 10) is configurable through the limit parameter
// func (db *BarkPostgresDb) FetchLimitedLogs(rowLimit int) ([]models.BarkLog, error) {
// 	var actualRowLimit int = 10 // default is 10, if its specified through params, the limit is overwritten
// 	if rowLimit > 0 {
// 		actualRowLimit = rowLimit
// 	}
// 	var barkLogs []models.BarkLog
//
// 	query := `
// 	SELECT id,log_time,log_level,service_name,code,msg,more_data
// 	FROM app_log LIMIT $1;`
//
// 	err := db.Client.Select(&barkLogs, query, actualRowLimit)
// 	if err != nil {
// 		return nil, fmt.Errorf("error fetching log rows: %w", err)
// 	}
//
// 	return barkLogs, nil
//
// }
