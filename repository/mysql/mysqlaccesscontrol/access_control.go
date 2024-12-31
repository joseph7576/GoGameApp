package mysqlaccesscontrol

import (
	"GoGameApp/entity"
	"GoGameApp/pkg/errmsg"
	"GoGameApp/pkg/richerror"
	"GoGameApp/pkg/slice"
	"GoGameApp/repository/mysql"
	"strings"
)

func (d *DB) GetUserPermissionTitles(userID uint, role entity.Role) ([]entity.PermissionTitle, error) {
	const op = "mysql.GetUserPermissionTitles"

	roleACL := make([]entity.AccessControl, 0)
	rows, err := d.conn.Conn().Query(`select * from access_controls where actor_type = ? and actor_id = ?`,
		entity.RoleActorType, role)
	if err != nil {
		return nil, richerror.New(op).WithErr(err).
			WithKind(richerror.KindUnexpected).WithMessage(errmsg.ErrMsgSomethingWentWrong)
	}

	defer rows.Close()

	for rows.Next() {
		acl, err := mysql.ScanAccessControl(rows)
		if err != nil {
			return nil, richerror.New(op).WithErr(err).
				WithKind(richerror.KindUnexpected).WithMessage(errmsg.ErrMsgSomethingWentWrong)
		}

		roleACL = append(roleACL, acl)
	}

	if err := rows.Err(); err != nil {
		return nil, richerror.New(op).WithErr(err).
			WithKind(richerror.KindUnexpected).WithMessage(errmsg.ErrMsgSomethingWentWrong)
	}

	userACL := make([]entity.AccessControl, 0)

	rows, err = d.conn.Conn().Query(`select * from access_controls where actor_type = ? and actor_id = ?`,
		entity.UserActorType, userID)
	if err != nil {
		return nil, richerror.New(op).WithErr(err).
			WithKind(richerror.KindUnexpected).WithMessage(errmsg.ErrMsgSomethingWentWrong)
	}

	defer rows.Close()

	for rows.Next() {
		acl, err := mysql.ScanAccessControl(rows)
		if err != nil {
			return nil, richerror.New(op).WithErr(err).
				WithKind(richerror.KindUnexpected).WithMessage(errmsg.ErrMsgSomethingWentWrong)
		}

		userACL = append(userACL, acl)
	}

	if err := rows.Err(); err != nil {
		return nil, richerror.New(op).WithErr(err).
			WithKind(richerror.KindUnexpected).WithMessage(errmsg.ErrMsgSomethingWentWrong)
	}

	permissionIDs := make([]uint, 0)
	for _, r := range roleACL {
		if !slice.DoesExist(permissionIDs, r.PermissionID) {
			permissionIDs = append(permissionIDs, r.PermissionID)
		}
	}

	if len(permissionIDs) == 0 {
		return nil, nil
	}

	// select * from permissions where id in (?,?,?) -> dynamic args
	args := make([]any, len(permissionIDs))

	for i, id := range permissionIDs {
		args[i] = id
	}

	rows, err = d.conn.Conn().Query("select * from permissions where id in (?"+strings.Repeat(",?", len(permissionIDs)-1)+")", args...)
	if err != nil {
		return nil, richerror.New(op).WithErr(err).
			WithKind(richerror.KindUnexpected).WithMessage(errmsg.ErrMsgSomethingWentWrong)
	}

	defer rows.Close()

	permissionTitles := make([]entity.PermissionTitle, 0)

	for rows.Next() {
		permission, err := mysql.ScanPermission(rows)
		if err != nil {
			return nil, richerror.New(op).WithErr(err).
				WithKind(richerror.KindUnexpected).WithMessage(errmsg.ErrMsgSomethingWentWrong)
		}

		permissionTitles = append(permissionTitles, permission.Title)
	}

	if err := rows.Err(); err != nil {
		return nil, richerror.New(op).WithErr(err).
			WithKind(richerror.KindUnexpected).WithMessage(errmsg.ErrMsgSomethingWentWrong)
	}

	return permissionTitles, nil
}
