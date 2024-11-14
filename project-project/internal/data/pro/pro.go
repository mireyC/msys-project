package pro

type Project struct {
	Id                 int64
	Cover              string
	Name               string
	Description        string
	AccessControlType  int
	WhiteList          string
	Sort               int
	Deleted            int
	TemplateCode       string
	Schedule           float64
	CreateTime         int64
	OrganizationCode   int64
	DeletedTime        string
	Private            int
	Prefix             string
	OpenPrefix         int
	Archive            int
	ArchiveTime        int64
	OpenBeginTime      int
	OpenTaskPrivate    int
	TaskBoardTheme     string
	BeginTime          int64
	EndTime            int64
	AutoUpdateSchedule int
}

func (*Project) TableName() string {
	return "ms_project"
}

type ProjectMember struct {
	Id          int64
	ProjectCode int64
	MemberCode  int64
	JoinTime    int64
	IsOwner     int64
	Authorize   string
}

func (*ProjectMember) TableName() string {
	return "ms_project_member"
}

type ProjectAndMember struct {
	Project
	ProjectCode int64
	MemberCode  int64
	JoinTime    int64
	IsOwner     int64
	Authorize   string
}

func (m *ProjectAndMember) GetAccessControlType() string {
	if m.AccessControlType == 0 {
		return "open"
	}
	if m.AccessControlType == 1 {
		return "private"
	}
	if m.AccessControlType == 2 {
		return "custom"
	}
	return ""
}

func ToMap(orgs []*ProjectAndMember) map[int64]*ProjectAndMember {
	m := make(map[int64]*ProjectAndMember)
	for _, v := range orgs {
		m[v.Id] = v
	}
	return m
}

type CollectionProject struct {
	Id          int64
	ProjectCode int64
	MemberCode  int64
	CreateTime  int64
}

func (*CollectionProject) TableName() string {
	return "ms_project_collection"
}
