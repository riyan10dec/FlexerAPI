package model

import (
	query "FlexerAPI/Query"
	"database/sql"
)

type Position struct {
	ClientID     int    `json:"clientID"`
	PositionName string `json:"positionName"`
}

func (p *Position) GetAllPositions(db *sql.DB) (error, []Position) {
	//rows, err := db.Query(query.SearchQuery("loginQuery"), l.UserLogin, l.Password)
	rows, err := db.Query(query.SearchQuery("cmsGetAllPositions"),
		p.ClientID)
	defer rows.Close()
	if err != nil {
		return err, nil
	}
	var ps []Position

	for rows.Next() {
		var p Position
		err := rows.Scan(&p.PositionName)
		if err != nil {
			return err, nil
		}
		ps = append(ps, p)
	}
	err = rows.Err()
	if err != nil {
		return err, nil
	}
	return nil, ps
}
