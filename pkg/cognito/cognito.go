package pkg

import (
	"backend/pkg/config"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type Client interface {
	SignUp(input *cognitoidentityprovider.SignUpInput) (*cognitoidentityprovider.SignUpOutput, error)
	InitateAuth(input *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error)
	ConfirmSignUp(input *cognitoidentityprovider.ConfirmSignUpInput) (*cognitoidentityprovider.ConfirmSignUpOutput, error)
	ForgotPassword(input *cognitoidentityprovider.ForgotPasswordInput) (*cognitoidentityprovider.ForgotPasswordOutput, error)
	GetUser(input *cognitoidentityprovider.GetUserInput) (*cognitoidentityprovider.GetUserOutput, error)
	ConfirmForgotPassword(input *cognitoidentityprovider.ConfirmForgotPasswordInput) (*cognitoidentityprovider.ConfirmForgotPasswordOutput, error)
}

type Cognito struct {
	Client   *cognitoidentityprovider.CognitoIdentityProvider
	UserPool string
}

func NewCognitoClient() (Client, error) {
	cfg := config.InitConfig()
	sess := cfg.AWS.GetAwsSession()

	cognitoClient := cognitoidentityprovider.New(sess)
	return &Cognito{Client: cognitoClient, UserPool: cfg.AWS.COGNITO.USERPOOL_ID}, nil
}

func (c *Cognito) SignUp(input *cognitoidentityprovider.SignUpInput) (*cognitoidentityprovider.SignUpOutput, error) {
	return c.Client.SignUp(input)
}

func (c *Cognito) InitateAuth(input *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	return c.Client.InitiateAuth(input)
}

func (c *Cognito) ConfirmSignUp(input *cognitoidentityprovider.ConfirmSignUpInput) (*cognitoidentityprovider.ConfirmSignUpOutput, error) {
	return c.Client.ConfirmSignUp(input)
}

func (c *Cognito) ForgotPassword(input *cognitoidentityprovider.ForgotPasswordInput) (*cognitoidentityprovider.ForgotPasswordOutput, error) {
	return c.Client.ForgotPassword(input)
}

func (c *Cognito) GetUser(input *cognitoidentityprovider.GetUserInput) (*cognitoidentityprovider.GetUserOutput, error) {
	return c.Client.GetUser(input)
}

func (c *Cognito) ConfirmForgotPassword(input *cognitoidentityprovider.ConfirmForgotPasswordInput) (*cognitoidentityprovider.ConfirmForgotPasswordOutput, error) {
	return c.Client.ConfirmForgotPassword(input)
}
