package pro

type Project struct {
	Id                 int64   `json:"id" form:"id"`
	Cover              string  `json:"cover" form:"cover"`
	Name               string  `json:"name" form:"name"`
	Description        string  `json:"description" form:"description"`
	AccessControlType  string  `json:"access_control_type" form:"access_control_type"`
	WhiteList          string  `json:"white_list" form:"white_list"`
	Order              int     `json:"order" form:"order"`
	Deleted            int     `json:"deleted" form:"deleted"`
	TemplateCode       int     `json:"template_code" form:"template_code"`
	Schedule           float64 `json:"schedule" form:"schedule"`
	CreateTime         string  `json:"create_time" form:"create_time"`
	OrganizationCode   string  `json:"organization_code" form:"organization_code"`
	DeletedTime        string  `json:"deleted_time" form:"deleted_time"`
	Private            int     `json:"private" form:"private"`
	Prefix             string  `json:"prefix" form:"prefix"`
	OpenPrefix         int     `json:"open_prefix" form:"open_prefix"`
	Archive            int     `json:"archive" form:"archive"`
	ArchiveTime        int64   `json:"archive_time" form:"archive_time"`
	OpenBeginTime      int     `json:"open_begin_time" form:"open_begin_time"`
	OpenTaskPrivate    int     `json:"open_task_private" form:"open_task_private"`
	TaskBoardTheme     string  `json:"task_board_theme" form:"task_board_theme"`
	BeginTime          int64   `json:"begin_time" form:"begin_time"`
	EndTime            int64   `json:"end_time" form:"end_time"`
	AutoUpdateSchedule int     `json:"auto_update_schedule" form:"auto_update_schedule"`
	Code               string  `json:"code" form:"code"`
	ProjectCode        string  `json:"projectCode" form:"projectCode"`
}

type ProjectDetail struct {
	Project
	Collected   int    `json:"collected"`
	OwnerName   string `json:"owner_name"`
	OwnerAvatar string `json:"owner_avatar"`
}

type ProjectMember struct {
	Id          int64  `json:"id"`
	ProjectCode int64  `json:"project_code"`
	MemberCode  int64  `json:"member_code"`
	JoinTime    string `json:"join_time"`
	IsOwner     int64  `json:"is_owner"`
	Authorize   string `json:"authorize"`
}

type ProjectAndMember struct {
	Project
	ProjectCode int64  `json:"project_code"`
	MemberCode  int64  `json:"member_code"`
	JoinTime    int64  `json:"join_time"`
	IsOwner     int64  `json:"is_owner"`
	Authorize   string `json:"authorize"`
	OwnerName   string `json:"owner_name"`
	Collected   int    `json:"collected"`
}

type Menu struct {
	Id         int64   `json:"id"`
	Pid        int64   `json:"pid"`
	Title      string  `json:"title"`
	Icon       string  `json:"icon"`
	Url        string  `json:"url"`
	FilePath   string  `json:"filePath"`
	Params     string  `json:"params"`
	Node       string  `json:"node"`
	Sort       int32   `json:"sort"`
	Status     int32   `json:"status"`
	CreateBy   int64   `json:"create_by"`
	IsInner    int32   `json:"is_inner"`
	Values     string  `json:"values"`
	ShowSlider int32   `json:"show_slider"`
	StatusText string  `json:"statusText"`
	InnerText  string  `json:"innerText"`
	FullUrl    string  `json:"fullUrl"`
	Children   []*Menu `json:"children"`
}

type ProjectTemplateQueryReq struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"pageSize" form:"pageSize"`
	ViewType string `json:"viewType" form:"viewType"`
}

type ProjectTemplate struct {
	Id               int                   `json:"id"`
	Name             string                `json:"name"`
	Description      string                `json:"description"`
	Sort             int                   `json:"sort"`
	CreateTime       string                `json:"create_time"`
	OrganizationCode string                `json:"organization_code"`
	Cover            string                `json:"cover"`
	MemberCode       string                `json:"member_code"`
	IsSystem         int                   `json:"is_system"`
	TaskStages       []*TaskStagesOnlyName `json:"task_stages"`
	Code             string                `json:"code"`
}

type TaskStagesOnlyName struct {
	Name string `json:"name"`
}

type ProjectSaveReq struct {
	Name         string `json:"name" form:"name"`
	TemplateCode string `json:"templateCode" form:"templateCode"`
	Description  string `json:"description" form:"description"`
	Id           int    `json:"id" form:"id"`
}

type ProjectSaveResp struct {
	CreateTime       string `json:"create_time"`
	Code             string `json:"code"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	OrganizationCode string `json:"organizationCode"`
	TaskBoardTheme   string `json:"taskBoardTheme"`
	Cover            string `json:"cover"`
	Id               int64  `json:"id"`
}
