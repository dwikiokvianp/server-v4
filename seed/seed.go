package main

import (
	"fmt"
	"github.com/jaswdr/faker"
	"github.com/joho/godotenv"
	"os"
	"server-v2/config"
	"server-v2/models"
)

func GenerateStatusType() {
	for _, statusType := range []string{"Delivery", "Pickup"} {
		statusType := models.StatusType{
			Name: statusType,
		}
		err := config.DB.Create(&statusType).Error
		if err != nil {
			fmt.Println(err)
		}
	}
}
func GenerateFakeDetail() models.Detail {
	detail := models.Detail{
		Balance: 0,
		Credit:  0,
	}
	return detail
}

func GenerateStatus(statusArr []string) {
	for _, status := range statusArr {
		status := models.Status{
			Name: status,
		}
		err := config.DB.Create(&status).Error
		if err != nil {
			fmt.Println(err)
		}
	}

	GenerateStatusType()

	statusTypeMappingData := []models.StatusTypeMapping{
		{
			StatusID:     1,
			StatusTypeID: 1,
		},
		{
			StatusID:     1,
			StatusTypeID: 2,
		},
		{
			StatusID:     2,
			StatusTypeID: 1,
		},
		{
			StatusID:     3,
			StatusTypeID: 1,
		},
		{
			StatusID:     3,
			StatusTypeID: 2,
		},
		{
			StatusID:     4,
			StatusTypeID: 2,
		},
		{
			StatusID:     5,
			StatusTypeID: 1,
		},
		{
			StatusID:     6,
			StatusTypeID: 1,
		},
		{
			StatusID:     7,
			StatusTypeID: 1,
		},
		{
			StatusID:     7,
			StatusTypeID: 2,
		},
		{
			StatusID:     8,
			StatusTypeID: 1,
		},
		{
			StatusID:     8,
			StatusTypeID: 2,
		},
	}

	for _, statusTypeMapping := range statusTypeMappingData {
		dataMapping := models.StatusTypeMapping{
			StatusID:     statusTypeMapping.StatusID,
			StatusTypeID: statusTypeMapping.StatusTypeID,
		}
		err := config.DB.Create(&dataMapping).Error
		if err != nil {
			fmt.Println(err)
		}
	}

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

func GenerateCustomerType() {
	for _, customerType := range []string{"Retail", "Corporate"} {
		customerType := models.CustomerType{
			Name: customerType,
		}
		err := config.DB.Create(&customerType).Error
		if err != nil {
			fmt.Println(err)
		}
	}
}

func GenerateFakeUsers(count int) {
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
			}
			err := config.DB.Create(&user).Error
			if err != nil {
				fmt.Println(err)
			}
			customer := models.Customer{
				UserId:         user.Id,
				CustomerTypeId: 1,
				DetailId:       1,
				CompanyID:      1,
				Phone:          fake.Internet().User(),
			}

			err = config.DB.Create(&customer).Error
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
	for i := 0; i < num; i++ {
		fake := faker.New()

		identifier := models.VehicleIdentifier{
			Identifier: fake.RandomLetter(),
		}

		err := config.DB.Create(&identifier).Error
		if err != nil {
			vehicle := models.Vehicle{
				Name:                fake.Car().Model(),
				VehicleTypeId:       vehicleType,
				VehicleIdentifierId: identifier.Id,
			}
			err := config.DB.Create(&vehicle).Error
			if err != nil {
				fmt.Println(err)
			}
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

func generateFakeEmployee(num, roleInput int) {
	for i := 0; i < num; i++ {
		fake := faker.New()

		user := models.User{
			Username: fake.Person().Name(),
			Password: fake.Internet().Password(),
			Email:    fake.Internet().Email(),
		}

		err := config.DB.Create(&user).Error
		if err != nil {
			fmt.Println(err)
		}
		employee := models.Employee{
			UserId: user.Id,
			RoleId: roleInput,
		}
		errCreateEmployee := config.DB.Create(&employee).Error
		if err != nil {
			fmt.Println(errCreateEmployee)
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

	status := []string{
		"PENDING APPROVAL",
		"APPROVED WAITING SCHEDULE",
		"APPROVED SCHEDULED",
		"WAITING CUSTOMER TO PICK UP",
		"PICKING UP FROM SUPPLY",
		"ON DELIVERY TO CUSTOMER",
		"FINISHED",
		"REJECTED",
	}

	GenerateStatus(status)
	GenerateCustomerType()

	vehicleTypes := []string{"SHIP", "TRUCK"}
	GenerateRoles(roles)
	GenerateVehicleType(vehicleTypes)

	generateFakeCompany(10)

	generateFakeEmployee(10, 1)
	generateFakeEmployee(10, 2)
	generateFakeEmployee(10, 3)
	generateFakeEmployee(10, 4)
	generateFakeEmployee(10, 5)
	oils := []string{"HSD", "MFO"}
	GenerateOil(oils, 1)

	generateVehicle(2, 1)
	generateVehicle(2, 2)

	GenerateFakeUsers(5)
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
		Username:  "driver1",
		Password:  "driver1",
		Email:     "driver1@gmail.com",
		RoleId:    5,
		CompanyID: 1,
		Phone:     "08123456789",
	})
	generateSomeUser(User{
		Username:  "driver2",
		Password:  "driver2",
		Email:     "driver2@gmail.com",
		RoleId:    5,
		CompanyID: 1,
		Phone:     "08123456789",
	})
	generateSomeUser(User{
		Username:  "driver3",
		Password:  "driver3",
		Email:     "driver3@gmail.com",
		RoleId:    5,
		CompanyID: 1,
		Phone:     "08123456789",
	})

	fmt.Println("Migration finished")
}
