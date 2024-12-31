package mysql

import (
	"GoGameApp/entity"
)

// polymorphism
type Scanner interface {
	Scan(dest ...any) error
}

func ScanAccessControl(scanner Scanner) (entity.AccessControl, error) {
	var (
		acl       entity.AccessControl
		createdAt []uint8
	)

	err := scanner.Scan(&acl.ID, &acl.ActorID, &acl.ActorType, &acl.PermissionID, &createdAt)
	return acl, err
}

func ScanPermission(scanner Scanner) (entity.Permission, error) {
	var (
		perm      entity.Permission
		createdAt []uint8
	)

	err := scanner.Scan(&perm.ID, &perm.Title, &createdAt)
	return perm, err
}

func ScanUser(scanner Scanner) (entity.User, error) {
	var (
		user      entity.User
		createdAt []uint8
		roleStr   string
	)

	err := scanner.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Password, &createdAt, &roleStr)

	user.Role = entity.MapToRoleEntity(roleStr)

	return user, err
}
