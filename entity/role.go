package entity

type Role uint8

const (
	UserRole Role = iota + 1
	AdminRole
)

const (
	UseRoleStr   = "user"
	AdminRoleStr = "admin"
)

func (r Role) String() string {
	switch r {
	case UserRole:
		return UseRoleStr
	case AdminRole:
		return AdminRoleStr
	default:
		return ""
	}
}

func MapToRoleEntity(role string) Role {
	switch role {
	case UseRoleStr:
		return UserRole
	case AdminRoleStr:
		return AdminRole
	default:
		return Role(0)
	}
}
