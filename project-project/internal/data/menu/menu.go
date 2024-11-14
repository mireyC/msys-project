package menu

import "github.com/jinzhu/copier"

type ProjectMenu struct {
	Id         int64
	Pid        int64
	Title      string
	Icon       string
	Url        string
	FilePath   string
	Params     string
	Node       string
	Sort       int
	Status     int
	CreateBy   int64
	IsInner    int
	Values     string
	ShowSlider int
}

func (*ProjectMenu) TableName() string {
	return "ms_project_menu"
}

type ProjectMenuTree struct {
	ProjectMenu
	StatusText string
	InnerText  string
	FullUrl    string
	Children   []*ProjectMenuTree
}

func BuildProjectMenuTree(menus []*ProjectMenu, Pid int64, res *ProjectMenuTree) {
	for _, menu := range menus {
		if menu.Pid == Pid {
			child := &ProjectMenuTree{}
			child.FullUrl = getFullUrl(menu.Url, menu.Params, menu.Values)
			child.StatusText = getStatus(child.Status)
			child.InnerText = getInnerText(child.IsInner)
			copier.Copy(child, menu)
			res.Children = append(res.Children, child)
			BuildProjectMenuTree(menus, menu.Id, child)
		}
	}
}

func getFullUrl(url string, params string, values string) string {
	if (params != "" && values != "") || values != "" {
		return url + "/" + values
	}
	return url
}

func getInnerText(inner int) string {
	if inner == 0 {
		return "导航"
	}
	if inner == 1 {
		return "内页"
	}
	return ""
}

func getStatus(status int) string {
	if status == 0 {
		return "禁用"
	}
	if status == 1 {
		return "使用中"
	}
	return ""
}
