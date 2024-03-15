package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
)

// EnsureValidToken is a middleware for the Gin framework that validates JWT tokens.
func EnsureValidToken() gin.HandlerFunc {
	// Parse the issuer URL from the AUTH0_DOMAIN environment variable.
	issuerURL, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/")
	if err != nil {
		// If the URL cannot be parsed, panic and stop the application startup.
		panic(fmt.Sprintf("Failed to parse the issuer url: %v", err))
	}

	// Create a JWKS provider that caches keys for 5 minutes to minimize requests to the JWKS endpoint.
	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	// Initialize the JWT validator using the JWKS provider for key resolution.
	jwtValidator, err := validator.New(
		provider.KeyFunc,                      // Key resolution function from the JWKS provider.
		validator.RS256,                       // Expected signing algorithm.
		issuerURL.String(),                    // Issuer URL to validate the "iss" claim.
		[]string{os.Getenv("AUTH0_AUDIENCE")}, // Audience to validate the "aud" claim.
		// Specify custom claims to be included in the validation process.
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				// Return an instance of CustomClaims, which must implement the validator.CustomClaims interface.
				return &CustomClaims{}
			},
		),
		// Allow a 1-minute clock skew to account for potential time differences between the server and token issuer.
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		// If the validator cannot be initialized, panic and stop the application startup.
		panic("Failed to set up the jwt validator")
	}

	// Return a Gin handler function that encapsulates the JWT validation logic.
	return func(c *gin.Context) {
		// Extract the Authorization header from the request.
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// If the header is missing, return an unauthorized response and stop processing the request.
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		// Split the Authorization header to separate the bearer token prefix from the token itself.
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			// If the header is malformed, return an unauthorized response and stop processing the request.
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is malformed"})
			c.Abort()
			return
		}

		// Validate the token using the configured JWT validator.
		_, err := jwtValidator.ValidateToken(c, headerParts[1])
		if err != nil {
			// If the token is invalid, return an unauthorized response and stop processing the request.
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// If the token is valid, proceed with processing the request.
		c.Next()
	}
}

// CustomClaims represents the custom claims you expect in the JWT token.
// You need to define this struct and ensure it implements the validator.CustomClaims interface,
// including the Validate method.
type CustomClaims struct {
	// Define your custom claims here, for example:
	Scope string `json:"scope"`
}

// Validate checks the custom claims. In this example, it does nothing but you can implement
// your own validation logic here.
func (c CustomClaims) Validate(ctx context.Context) error {
	// Implement custom validation logic if needed.
	return nil
}
