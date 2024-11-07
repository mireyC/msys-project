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
	Children []*ProjectMenuTree
}

func BuildProjectMenuTree(menus []*ProjectMenu, Pid int64, res *ProjectMenuTree) {
	for _, menu := range menus {
		if menu.Pid == Pid {
			child := &ProjectMenuTree{}
			copier.Copy(child, menu)
			res.Children = append(res.Children, child)
			BuildProjectMenuTree(menus, menu.Id, child)
		}
	}
}
