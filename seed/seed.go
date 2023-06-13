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

func GenerateVehicleType(vehicleType []string) {
	for _, vehicleType := range vehicleType {
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
				Username: fake.Person().LastName(),
				Password: fake.Internet().Password(),
				Email:    fake.Internet().Email(),
				RoleId:   role,
				DetailId: detail.Id,
			}
			err := config.DB.Create(&user).Error
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func GenerateOil(num int) {
	fake := faker.New()
	for i := 0; i < num; i++ {
		oil := models.Oil{
			Name: fake.Person().LastName(),
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
	roles := []string{"ADMIN_MADAM", "DRIVER", "PETUGAS", "USER"}
	vehicleType := []string{"SHIP", "TRUCK"}
	GenerateRoles(roles)
	GenerateVehicleType(vehicleType)
	GenerateOil(100)

	generateVehicle(100, 1)
	generateVehicle(100, 2)

	GenerateFakeUsers(100, 1)
	GenerateFakeUsers(100, 2)
	GenerateFakeUsers(100, 3)
	GenerateFakeUsers(100, 4)
	fmt.Println("Migration finished")
}
