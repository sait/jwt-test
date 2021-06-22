package models

import (
	"fmt"
	"os"
	"testing"

	"github.com/AlanHerediaG/test-jwt/database"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestModels(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Models Suite")
}

var _ = Describe("Models", func() {
	Context("hash a password", func() {
		user := User{
			Password: "secret",
		}

		err := user.HashPassword(user.Password)
		It("hashes the password", func() {
			Expect(err).NotTo(HaveOccurred())
			os.Setenv("passwordHash", user.Password)
		})
	})

	Context("when creating a user", func() {
		var (
			user       User
			userResult User
			err        error
		)

		It("should create a successfull connection to the database", func() {
			err := database.InitDatabase()
			Expect(err).ShouldNot(HaveOccurred())

			By("and migrate the users table")
			err = database.GlobalDB.AutoMigrate(&User{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should create the new user in database", func() {
			user = User{
				Name:     "Test User",
				Email:    "test@email.com",
				Password: os.Getenv("passwordHash"),
			}

			err = user.CreateUserRecord()
			Expect(err).NotTo(HaveOccurred())

			database.GlobalDB.Where("email = ?", user.Email).Find(&userResult)
			database.GlobalDB.Unscoped().Delete(&user)
		})

		It("should return the user info correctly", func() {
			Expect("Test User").To(Equal(user.Name))
			Expect("test@email.com").To(Equal(user.Email))
		})
	})

	Context("when checking a user password", func() {
		It("then it should return a valid password", func() {
			By("get the user password hash from the enviroment")
			hash := os.Getenv("passwordHash")
			fmt.Fprintf(GinkgoWriter, "Password Hash: %v\n\n", hash)
			By("creating the user with it's hash")
			user := User{
				Password: hash,
			}

			By("then compare the password with provided")
			err := user.CheckPassword("secret")
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
