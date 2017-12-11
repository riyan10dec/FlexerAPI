package query

type Query struct {
	key   string
	query string
}

var q = []Query{
	Query{
		key: "loginQuery",
		query: `CALL sp_front_login(
			?,?,?,?,?,?,?,?,?)`,
	}, Query{
		key:   "getClientQuery",
		query: "SELECT client_name FROM Master_Client where client_id = ? LIMIT 1",
	}, Query{
		key:   "frontCheckSession",
		query: "CALL sp_front_check_session (?)",
	}, Query{
		key:   "logoutQuery",
		query: "CALL sp_front_logout (?,?)",
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
		key:   "getTask",
		query: `call sp_front_get_tasks(?)`,
	}, Query{
		key:   "addTask",
		query: `call sp_front_sync_tasks(?,?,?,?,?)`,
	}, Query{
		key:   "createApplicationQuery",
		query: `Insert into Master_Application (application_name) values(?)`,
	},
	//CMS
	Query{
		key:   "cmsAddUser",
		query: `call sp_cms_add_user (?,?,?,?,?,?,?,?,?,?)`,
	},
	Query{
		key:   "cmsEditUser",
		query: `call sp_cms_edit_user(?,?,?,?,?,?,?,?,?,?,?,?)`,
	},
	//todo
	Query{
		key:   "cmsCheckSubscription",
		query: `call sp_cms_check_subscription(?)`,
	}, Query{
		key:   "cmsGetActiveSubs",
		query: `call sp_cms_get_active_subs(?,?,?)`,
	},
	Query{
		key: "cmsEmployeeTreeFirstLevelAll",
		query: `select User_id, User_name, (select count(1) from Master_User where Client_id = ? and Superior_id = a.User_id) as Subs,
				CASE WHEN curdate() BETWEEN Active_start AND Active_end THEN 'Active' ELSE 'Inactive' END AS Status
				from Master_User a
				where Client_id = ? and Superior_id is null
				and Position_name <> '#Admin'
				order by Active_end desc, User_id asc;`,
	}, Query{
		key: "cmsEmployeeTreeFirstLevelActive",
		query: `select User_id, User_name, (select count(1) from Master_User where Client_id = ? and Superior_id = a.User_id) as Subs,
				CASE WHEN curdate() BETWEEN Active_start AND Active_end THEN 'Active' ELSE 'Inactive' END AS Status
				from Master_User a
				where Client_id = ? and Superior_id is null
				and Position_name <> '#Admin'
				and curdate() BETWEEN Active_start AND Active_end -- active only
				order by Active_end desc, User_id asc;`,
	},
	Query{
		key: "cmsEmployeeTreeSubsActive",
		query: `select User_id, User_name, (select count(1) from Master_User where Superior_id = a.User_id) as Subs,
				CASE WHEN curdate() BETWEEN Active_start AND Active_end THEN 'Active' ELSE 'Inactive' END AS Status
				from Master_User a
				where Superior_id = ?
				and Position_name <> '#Admin'
				and curdate() BETWEEN Active_start AND Active_end -- active only
				order by Active_end desc, User_id asc;`,
	},
	Query{
		key: "cmsEmployeeTreeSubsAll",
		query: `select User_id, User_name, (select count(1) from Master_User where Superior_id = a.User_id) as Subs,
				CASE WHEN curdate() BETWEEN Active_start AND Active_end THEN 'Active' ELSE 'Inactive' END AS Status
				from Master_User a
				where Superior_id = ?
				and Position_name <> '#Admin'
				order by Active_end desc, User_id asc;`,
	},
	Query{
		key:   "cmsEmployeeTreeChangeSuperior",
		query: `update Master_User set Superior_id = ? where User_id = ?;`,
	},
	Query{
		key:   "cmsEmployeeGrid",
		query: `call sp_cms_get_all_employees (?,?)`,
	},
	Query{
		key: "cmsEmailValidation",
		query: `SELECT IF (
				(SELECT COUNT(1) FROM Master_User
				WHERE Client_id = ?
				AND Email = ?
				AND User_id <> ? ) > 0,
				'This Email Address is already registered', 'OK') AS Result;`,
	},
	Query{
		key:   "cmsSaveDepartments",
		query: `call sp_cms_save_departments (?,?,?)`,
	},
	Query{
		key:   "cmsEditDepartment",
		query: `call sp_cms_edit_department (?,?,?,?)`,
	},
	Query{
		key:   "cmsGetAllDepartments",
		query: `call sp_cms_get_all_departments (?,?)`,
	},
	Query{
		key:   "cmsGetActiveDepartments",
		query: `call sp_cms_get_active_departments (?)`,
	},
	Query{
		key:   "cmsChangePassword",
		query: `call sp_cms_change_password (?,?,?)`,
	},
	Query{
		key:   "loginCMSQuery",
		query: `call sp_cms_login (?,?,?)`,
	},
	Query{
		key:   "cmsGetFeatures",
		query: `call sp_cms_get_features (?,?,?)`,
	},
	Query{
		key:   "cmsGetSubs",
		query: `call sp_cms_get_subs (?, ?, ?)`,
	},
	Query{
		key:   "cmsGetFeatures",
		query: `call cmsGetActivities (?)`,
	},
	Query{
		key:   "cmsGetAllPositions",
		query: `call sp_cms_get_all_positions (?)`,
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
