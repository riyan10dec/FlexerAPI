package query

type Query struct {
	key   string
	query string
}

var q = []Query{
	Query{
		key: "loginQuery",
		query: `CALL sp_front_login(
			?,?,?,?,?,?,?)`,
	}, Query{
		key:   "getClientQuery",
		query: "SELECT client_name FROM Master_Client where client_id = ? LIMIT 1",
	}, Query{
		key:   "frontCheckSession",
		query: "CALL sp_front_check_session (?)",
	}, Query{
		key:   "logoutQuery",
		query: "CALL sp_front_logout (?)",
	}, Query{
		key:   "getUserID",
		query: "SELECT user_id FROM Master_User where user_login = ? LIMIT 1",
	}, Query{
		key:   "createTransactionQuery",
		query: `CALL sp_front_add_transaction(?,?,?,?,?,?,?)`,
	}, Query{
		key:   "getScreenshotParamQuery",
		query: `CALL sp_front_add_screenshot(?,?,?,?)`,
	}, Query{
		key:   "reportScreenshotStatusQuery",
		query: `CALL sp_front_update_screenshot(?,?,?,?)`,
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
	}, Query{
		key:   "cmsAddUser",
		query: `call sp_add_user(?,?,?,?,?,?,?,?,?)`,
	}, Query{
		key:   "cmsEditUser",
		query: `call sp_edit_user(?,?,?,?,?,?,?,?,?)`,
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
