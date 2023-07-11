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
				CompanyID: 2,
				Phone:     fake.Phone().Number(),
			}
			err := config.DB.Create(&user).Error
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func GenerateOil(oils []string, warehouseId int) {
	detailWarehouse := models.WarehouseDetail{
		WarehouseID: uint64(warehouseId),
	}

	err := config.DB.Create(&detailWarehouse).Error
	if err != nil {
		fmt.Println(err)
	}

	for _, oil := range oils {

		storage := models.Storage{
			Name:              fmt.Sprintf("%s Storage", oil),
			Quantity:          80_000,
			WarehouseDetailID: detailWarehouse.ID,
			OilID:             1,
		}

		err := config.DB.Create(&storage).Error
		if err != nil {
			fmt.Println(err)
		}

		oil := models.Oil{
			Name:      oil,
			StorageId: int(storage.ID),
		}

		err = config.DB.Create(&oil).Error
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
			Address:        fake.Address().StreetAddress(),
			CompanyDetail:  fake.Company().BS(),
			CompanyZipCode: fake.RandomNumber(5),
			ProvinceId:     1,
			CityId:         3,
		}
		err := config.DB.Create(&company).Error
		if err != nil {
			fmt.Println(err)
		}
	}
}

func generateFakeEmployee(num int) {
	fake := faker.New()
	var warehouse models.Warehouse
	warehouse = models.Warehouse{
		Name:       fake.Car().Model(),
		ProvinceId: 2,
		CityId:     2,
		Location:   fake.Address().City(),
	}
	config.DB.Create(&warehouse)

	for i := 0; i < num; i++ {
		employee := models.Officer{
			Username:    fake.Person().LastName(),
			Password:    fake.Company().Suffix(),
			Email:       fake.Internet().Email(),
			PhoneNumber: fake.Phone().Number(),
			WarehouseId: warehouse.Id,
		}
		err := config.DB.Create(&employee).Error
		if err != nil {
			fmt.Println(err)
		}
	}
}

func generateDrivers(num int8) {
	fake := faker.New()
	for i := 0; i < int(num); i++ {
		driver := models.Driver{
			Username: fake.Person().LastName(),
			Password: fake.Company().Suffix(),
		}
		err := config.DB.Create(&driver).Error
		if err != nil {
			fmt.Println(err)
		}
	}
}

type User struct {
	Username  string
	Password  string
	Email     string
	RoleId    int
	CompanyID int
	Phone     string
}

func generateSomeUser(data User) {
	var admin models.User

	var detail models.Detail

	config.DB.Last(&detail)
	detail.Id = detail.Id + 1

	err := config.DB.Create(&detail).Error
	if err != nil {
		fmt.Println(err)
	}

	admin.Username = data.Username
	admin.Password = data.Password
	admin.Email = data.Email
	admin.RoleId = data.RoleId
	admin.CompanyID = data.CompanyID
	admin.Phone = data.Phone
	admin.DetailId = detail.Id
	err = config.DB.Create(&admin).Error
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	failedLoadEnv := godotenv.Load("./.env")
	if failedLoadEnv != nil {
		fmt.Println("Error loading .env file")
	}

	errDropDb := config.DropDatabase(os.Getenv("DB_URL"))
	if errDropDb != nil {
		fmt.Println(errDropDb)
	}

	err := config.InitDatabase(os.Getenv("DB_URL"))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Migration started")

	roles := []string{"ADMIN_PUSAT", "ADMIN_SALES", "OFFICER", "USER", "DRIVER"}
	vehicleTypes := []string{"SHIP", "TRUCK"}
	GenerateRoles(roles)
	GenerateVehicleType(vehicleTypes)

	generateFakeCompany(10)

	generateFakeEmployee(10)
	oils := []string{"SOLAR", "MFO"}
	GenerateOil(oils, 1)

	generateVehicle(10, 1)
	generateVehicle(10, 2)

	GenerateFakeUsers(100, 1)
	GenerateFakeUsers(100, 2)
	GenerateFakeUsers(100, 3)
	GenerateFakeUsers(100, 4)
	GenerateFakeUsers(100, 5)
	generateSomeUser(User{
		Username:  "admin",
		Password:  "admin",
		Email:     "admin@gmail.com",
		RoleId:    1,
		CompanyID: 1,
		Phone:     "08123456789",
	})
	generateSomeUser(User{
		Username:  "adminsales",
		Password:  "adminsales",
		Email:     "adminsales@gmail.com",
		RoleId:    2,
		CompanyID: 1,
		Phone:     "08123456789",
	})

	generateSomeUser(User{
		Username:  "petugas",
		Password:  "petugas",
		Email:     "petugas@gmail.com",
		RoleId:    3,
		CompanyID: 1,
		Phone:     "08123456789",
	})

	generateSomeUser(User{
		Username:  "dwiki",
		Password:  "dwiki",
		Email:     "dwikiokvianp1999@gmail.com",
		RoleId:    4,
		CompanyID: 1,
		Phone:     "08123456789",
	})
	generateSomeUser(User{
		Username:  "driver",
		Password:  "driver",
		Email:     "driver@gmail.com",
		RoleId:    5,
		CompanyID: 1,
		Phone:     "08123456789",
	})
	generateDrivers(100)

	fmt.Println("Migration finished")
}
