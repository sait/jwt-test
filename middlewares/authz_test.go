package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlanHerediaG/test-jwt/auth"
	"github.com/AlanHerediaG/test-jwt/controllers"
	"github.com/AlanHerediaG/test-jwt/database"
	"github.com/AlanHerediaG/test-jwt/models"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestToken(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Token Suite")
}

var _ = Describe("Token", func() {
	var (
		router   *gin.Engine
		err      error
		response models.User
		req      *http.Request
	)

	BeforeEach(func() {
		router = gin.Default()
		router.Use(Authz())
		router.GET("/api/protected/profile", controllers.Profile)
		err := database.InitDatabase()
		Expect(err).ShouldNot(HaveOccurred())
	})

	JustBeforeEach(func() {
		err = database.GlobalDB.AutoMigrate(&models.User{})
		Expect(err).ShouldNot(HaveOccurred())
		req, err = http.NewRequest("GET", "/api/protected/profile", nil)
		Expect(err).ShouldNot(HaveOccurred())
	})

	Context("when using a headless token", func() {
		w := httptest.NewRecorder()

		It("make the http call", func() {
			router.ServeHTTP(w, req)

			By("and it reaches a forbidden error code")
			Expect(w.Code).To(Equal(403))
		})
	})

	Context("when using an invalid token format", func() {
		It("create the http request", func() {
			w := httptest.NewRecorder()

			By("insert an 'Authorization' header parameter, using an invalid value")
			req.Header.Add("Authorization", "test")

			By("make the http call")
			router.ServeHTTP(w, req)

			By("get the error code")
			Expect(w.Code).To(Equal(400))
		})
	})

	Context("when using an invalid token", func() {
		It("create the http request", func() {
			By("first define the custom handler for gin")
			invalidToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

			w := httptest.NewRecorder()
			By("insert an 'Authorization' header parameter, using an invalid token")
			req.Header.Add("Authorization", invalidToken)

			By("make the http call")
			router.ServeHTTP(w, req)

			By("and receive the 401 error code")
			Expect(w.Code).To(Equal(401))
		})
	})

	Context("when testing a valid token", func() {
		It("make the migration for users table", func() {
			err = database.GlobalDB.AutoMigrate(&models.User{})
			Expect(err).ToNot(HaveOccurred())
		})

		Describe("when generating the token", func() {
			By("first create the user")
			user := models.User{
				Email:    "test@email.com",
				Password: "secret",
				Name:     "Test User",
			}

			By("then create the JWT wrapper")
			jwtWrapper := auth.JwtWrapper{
				SecretKey:       "verysecretkey",
				Issuer:          "AuthService",
				ExpirationHours: 24,
			}

			By("generate the token using the provided user email")
			token, err := jwtWrapper.GenerateToken(user.Email)

			It("the token generated successfully", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("and the user was created with their password hashed", func() {
				err = user.HashPassword(user.Password)
				Expect(err).ToNot(HaveOccurred())

				_ = database.GlobalDB.Create(&user)
			})

			Context("After that, when calling the API and receiving the response", func() {
				w := httptest.NewRecorder()

				It("create the http request", func() {
					By("add the generated token to the header")
					req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

					By("make the http call")
					router.ServeHTTP(w, req)
					Expect(w.Code).To(Equal(200))
				})

				It("build the http request correctly", func() {
					err = json.Unmarshal(w.Body.Bytes(), &response)
					Expect(err).ToNot(HaveOccurred())

					By("delete the user from database")
					database.GlobalDB.Unscoped().Where("email = ?", user.Email).Delete(&models.User{})
				})

				It("and the response data is complete", func() {
					Expect(response.Email).To(Equal("test@email.com"))
					Expect(response.Name).To(Equal("Test User"))
				})
			})
		})
	})
})
