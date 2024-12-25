package mysql

import (
	"GoGameApp/entity"
	"GoGameApp/pkg/errmsg"
	"GoGameApp/pkg/richerror"
	"database/sql"
)

func (d *MySQLDB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	const op = "mysql.IsPhoneNumberUnique"

	row := d.db.QueryRow(`select * from users where phone_number = ?`, phoneNumber)
	_, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}

		return false, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrMsgCantScanQuery).WithKind(richerror.KindUnexpected)
	}

	return false, nil
}

func (d *MySQLDB) CreateUser(u entity.User) (entity.User, error) {
	const op = "mysql.CreateUser"

	res, err := d.db.Exec(`insert into users(name, phone_number, password) values(?, ?, ?)`, u.Name, u.PhoneNumber, u.Password)
	if err != nil {
		return entity.User{}, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrMsgCantExecCommand).WithKind(richerror.KindUnexpected)
	}

	id, _ := res.LastInsertId()
	u.ID = uint(id)

	return u, nil
}

func (d *MySQLDB) GetUserByPhoneNumber(phoneNumber string) (entity.User, bool, error) {
	const op = "mysql.GetUserByPhoneNumber"

	row := d.db.QueryRow(`select * from users where phone_number = ?`, phoneNumber)
	user, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, false, nil
		}

		return entity.User{}, false, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrMsgCantScanQuery).WithKind(richerror.KindUnexpected)
	}

	return user, true, nil
}

func (d *MySQLDB) GetUserByID(userID uint) (entity.User, error) {
	const op = "mysql.GetUserByID"

	row := d.db.QueryRow(`select * from users where id = ?`, userID)
	user, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, richerror.New(op).WithErr(err).
				WithMessage(errmsg.ErrMsgNotFound).WithKind(richerror.KindNotFound)
		}

		return entity.User{}, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrMsgCantScanQuery).WithKind(richerror.KindUnexpected)
	}

	return user, nil
}

func scanUser(row *sql.Row) (entity.User, error) {
	var (
		createdAt []uint8
		user      entity.User
	)

	err := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Password, &createdAt)

	return user, err
}
