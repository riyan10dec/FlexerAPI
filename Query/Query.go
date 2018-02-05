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
	}, Query{
		key:   "cmsEmployeeTreeFirstLevel",
		query: `call sp_cms_get_tree_first(?, ?, ?)`,
	},
	Query{
		key:   "cmsEmployeeTreeSubs",
		query: `call sp_cms_get_tree_subs (?,?,?)`,
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
		key:   "cmsGetActivities",
		query: `call sp_cms_get_all_activities (?,?)`,
	},

	Query{
		key:   "cmsSaveActivity",
		query: `call sp_cms_save_activity (?,?,?,?,?)`,
	},
	Query{
		key:   "cmsGetAllPositions",
		query: `call sp_cms_get_all_positions (?)`,
	},
	Query{
		key:   "cmsGetUserPerformance",
		query: `call sp_cms_get_user_performance (?,?,?)`,
	},
	Query{
		key:   "cmsGetUserDaily",
		query: `call sp_cms_get_user_daily (?,?,?,?)`,
	},
	Query{
		key:   "cmsGetUserDailyActivity",
		query: `call sp_cms_get_user_daily_activity (?,?,?)`,
	},
	Query{
		key:   "cmsGetUserDailyTimeline",
		query: `call sp_cms_get_user_daily_timeline (?,?)`,
	},
	Query{
		key:   "getUserTask",
		query: `call sp_cms_get_user_tasks (?,?,?,?)`,
	},
	Query{
		key:   "getNotificationQuery",
		query: `Select Notification_ID, Notification_Message, Page_URL, Seen from Notifications where User_ID = ?`,
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
