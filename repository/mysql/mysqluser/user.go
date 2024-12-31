package mysqluser

import (
	"GoGameApp/entity"
	"GoGameApp/pkg/errmsg"
	"GoGameApp/pkg/richerror"
	"GoGameApp/repository/mysql"
	"database/sql"
)

func (d *DB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	const op = "mysql.IsPhoneNumberUnique"

	row := d.conn.Conn().QueryRow(`select * from users where phone_number = ?`, phoneNumber)
	_, err := mysql.ScanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}

		return false, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
	}

	return false, nil
}

func (d *DB) CreateUser(u entity.User) (entity.User, error) {
	const op = "mysql.CreateUser"

	res, err := d.conn.Conn().Exec(`insert into users(name, phone_number, password, role) values(?, ?, ?, ?)`,
		u.Name, u.PhoneNumber, u.Password, u.Role.String())
	if err != nil {
		return entity.User{}, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrMsgCantExecCommand).WithKind(richerror.KindUnexpected)
	}

	id, _ := res.LastInsertId()
	u.ID = uint(id)

	return u, nil
}

func (d *DB) GetUserByPhoneNumber(phoneNumber string) (entity.User, error) {
	const op = "mysql.GetUserByPhoneNumber"

	row := d.conn.Conn().QueryRow(`select * from users where phone_number = ?`, phoneNumber)
	user, err := mysql.ScanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, richerror.New(op).WithErr(err).
				WithMessage(errmsg.ErrMsgNotFound).WithKind(richerror.KindNotFound)
		}

		//TODO: log the unexpected error
		return entity.User{}, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
	}

	return user, nil
}

func (d *DB) GetUserByID(userID uint) (entity.User, error) {
	const op = "mysql.GetUserByID"

	row := d.conn.Conn().QueryRow(`select * from users where id = ?`, userID)
	user, err := mysql.ScanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, richerror.New(op).WithErr(err).
				WithMessage(errmsg.ErrMsgNotFound).WithKind(richerror.KindNotFound)
		}

		return entity.User{}, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
	}

	return user, nil
}
