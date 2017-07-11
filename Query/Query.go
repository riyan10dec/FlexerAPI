package query

type Query struct {
	key   string
	query string
}

var q = []Query{
	Query{
		key:   "loginQuery",
		query: "SELECT user_name, client_id FROM Master_User where email = ? and password = ? LIMIT 1",
	}, Query{
		key:   "getClientQuery",
		query: "SELECT client_name FROM Master_Client where client_id = ? LIMIT 1",
	}, Query{
		key:   "getUserID",
		query: "SELECT user_id FROM Master_User where user_login = ? LIMIT 1",
	}, Query{
		key: "createTransactionQuery",
		query: `INSERT INTO Trx_transaction  (user_id, Date, Application_id, URL, Keystroke, Mouseclick, StartTime, EndTime)
    								Values(?, CURDATE(), ?, ?, ?, ?, ?, ?)`,
	}, Query{
		key: "getApplicationName",
		query: `SELECT 
				Application_name
			FROM
				Master_Application
			WHERE
				application_id = ?`,
	}, Query{
		key: "getApplicationID",
		query: `SELECT 
				Application_id
			FROM
				Master_Application
			WHERE
				application_name = ?`,
	}, Query{
		key: "checkApplicationName",
		query: `SELECT 
				1
			FROM
				Master_Application
			WHERE
				Application_name = ?`,
	}, Query{
		key:   "createApplicationQuery",
		query: `Insert into Master_Application (application_name) values(?)`,
	},
}

func SearchQuery(key string) string {
	for _, obj := range q {
		if obj.key == key {
			return obj.query
		}

	}
	return ""
}
