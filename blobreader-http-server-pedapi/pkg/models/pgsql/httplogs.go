package pgsql

import (
	"database/sql"
	"errors"

	"montgomery.wg/pedapilog/pkg/models"
)

// Define a PEDAPIRecordModel type which wraps a sql.DB connection pool.
type PEDAPIRecordModel struct {
    DB *sql.DB
}

// This will return a specific request based on its id.
func (m *PEDAPIRecordModel) Get(id int) (*models.PEDAPIRecord, error) {
    stmt := `SELECT requesttime, partnerid, merchantid, terminalid, url, responsestatus, requestbody, responsebody FROM httplogs
	WHERE id = $1`
	
	row := m.DB.QueryRow(stmt, id)

	req := &models.PEDAPIRecord{}

	err := row.Scan(&req.RequestTime, &req.PartnerID, &req.MerchantID, &req.TerminalID, &req.URL, &req.ResponseStatus, &req.RequestBody, &req.ResponseBody)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, models.ErrNoRecord
        } else {
             return nil, err
        }
	}
	
	return req, nil
}

// This will return the 100 most recently created requests.
func (m *PEDAPIRecordModel) Latest() ([]*models.PEDAPIRecord, error) {
    // Write the SQL statement we want to execute.
    stmt := `SELECT requesttime, partnerid, merchantid, terminalid, url, responsestatus, requestbody, responsebody FROM httplogs
    ORDER BY requesttime DESC LIMIT 100`

    rows, err := m.DB.Query(stmt)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    pedapirequests := []*models.PEDAPIRecord{}

    for rows.Next() {
 
        req := &models.PEDAPIRecord{}

        err = rows.Scan(&req.RequestTime, &req.PartnerID, &req.MerchantID, &req.TerminalID, &req.URL, &req.ResponseStatus, &req.RequestBody, &req.ResponseBody)
        if err != nil {
            return nil, err
        }

        // Append it to the slice of pedapirequests.
        pedapirequests = append(pedapirequests, req)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return pedapirequests, nil
}
