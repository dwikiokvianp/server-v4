package main

import (
	"fmt"
	"github.com/jaswdr/faker"
	"github.com/joho/godotenv"
	"os"
	"server-v2/config"
	"server-v2/models"
)

func GenerateFakeDetail() models.Detail {
	detail := models.Detail{
		Balance: 0,
		Credit:  0,
	}
	return detail
}

func GenerateRoles(roles []string) {
	for _, role := range roles {
		role := models.Role{
			Role: role,
		}
		err := config.DB.Create(&role).Error
		if err != nil {
			fmt.Println(err)
		}
	}
}

func GenerateVehicleType(vehicleTypes []string) {
	for _, vehicleType := range vehicleTypes {
		vehicleType := models.VehicleType{
			Name: vehicleType,
		}
		err := config.DB.Create(&vehicleType).Error
		if err != nil {
			fmt.Println(err)
		}
	}
}

func GenerateFakeUsers(count, role int) {
	fake := faker.New()
	for i := 0; i < count; i++ {
		detail := GenerateFakeDetail()
		err := config.DB.Create(&detail).Error
		if err != nil {
			fmt.Println(err)
		} else {
			user := models.User{
				Username:  fake.Person().LastName(),
				Password:  fake.Internet().Password(),
				Email:     fake.Internet().Email(),
				RoleId:    role,
				DetailId:  detail.Id,
				CompanyID: count,
				Phone:     fake.Phone().Number(),
			}
			err := config.DB.Create(&user).Error
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func GenerateOil(oils []string) {
	for _, oil := range oils {
		oil := models.Oil{
			Name: oil,
		}
		err := config.DB.Create(&oil).Error
		if err != nil {
			fmt.Println(err)
		}
	}
}

func generateVehicle(num, vehicleType int) {
	fake := faker.New()
	for i := 0; i < num; i++ {
		vehicle := models.Vehicle{
			Name:          fake.Car().Model(),
			VehicleTypeId: vehicleType,
		}
		err := config.DB.Create(&vehicle).Error
		if err != nil {
			fmt.Println(err)
		}
	}
}

func generateFakeCompany(num int) {
	fake := faker.New()
	for i := 0; i < num; i++ {
		company := models.Company{
			CompanyName:    fake.Company().Name(),
			Address:        fake.Address().Address(),
			CompanyDetail:  fake.Company().Suffix(),
			CompanyZipCode: fake.RandomNumber(5),
		}
		err := config.DB.Create(&company).Error
		if err != nil {
			fmt.Println(err)
		}
	}
}

func generateFakeEmployee(num int) {
	fake := faker.New()
	for i := 0; i < num; i++ {
		employee := models.Officer{
			Username: fake.Person().LastName(),
			Password: fake.Company().Suffix(),
			Email:    fake.Internet().Email(),
		}
		err := config.DB.Create(&employee).Error
		if err != nil {
			fmt.Println(err)
		}
	}
}

func main() {
	failedLoadEnv := godotenv.Load()
	if failedLoadEnv != nil {
		fmt.Println("Error loading .env file")
	}
	err := config.InitDatabase(os.Getenv("DB_URL"))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Migration started")
	oils := []string{"Bio Solar", "Premium", "Air"}
	roles := []string{"ADMIN", "PETUGAS", "USER"}
	vehicleTypes := []string{"SHIP", "TRUCK"}
	GenerateRoles(roles)
	GenerateVehicleType(vehicleTypes)
	GenerateOil(oils)
	generateFakeCompany(10)
	generateFakeEmployee(10)

	generateVehicle(10, 1)
	generateVehicle(10, 2)

	GenerateFakeUsers(10, 1)
	GenerateFakeUsers(10, 2)
	GenerateFakeUsers(10, 3)
	fmt.Println("Migration finished")
}
