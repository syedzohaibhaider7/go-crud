package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MARK: Models
type User struct {
	ID       int       `json:"id" gorm:"primaryKey"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Gender   string    `json:"gender"`
	Age      int       `json:"age"`
	Products []Product `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

type Product struct {
	ID     int    `json:"id" gorm:"primaryKey"`
	UserID int    `json:"user_id"`
	User   User   `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name   string `json:"name"`
	Price  int    `json:"price"`
}

// MARK: Main Function
func main() {
	dsn := "root:root@tcp(localhost:3307)/go_crud?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}

	router := gin.Default()
	db.AutoMigrate(&User{}, &Product{})

	userRoutes := router.Group("/user")
	{
		userRoutes.GET("/get/:id", func(context *gin.Context) {
			getUser(context, db)
		})
		userRoutes.GET("/list", func(context *gin.Context) {
			listUsers(context, db)
		})
		userRoutes.POST("/create", func(context *gin.Context) {
			createUser(context, db)
		})
		userRoutes.PATCH("/update/:id", func(context *gin.Context) {
			updateUser(context, db)
		})
		userRoutes.DELETE("/delete/:id", func(context *gin.Context) {
			deleteUser(context, db)
		})
	}

	productRoutes := router.Group("/product")
	{
		productRoutes.GET("/get/:id", func(context *gin.Context) {
			getProduct(context, db)
		})
		productRoutes.GET("/list", func(context *gin.Context) {
			listProducts(context, db)
		})
		productRoutes.POST("/create", func(context *gin.Context) {
			createProduct(context, db)
		})
		productRoutes.PATCH("/update/:id", func(context *gin.Context) {
			updateProduct(context, db)
		})
		productRoutes.DELETE("/delete/:id", func(context *gin.Context) {
			deleteProduct(context, db)
		})
	}

	router.GET("/user-products/:id", func(context *gin.Context) {
		userProducts(context, db)
	})
	router.GET("/product-owner/:id", func(context *gin.Context) {
		productOwner(context, db)
	})

	router.Run(":7000")
}

// MARK: User Functions
func getUser(context *gin.Context, db *gorm.DB) {
	id := context.Param("id")
	var user User

	err := db.Where("id = ?", id).First(&user).Error
	if err != nil {
		fmt.Println(err.Error())
		context.JSON(404, gin.H{
			"data":    nil,
			"message": nil,
			"error":   "user not found",
		})
		return
	}

	context.JSON(200, gin.H{
		"data":    user,
		"message": "user found",
		"error":   nil,
	})
}

func listUsers(context *gin.Context, db *gorm.DB) {
	var users []User

	err := db.Select("*").Find(&users).Error

	if err != nil {
		fmt.Println(err.Error())
		context.JSON(404, gin.H{
			"data":    nil,
			"message": nil,
			"error":   "no record found",
		})
		return
	}

	var message string
	if len(users) == 0 {
		message = "no record found"
	} else {
		message = "users found"
	}

	context.JSON(200, gin.H{
		"data":       users,
		"totalCount": len(users),
		"message":    message,
		"error":      nil,
	})
}

func createUser(context *gin.Context, db *gorm.DB) {
	var user User

	context.Request.ParseForm()

	user.Name = context.PostForm("name")
	user.Email = context.PostForm("email")
	user.Gender = context.PostForm("gender")
	userAge, err := strconv.Atoi(context.PostForm("age"))
	if err != nil {
		fmt.Println(err.Error())
		context.JSON(400, gin.H{
			"data":    nil,
			"message": nil,
			"error":   "invalid age format",
		})
		return
	}
	user.Age = userAge

	err = db.Create(&user).Error
	if err != nil {
		fmt.Println(err.Error())
		context.JSON(404, gin.H{
			"data":    nil,
			"message": nil,
			"error":   "user couldn't be created",
		})
		return
	}

	context.JSON(200, gin.H{
		"data":    user,
		"message": "user created successfully",
		"error":   nil,
	})
}

func updateUser(context *gin.Context, db *gorm.DB) {
	id := context.Param("id")
	var user User

	err := db.Where("id = ?", id).First(&user).Error
	if err != nil {
		fmt.Println(err.Error())
		context.JSON(404, gin.H{
			"data":    nil,
			"message": nil,
			"error":   "user not found",
		})
		return
	}

	context.Request.ParseForm()

	userName := context.PostForm("name")
	if userName != "" {
		user.Name = userName
	}
	userEmail := context.PostForm("email")
	if userEmail != "" {
		user.Email = userEmail
	}
	userGender := context.PostForm("gender")
	if userGender != "" {
		user.Gender = userGender
	}
	age := context.PostForm("age")
	if age != "" {
		userAge, err := strconv.Atoi(age)
		if err != nil {
			fmt.Println(err.Error())
			context.JSON(400, gin.H{
				"data":    nil,
				"message": nil,
				"error":   "invalid age format",
			})
			return
		}
		user.Age = userAge
	}

	err = db.Where("id = ?", id).Updates(&user).Error
	if err != nil {
		fmt.Println(err.Error())
		context.JSON(404, gin.H{
			"data":    nil,
			"message": nil,
			"error":   "user couldn't be updated",
		})
		return
	}

	context.JSON(200, gin.H{
		"data":    user,
		"message": "user updated successfully",
		"error":   nil,
	})
}

func deleteUser(context *gin.Context, db *gorm.DB) {
	id := context.Param("id")
	var user User

	err := db.Where("id = ?", id).First(&user).Error
	if err != nil {
		fmt.Println(err.Error())
		context.JSON(404, gin.H{
			"data":    nil,
			"message": nil,
			"error":   "user not found",
		})
		return
	}

	err = db.Where("id = ?", id).Delete(&user).Error
	if err != nil {
		fmt.Println(err.Error())
		context.JSON(404, gin.H{
			"data":    nil,
			"message": nil,
			"error":   "user couldn't be deleted",
		})
		return
	}

	context.JSON(200, gin.H{
		"data":    nil,
		"message": "user deleted successfully",
		"error":   nil,
	})
}

// MARK: Product Functions
func getProduct(context *gin.Context, db *gorm.DB) {
	id := context.Param("id")
	var product Product

	// Fetch product with associated user details
	// err := db.Preload("User").Where("id = ?", id).First(&product).Error
	err := db.Where("id = ?", id).First(&product).Error
	if err != nil {
		fmt.Println(err.Error())
		context.JSON(404, gin.H{
			"data":    nil,
			"message": nil,
			"error":   "product not found",
		})
		return
	}

	context.JSON(200, gin.H{
		"data":    product,
		"message": "product found",
		"error":   nil,
	})
}

func listProducts(context *gin.Context, db *gorm.DB) {
	var products []Product

	err := db.Select("*").Find(&products).Error

	if err != nil {
		fmt.Println(err.Error())
		context.JSON(404, gin.H{
			"data":    nil,
			"message": nil,
			"error":   "no record found",
		})
		return
	}

	var message string
	if len(products) == 0 {
		message = "no record found"
	} else {
		message = "products found"
	}

	context.JSON(200, gin.H{
		"data":       products,
		"totalCount": len(products),
		"message":    message,
		"error":      nil,
	})
}

func createProduct(context *gin.Context, db *gorm.DB) {
	var product Product
	var user User

	context.Request.ParseForm()

	productUserID, err := strconv.Atoi(context.PostForm("user_id"))
	if err != nil {
		fmt.Println(err.Error())
		context.JSON(400, gin.H{
			"data":    nil,
			"message": nil,
			"error":   "invalid user ID format",
		})
		return
	}

	err = db.Where("id = ?", productUserID).First(&user).Error
	if err != nil {
		fmt.Println(err.Error())
		context.JSON(404, gin.H{
			"data":    nil,
			"message": nil,
			"error":   "user not found",
		})
		return
	}

	product.UserID = productUserID
	product.Name = context.PostForm("name")
	productPrice, err := strconv.Atoi(context.PostForm("price"))
	if err != nil {
		fmt.Println(err.Error())
		context.JSON(400, gin.H{
			"data":    nil,
			"message": nil,
			"error":   "invalid age format",
		})
		return
	}
	product.Price = productPrice

	err = db.Create(&product).Error
	if err != nil {
		fmt.Println(err.Error())
		context.JSON(404, gin.H{
			"data":    nil,
			"message": nil,
			"error":   "product couldn't be created",
		})
		return
	}

	// var createdProduct Product
	// err = db.Preload("User").First(&createdProduct, product.ID).Error
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	context.JSON(404, gin.H{
	// 		"data":    nil,
	// 		"message": nil,
	// 		"error":   "user not found",
	// 	})
	// 	return
	// }

	context.JSON(200, gin.H{
		"data":    product,
		"message": "product created successfully",
		"error":   nil,
	})
}

func updateProduct(context *gin.Context, db *gorm.DB) {
	id := context.Param("id")
	var product Product
	var user User

	err := db.Where("id = ?", id).First(&product).Error
	if err != nil {
		fmt.Println(err.Error())
		context.JSON(404, gin.H{
			"data":    nil,
			"message": nil,
			"error":   "product not found",
		})
		return
	}

	context.Request.ParseForm()

	userID := context.PostForm("user_id")
	if userID != "" {
		productUserID, err := strconv.Atoi(userID)
		if err != nil {
			fmt.Println(err.Error())
			context.JSON(400, gin.H{
				"data":    nil,
				"message": nil,
				"error":   "invalid age format",
			})
			return
		}

		err = db.Where("id = ?", productUserID).First(&user).Error
		if err != nil {
			fmt.Println(err.Error())
			context.JSON(404, gin.H{
				"data":    nil,
				"message": nil,
				"error":   "user not found",
			})
			return
		}
		product.UserID = productUserID
	}
	productName := context.PostForm("name")
	if productName != "" {
		product.Name = productName
	}
	price := context.PostForm("price")
	if price != "" {
		productPrice, err := strconv.Atoi(price)
		if err != nil {
			fmt.Println(err.Error())
			context.JSON(400, gin.H{
				"data":    nil,
				"message": nil,
				"error":   "invalid age format",
			})
			return
		}
		product.Price = productPrice
	}

	err = db.Where("id = ?", id).Updates(&product).Error
	if err != nil {
		fmt.Println(err.Error())
		context.JSON(404, gin.H{
			"data":    nil,
			"message": nil,
			"error":   "product couldn't be updated",
		})
		return
	}

	// var updatedProduct Product
	// err = db.Preload("User").First(&updatedProduct, id).Error
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	context.JSON(404, gin.H{
	// 		"data":    nil,
	// 		"message": nil,
	// 		"error":   "user not found",
	// 	})
	// 	return
	// }

	context.JSON(200, gin.H{
		"data":    product,
		"message": "product updated successfully",
		"error":   nil,
	})
}

func deleteProduct(context *gin.Context, db *gorm.DB) {
	id := context.Param("id")
	var product Product

	err := db.Where("id = ?", id).First(&product).Error
	if err != nil {
		fmt.Println(err.Error())
		context.JSON(404, gin.H{
			"data":    nil,
			"message": nil,
			"error":   "product not found",
		})
		return
	}

	err = db.Where("id = ?", id).Delete(&product).Error
	if err != nil {
		fmt.Println(err.Error())
		context.JSON(404, gin.H{
			"data":    nil,
			"message": nil,
			"error":   "product couldn't be deleted",
		})
		return
	}

	context.JSON(200, gin.H{
		"data":    nil,
		"message": "product deleted successfully",
		"error":   nil,
	})
}

// MARK: User Products API
func userProducts(context *gin.Context, db *gorm.DB) {
	id := context.Param("id")
	var products []Product
	var user User

	err := db.Where("id = ?", id).First(&user).Error
	if err != nil {
		fmt.Println(err.Error())
		context.JSON(404, gin.H{
			"data":    nil,
			"message": nil,
			"error":   "user not found",
		})
		return
	}

	err = db.Where("user_id = ?", id).Find(&products).Error
	if err != nil {
		context.JSON(404, gin.H{
			"data":    nil,
			"error":   "no record found",
			"message": nil,
		})
		return
	}

	context.JSON(200, gin.H{
		"data":       products,
		"totalCount": len(products),
		"error":      nil,
		"message":    "products found",
	})
}

// MARK: Product Owner API
func productOwner(context *gin.Context, db *gorm.DB) {
	id := context.Param("id")
	var product Product
	var user User

	err := db.Where("id = ?", id).First(&product).Error
	if err != nil {
		fmt.Println(err.Error())
		context.JSON(404, gin.H{
			"data":    nil,
			"message": nil,
			"error":   "product not found",
		})
		return
	}

	err = db.Where("id = ?", product.UserID).First(&user).Error
	if err != nil {
		context.JSON(404, gin.H{
			"data":    nil,
			"error":   "user not found",
			"message": nil,
		})
		return
	}

	context.JSON(200, gin.H{
		"data":    user,
		"error":   nil,
		"message": "user found",
	})
}
