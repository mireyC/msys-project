package repo

import (
	"context"
	"mirey7/project-user/internal/data/member"
	"mirey7/project-user/internal/database"
)

type MemberRepo interface {
	GetMemberByEmail(ctx context.Context, email string) (bool, error)
	GetMemberByMobile(ctx context.Context, mobile string) (bool, error)
	GetMemberByAccount(ctx context.Context, account string) (bool, error)
	SaveMember(conn database.DBConn, ctx context.Context, mem *member.Member) error
	FindMember(ctx context.Context, account string, password string) (*member.Member, error)
}
