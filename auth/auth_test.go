package auth

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var token string

func TestToken(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Token Suite")
}

var _ = Describe("Token", func() {
	Context("when created", func() {
		jwtWrapper := JwtWrapper{
			SecretKey:       "verysecretkey",
			Issuer:          "AuthService",
			ExpirationHours: 24,
		}

		token, _ = jwtWrapper.GenerateToken("jwt@email.com")
		fmt.Fprintf(GinkgoWriter, "Generated token: %v\n", token)

		It("is", func() {
			Expect(token).ToNot(Equal(""))
		})

		Context("when validated", func() {
			get_token := token
			jwtWrapper := JwtWrapper{
				SecretKey: "verysecretkey",
				Issuer:    "AuthService",
			}

			claims, err := jwtWrapper.ValidateToken(get_token)

			It("doesn't return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})
			It("has the email", func() {
				Expect(claims.Email).To(Equal("jwt@email.com"))
			})
			It("has the issuer", func() {
				Expect(claims.Issuer).To(Equal("AuthService"))
			})
		})
	})
})
