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
	for _, customerType := range []string{"Internal", "Pertamina"} {
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

func GenerateOil(oils, warehouse []string) {

	for _, wareHouseData := range warehouse {
		wareHouse := models.Warehouse{
			Name:       fmt.Sprintf("Warehouse %v", wareHouseData),
			Location:   "Jakarta",
			ProvinceId: 1,
			CityId:     2,
		}

		err := config.DB.Create(&wareHouse).Error
		if err != nil {
			fmt.Println(err)
		}

	}
	for _, oil := range oils {

		storage := models.Storage{
			Name:        fmt.Sprintf("%s Storage", oil),
			Quantity:    80_000,
			WarehouseID: 1,
			Capacity:    1000_000,
			OilID:       1,
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
		identifier := models.VehicleIdentifier{
			Identifier: fake.RandomLetter(),
		}

		err := config.DB.Create(&identifier).Error
		if err != nil {
			fmt.Println(err)
		}
		vehicle := models.Vehicle{
			Name:                fake.Car().Model(),
			VehicleTypeId:       vehicleType,
			VehicleIdentifierId: identifier.Id,
		}

		err = config.DB.Create(&vehicle).Error
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

func generateFakeCustomer(num int) {
	fake := faker.New()
	for i := 0; i < num/2; i++ {
		detail := GenerateFakeDetail()

		errCreateDetail := config.DB.Create(&detail).Error
		if errCreateDetail != nil {
			fmt.Println(errCreateDetail)
		}

		user := models.User{
			Username: fake.Person().Name(),
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
			DetailId:       detail.Id,
			CompanyID:      2,
			Phone:          fake.Phone().Number(),
		}
		errCreateCustomer := config.DB.Create(&customer).Error
		if err != nil {
			fmt.Println(errCreateCustomer)
		}
	}
	for i := 0; i < num/2; i++ {
		detail := GenerateFakeDetail()

		errCreateDetail := config.DB.Create(&detail).Error
		if errCreateDetail != nil {
			fmt.Println(errCreateDetail)
		}

		user := models.User{
			Username: fmt.Sprintf("Pertamina %v", fake.Person().Name()),
			Password: fake.Internet().Password(),
			Email:    fake.Internet().Email(),
		}

		err := config.DB.Create(&user).Error
		if err != nil {
			fmt.Println(err)
		}
		customer := models.Customer{
			UserId:         user.Id,
			CustomerTypeId: 2,
			DetailId:       detail.Id,
			CompanyID:      fake.IntBetween(1, 10),
			Phone:          fake.Phone().Number(),
		}
		errCreateCustomer := config.DB.Create(&customer).Error
		if err != nil {
			fmt.Println(errCreateCustomer)
		}
	}
}

type User struct {
	Username string
	Password string
	Email    string
	Phone    string
}

func generateSomeUser(data User, role int) {
	fake := faker.New()
	var user models.User

	var detail models.Detail
	err := config.DB.Create(&detail).Error
	if err != nil {
		fmt.Println(err)
	}

	user.Username = data.Username
	user.Password = data.Password
	user.Email = data.Email
	err = config.DB.Create(&user).Error
	if err != nil {
		fmt.Println(err)
	}

	employee := models.Employee{
		UserId:      user.Id,
		RoleId:      role,
		PhoneNumber: fake.Phone().Number(),
	}

	err = config.DB.Create(&employee).Error
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
	oils := []string{"HSD", "MFO"}
	warehouse := []string{"Jetty", "SPOB"}
	GenerateOil(oils, warehouse)

	generateVehicle(2, 1)
	generateVehicle(2, 2)

	GenerateFakeUsers(5)
	generateFakeCustomer(20)
	generateSomeUser(User{
		Username: "admin",
		Password: "admin",
		Email:    "admin@gmail.com",
		Phone:    "08123456789",
	}, 1)
	generateSomeUser(User{
		Username: "adminsales",
		Password: "adminsales",
		Email:    "adminsales@gmail.com",
		Phone:    "08123456789",
	}, 2)

	generateSomeUser(User{
		Username: "petugas",
		Password: "petugas",
		Email:    "petugas@gmail.com",
		Phone:    "08123456789",
	}, 3)

	generateSomeUser(User{
		Username: "dwiki",
		Password: "dwiki",
		Email:    "dwikiokvianp1999@gmail.com",
		Phone:    "08123456789",
	}, 4)
	generateSomeUser(User{
		Username: "driver1",
		Password: "driver1",
		Email:    "driver1@gmail.com",
		Phone:    "08123456789",
	}, 5)
	generateSomeUser(User{
		Username: "driver2",
		Password: "driver2",
		Email:    "driver2@gmail.com",
		Phone:    "08123456789",
	}, 5)
	generateSomeUser(User{
		Username: "driver3",
		Password: "driver3",
		Email:    "driver3@gmail.com",
		Phone:    "08123456789",
	}, 5)

	fmt.Println("Migration finished")
}
