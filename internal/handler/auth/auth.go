package auth

import (
	"net/http"

	"backend/internal/svc"
	cognito "backend/pkg/cognito"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
)

type UserAttributes struct {
	Username           string `form:"username"`
	FirstName          string `form:"firstName"`
	LastName           string `form:"lastName"`
	Password           string `form:"password"`
	SubscriptionStatus string `form:"subscriptionStatus"`
}

type ErrorResponse struct {
	Message string
	Error   string
}

type SuccessResponse struct {
	Message string
	Data    interface{}
}

// @Summary Sign Up
// @Description Endpoint for signing up a user
// @Tags Auth
// @Accept multipart/form-data
// @Param username formData string true "Username"
// @Param nickname formData string true "Nickname"
// @Param fullName formData string true "Full Name"
// @Param phoneNumber formData string true "Phone Number"
// @Param password formData string true "Password"
// @Param email formData string true "Email"
// @Param photo formData file true "Profile Photo"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /signup [post]
func SignUp(s *svc.ServiceContext) echo.HandlerFunc {
	return func(c echo.Context) error {
		tracer := *s.Tracer
		_, span := tracer.Start(c.Request().Context(), "handler.SignUp")
		defer span.End()
		span.SetAttributes(attribute.String("http.method", "POST"), attribute.String("http.route", "/auth/signup"))

		var user UserAttributes
		err := c.Bind(&user)

		if err != nil {
			span.RecordError(err)
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "Something is missing, please check the form data and try again",
				"error":   err.Error(),
			})
		}

		if user.Username == "" || user.Password == "" || user.SubscriptionStatus == "" || user.FirstName == "" || user.LastName == "" {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "Username, password, first name and last name are required fields",
			})
		}

		cognito, _ := cognito.NewCognitoClient()
		input := &cognitoidentityprovider.SignUpInput{
			ClientId: aws.String(s.Config.AWS.COGNITO.CLIENT_ID),
			Username: aws.String(user.Username),
			Password: aws.String(user.Password),
			UserAttributes: []*cognitoidentityprovider.AttributeType{
				{
					Name:  aws.String("email"),
					Value: aws.String(user.Username),
				},
				{
					Name:  aws.String("given_name"),
					Value: aws.String(user.FirstName),
				},
				{
					Name:  aws.String("family_name"),
					Value: aws.String(user.LastName),
				},
				{
					Name:  aws.String("custom:subscription_status"),
					Value: aws.String("free"),
				},
			},
		}

		_, err = cognito.SignUp(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case cognitoidentityprovider.ErrCodeInvalidParameterException:
					span.RecordError(err)
					return c.JSON(http.StatusBadGateway, echo.Map{
						"message": "Email and Password is required arguments",
						"error":   err.Error(),
					})

				case cognitoidentityprovider.ErrCodeUsernameExistsException:
					span.RecordError(err)
					return c.JSON(http.StatusConflict, echo.Map{
						"message": "An account with the given email already exists",
						"error":   err.Error(),
					})

				case cognitoidentityprovider.ErrCodeInvalidPasswordException:
					span.RecordError(err)
					return c.JSON(http.StatusBadRequest, echo.Map{
						"message": "Password must include uppercase, special-character and number",
						"error":   err.Error(),
					})

				default:
					span.RecordError(err)
					return c.JSON(http.StatusInternalServerError, echo.Map{
						"message": "Something wen't wrong while sign up",
						"error":   err.Error(),
					})
				}
			} else {
				span.RecordError(err)
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"message": "Something wen't wrong while sign up",
					"error":   err.Error(),
				})
			}
		}

		return c.JSON(http.StatusOK, echo.Map{
			"message": "You have successfully signed up!",
		})
	}
}

// @Summary Sign In
// @Description Endpoint for signing in a user
// @Tags Auth
// @Accept multipart/form-data
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /signin [post]
func SignIn(s *svc.ServiceContext) echo.HandlerFunc {
	return func(c echo.Context) error {
		tracer := *s.Tracer
		_, span := tracer.Start(c.Request().Context(), "handler.SignIn")
		defer span.End()
		span.SetAttributes(attribute.String("http.method", "POST"), attribute.String("http.route", "/auth/signin"))

		var user UserAttributes
		err := c.Bind(&user)

		if err != nil {
			span.RecordError(err)
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "Something is missing, please check the form data and try again",
				"error":   err.Error(),
			})
		}

		if user.Username == "" || user.Password == "" {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "Username and Password are required fields",
			})
		}

		cognito, _ := cognito.NewCognitoClient()

		authInput := &cognitoidentityprovider.InitiateAuthInput{
			AuthFlow: aws.String(cognitoidentityprovider.AuthFlowTypeUserPasswordAuth),
			ClientId: aws.String(s.Config.AWS.COGNITO.CLIENT_ID),
			AuthParameters: map[string]*string{
				"USERNAME": aws.String(user.Username),
				"PASSWORD": aws.String(user.Password),
			},
		}

		authOutput, err := cognito.InitateAuth(authInput)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case cognitoidentityprovider.ErrCodeUserNotConfirmedException:
					span.RecordError(err)
					return c.JSON(http.StatusUnauthorized, echo.Map{
						"message": "Email is not confirmed.",
						"error":   err.Error(),
					})
				case cognitoidentityprovider.ErrCodeNotAuthorizedException:
					span.RecordError(err)
					return c.JSON(http.StatusUnauthorized, echo.Map{
						"message": "Incorrect email or password.",
						"error":   err.Error(),
					})

				default:
					span.RecordError(err)
					return c.JSON(http.StatusInternalServerError, echo.Map{
						"message": "Something wen't wrong while sign up",
						"error":   err.Error(),
					})
				}
			} else {
				span.RecordError(err)
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"message": "Something wen't wrong while sign up",
					"error":   err.Error(),
				})
			}
		}

		cookie := new(http.Cookie)
		cookie.Name = "token"
		cookie.Value = *authOutput.AuthenticationResult.IdToken
		c.SetCookie(cookie)

		return c.JSON(http.StatusOK, echo.Map{
			"message":      "You have successfully signed in!",
			"refreshToken": authOutput.AuthenticationResult.RefreshToken,
			"token":        authOutput.AuthenticationResult.IdToken,
		})
	}
}

// @Summary Verify Email
// @Description Endpoint for verifying a user's email
// @Tags Auth
// @Accept multipart/form-data
// @Param username formData string true "Username"
// @Param code formData string true "Verification Code"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /verify [post]
func VerifyEmail(s *svc.ServiceContext) echo.HandlerFunc {
	return func(c echo.Context) error {
		tracer := *s.Tracer
		_, span := tracer.Start(c.Request().Context(), "handler.VerifyEmail")
		defer span.End()
		span.SetAttributes(attribute.String("http.method", "POST"), attribute.String("http.route", "/auth/verify"))
		username := c.FormValue("username")
		code := c.FormValue("code")

		cognito, _ := cognito.NewCognitoClient()

		ConfirmSignUpInput := &cognitoidentityprovider.ConfirmSignUpInput{
			ClientId:         aws.String(s.Config.AWS.COGNITO.CLIENT_ID),
			Username:         aws.String(username),
			ConfirmationCode: aws.String(code),
		}
		_, err := cognito.ConfirmSignUp(ConfirmSignUpInput)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case cognitoidentityprovider.ErrCodeCodeMismatchException:
					span.RecordError(err)
					return c.JSON(http.StatusUnauthorized, echo.Map{
						"message": "Invalid verification code provided, please try again.",
						"error":   err.Error(),
					})
				case cognitoidentityprovider.ErrCodeExpiredCodeException:
					span.RecordError(err)
					return c.JSON(http.StatusUnauthorized, echo.Map{
						"message": "Verification code provided is expired please try again from start..",
						"error":   err.Error(),
					})

				default:
					span.RecordError(err)
					return c.JSON(http.StatusInternalServerError, echo.Map{
						"message": "Something wen't wrong while sign up",
						"error":   err.Error(),
					})
				}
			} else {
				span.RecordError(err)
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"message": "Something wen't wrong while sign up",
					"error":   err.Error(),
				})
			}
		}

		return c.JSON(http.StatusOK, echo.Map{
			"message": "Email verification successful!",
		})
	}
}

// @Summary Forgot Password
// @Description Endpoint for initiating the forgot password process
// @Tags Auth
// @Accept multipart/form-data
// @Param username formData string true "Username"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /forgot-password [post]
func ForgotPassword(s *svc.ServiceContext) echo.HandlerFunc {
	return func(c echo.Context) error {
		tracer := *s.Tracer
		_, span := tracer.Start(c.Request().Context(), "handler.ForgotPassword")
		defer span.End()
		span.SetAttributes(attribute.String("http.method", "POST"), attribute.String("http.route", "/auth/forgotpassword"))
		username := c.FormValue("username")
		cognito, _ := cognito.NewCognitoClient()

		ForgotPasswordInput := &cognitoidentityprovider.ForgotPasswordInput{
			ClientId: aws.String(s.Config.AWS.COGNITO.CLIENT_ID),
			Username: aws.String(username),
		}
		_, err := cognito.ForgotPassword(ForgotPasswordInput)
		if err != nil {
			span.RecordError(err)
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "Something went wrong!",
				"error":   err.Error(),
			})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"message": "Forgot password process initiated successfully!",
		})
	}
}

// @Summary Refresh Token
// @Description Endpoint for refreshing user token
// @Tags Auth
// @Accept multipart/form-data
// @Param username formData string true "Username"
// @Param accessToken formData string true "accessToken"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /refresh-token [post]
func RefreshToken(s *svc.ServiceContext) echo.HandlerFunc {
	return func(c echo.Context) error {
		tracer := *s.Tracer
		_, span := tracer.Start(c.Request().Context(), "handler.RefreshToken")
		defer span.End()
		span.SetAttributes(attribute.String("http.method", "POST"), attribute.String("http.route", "/auth/refresh-token"))

		var refreshTokenReq struct {
			RefreshToken string `form:"refreshToken"`
		}

		err := c.Bind(&refreshTokenReq)
		if err != nil {
			span.RecordError(err)
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "Invalid request format",
				"error":   err.Error(),
			})
		}

		if refreshTokenReq.RefreshToken == "" {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "Refresh token is required",
			})
		}

		cognito, _ := cognito.NewCognitoClient()

		// Use the AWS Cognito SDK to refresh the token
		refreshInput := &cognitoidentityprovider.InitiateAuthInput{
			AuthFlow: aws.String(cognitoidentityprovider.AuthFlowTypeRefreshToken),
			ClientId: aws.String(s.Config.AWS.COGNITO.CLIENT_ID),
			AuthParameters: map[string]*string{
				"REFRESH_TOKEN": aws.String(refreshTokenReq.RefreshToken),
			},
		}

		authOutput, err := cognito.InitateAuth(refreshInput)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case cognitoidentityprovider.ErrCodeNotAuthorizedException:
					span.RecordError(err)
					return c.JSON(http.StatusUnauthorized, echo.Map{
						"message": "Refresh token is invalid or expired",
					})
				case cognitoidentityprovider.ErrCodeUserNotConfirmedException:
					span.RecordError(err)
					return c.JSON(http.StatusUnauthorized, echo.Map{
						"message": "User email is not confirmed",
					})
				default:
					span.RecordError(err)
					return c.JSON(http.StatusInternalServerError, echo.Map{
						"message": "Something went wrong while refreshing the token",
					})
				}
			} else {
				// Handle other non-AWS Cognito errors
				span.RecordError(err)
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"message": "Something went wrong while refreshing the token",
				})
			}
		}

		// Set the new access token in the response
		token := *authOutput.AuthenticationResult.IdToken
		return c.JSON(http.StatusOK, echo.Map{
			"message": "Token refreshed successfully",
			"token":   token,
		})
	}
}

// @Summary Reset Password
// @Description Endpoint for resetting the password after initiating forgot password process
// @Tags Auth
// @Accept multipart/form-data
// @Param username formData string true "Username"
// @Param code formData string true "Verification Code"
// @Param newPassword formData string true "New Password"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /reset-password [post]
func ResetPassword(s *svc.ServiceContext) echo.HandlerFunc {
	return func(c echo.Context) error {
		tracer := *s.Tracer
		_, span := tracer.Start(c.Request().Context(), "handler.ResetPassword")
		defer span.End()
		span.SetAttributes(attribute.String("http.method", "POST"), attribute.String("http.route", "/auth/reset-password"))

		username := c.FormValue("username")
		code := c.FormValue("code")
		newPassword := c.FormValue("newPassword")

		if username == "" || code == "" || newPassword == "" {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "Username, code, and newPassword are required fields",
			})
		}

		cognito, _ := cognito.NewCognitoClient()

		resetInput := &cognitoidentityprovider.ConfirmForgotPasswordInput{
			ClientId:         aws.String(s.Config.AWS.COGNITO.CLIENT_ID),
			Username:         aws.String(username),
			ConfirmationCode: aws.String(code),
			Password:         aws.String(newPassword),
		}

		_, err := cognito.ConfirmForgotPassword(resetInput)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case cognitoidentityprovider.ErrCodeCodeMismatchException:
					span.RecordError(err)
					return c.JSON(http.StatusUnauthorized, echo.Map{
						"message": "Invalid verification code",
					})
				case cognitoidentityprovider.ErrCodeExpiredCodeException:
					span.RecordError(err)
					return c.JSON(http.StatusUnauthorized, echo.Map{
						"message": "Verification code has expired",
					})
				case cognitoidentityprovider.ErrCodeUserNotFoundException:
					span.RecordError(err)
					return c.JSON(http.StatusNotFound, echo.Map{
						"message": "User not found",
					})
				default:
					span.RecordError(err)
					return c.JSON(http.StatusInternalServerError, echo.Map{
						"message": "Something went wrong while resetting the password",
					})
				}
			} else {
				span.RecordError(err)
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"message": "Something went wrong while resetting the password",
				})
			}
		}

		return c.JSON(http.StatusOK, echo.Map{
			"message": "Password reset successfully!",
		})
	}
}
