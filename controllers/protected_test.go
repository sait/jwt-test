package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlanHerediaG/test-jwt/database"
	"github.com/AlanHerediaG/test-jwt/models"

	"github.com/gin-gonic/gin"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestProtected(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Protected Suite")
}

var _ = Describe("Protected", func() {
	Context("when trying to get a profile", func() {
		var profile models.User
		err := database.InitDatabase()

		It("it first reaches the database", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("then it migrates the user table")
		database.GlobalDB.AutoMigrate(&models.User{})

		user := models.User{
			Email:    "jwt@email.com",
			Password: "secret",
			Name:     "Test User",
		}

		It("hashes the user password", func() {
			err = user.HashPassword(user.Password)
			Expect(err).ToNot(HaveOccurred())
		})

		It("then it creates the user", func() {
			err = user.CreateUserRecord()
			Expect(err).ToNot(HaveOccurred())
		})

		It("then it makes the http request", func() {
			request, err := http.NewRequest("GET", "/api/protected/profile", nil)
			Expect(err).ToNot(HaveOccurred())

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = request

			By("inside the http request, it sets the email provided")
			c.Set("email", "jwt@email.com")

			By("and using that request, it fetches the profile")
			Profile(c)

			By("make sure we got a 200 code")
			Expect(w.Code).To(Equal(200))

			err = json.Unmarshal(w.Body.Bytes(), &profile)
			Expect(err).ToNot(HaveOccurred())
		})

		It("and that should get all the profile info correctly", func() {
			Expect(user.Email).To(Equal(profile.Email))
			Expect(user.Name).To(Equal(profile.Name))
		})
	})

	Context("when fetching an unexisting profile", func() {
		var profile models.User
		err := database.InitDatabase()

		It("first it tries to initialize the database", func() {
			Expect(err).ToNot(HaveOccurred())

			By("then it migrates the user table")
			database.GlobalDB.AutoMigrate(&models.User{})
		})

		It("then it makes an http request", func() {
			request, err := http.NewRequest("GET", "/api/protected/profile", nil)
			Expect(err).ToNot(HaveOccurred())

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = request

			By("inside of it, set the email address")
			c.Set("email", "notfound@email.com")

			By("and make the actual call")
			Profile(c)

			By("after getting the response, read all data")
			err = json.Unmarshal(w.Body.Bytes(), &profile)
			Expect(err).ToNot(HaveOccurred())

			By("delete the user form database")
			database.GlobalDB.Unscoped().Where("email = ?", "jwt@email.com").Delete(&models.User{})

			By("and make sure that the error code is NotFound")
			Expect(w.Code).To(Equal(404))
		})
	})
})
