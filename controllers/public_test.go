package controllers

import (
	"bytes"
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

func TestPublic(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Public Suite")
}

var _ = Describe("Public", func() {
	Context("when making a sign up action", func() {
		By("first it creates an user")

		var actualResult models.User
		user := models.User{
			Name:     "Test User",
			Email:    "jwt@email.com",
			Password: "secret",
		}

		By("then it makes the payload with the user object")
		payload, err := json.Marshal(&user)

		It("and verify that it didn't failed", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		By("then it creates an http request")
		request, err := http.NewRequest("POST", "/api/public/signup", bytes.NewBuffer(payload))

		It("and verify if it didn't failed", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		By("then it makes a database connection")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = request
		err = database.InitDatabase()

		It("and verifies that the connection didn't failed", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		database.GlobalDB.AutoMigrate(&models.User{})

		By("then it makes the signup call")
		Signup(c)

		It("and expect to receive a 200 http code response", func() {
			Expect(w.Code).To(Equal(200))
		})

		By("we unmarshal the response body")
		err = json.Unmarshal(w.Body.Bytes(), &actualResult)

		It("hoping that we didn't fail", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("we read the name", func() {
			Expect(user.Name).To(Equal(actualResult.Name))
		})
		It("and email", func() {
			Expect(user.Email).To(Equal(actualResult.Email))
		})
	})

	Context("when making a sign up with invalid json", func() {
		user := "test"
		payload, err := json.Marshal(&user)

		It("made a payload", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		request, err := http.NewRequest("POST", "/api/public/signup", bytes.NewBuffer(payload))

		It("created a http request", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = request

		Signup(c)

		It("and finally, received a 499 http error code", func() {
			Expect(w.Code).To(Equal(400))
		})
	})
	Context("when making a login request", func() {
		user := LoginPayload{
			Email:    "jwt@email.com",
			Password: "secret",
		}
		payload, err := json.Marshal(&user)

		It("created the payload", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		request, err := http.NewRequest("POST", "/api/public/login", bytes.NewBuffer(payload))

		It("made a request", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = request
		err = database.InitDatabase()

		It("made a successful call to the database", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		database.GlobalDB.AutoMigrate(&models.User{})
		Login(c)

		It("and received a 200 http code", func() {
			Expect(w.Code).To(Equal(200))
		})
	})
	Context("when making an invalid login request with json", func() {
		user := "test"
		payload, err := json.Marshal(&user)

		It("made the payload", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		request, err := http.NewRequest("POST", "/api/public/login", bytes.NewBuffer(payload))

		It("made an http request", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = request
		Login(c)

		It("and received a 400 error code", func() {
			Expect(w.Code).To(Equal(400))
		})
	})
	Context("when making a login action with invalid credentials", func() {
		By("first it creates a fake user")
		user := LoginPayload{
			Email:    "jwt@email.com",
			Password: "invalid",
		}

		By("then makes the payload")
		payload, err := json.Marshal(&user)

		It("and make sure it didn't fail", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		request, err := http.NewRequest("POST", "/api/public/login", bytes.NewBuffer(payload))

		It("then it makes the http request", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = request

		It("then it creates a connection to the database", func() {
			err = database.InitDatabase()

			Expect(err).ToNot(HaveOccurred())
		})

		database.GlobalDB.AutoMigrate(&models.User{})
		Login(c)
		database.GlobalDB.Unscoped().Where("email = ?", user.Email).Delete(&models.User{})

		It("and finally, the server responds with a 401 error code", func() {
			Expect(w.Code).To(Equal(401))
		})
	})
})
