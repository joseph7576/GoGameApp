package mysql

import (
	"GoGameApp/entity"
	"database/sql"
	"fmt"
	"time"
)

func (d *MySQLDB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	user := entity.User{}
	var createdAt time.Time

	row := d.db.QueryRow(`select * from users where phone_number = ?`, phoneNumber)
	err := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}

		return false, fmt.Errorf("can't scan query -> %w", err)
	}

	return false, nil
}

func (d *MySQLDB) CreateUser(u entity.User) (entity.User, error) {
	res, err := d.db.Exec(`insert into users(name, phone_number) values(?, ?)`, u.Name, u.PhoneNumber)
	if err != nil {
		return entity.User{}, fmt.Errorf("can't execute command -> %w", err)
	}

	id, _ := res.LastInsertId()
	u.ID = uint(id)

	return u, nil
}
