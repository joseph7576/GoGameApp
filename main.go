package main

import (
	"GoGameApp/entity"
	"GoGameApp/repository/mysql"
	"fmt"
)

func main() {
	testUserMySQLRepo()
}

func testUserMySQLRepo() {
	mysqlRepo := mysql.New()

	createdUser, err := mysqlRepo.CreateUser(entity.User{
		ID:          0,
		PhoneNumber: "09021650189",
		Name:        "Test User 2",
	})
	if err != nil {
		fmt.Println("can't create user -> ", err)
	} else {
		fmt.Println("created user -> ", createdUser)
	}

	isUnique, err := mysqlRepo.IsPhoneNumberUnique(createdUser.PhoneNumber + "13")
	if err != nil {
		fmt.Println("unique err", err)
	}

	fmt.Println("isUnique=", isUnique)
}
