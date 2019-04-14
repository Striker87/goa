package controllers

import (
	"database/sql"
	"github.com/goadesign/goa"
	"github.com/rymccue/golang-auth-microservice/app"
	"goa/models"
	"goa/utils/crypto"
	"goa/utils/jwt"
)

// AuthenticationController implements the authentication resource.
type AuthenticationController struct {
	*goa.Controller
	*sql.DB
}

// NewAuthenticationController creates a authentication controller.
func NewAuthenticationController(service *goa.Service, db *sql.DB) *AuthenticationController {
	return &AuthenticationController{
		Controller: service.NewController("AuthenticationController"),
		DB:         db,
	}
}

// Login runs the login action
func (c *AuthenticationController) Login(ctx *app.LoginAuthenticationContext) error {
	payload := ctx.Payload
	u, err := models.GetUserByEmail(c.DB, payload.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.BadRequest(goa.ErrBadRequest("Invalid email or password"))
		}
		c.Service.LogError("Login user", "err", err)
		return ctx.InternalServerError()
	}

	hashedPassword := crypto.HashPassword(payload.Password, *u.Salt)
	if *u.Password != hashedPassword {
		return ctx.BadRequest(goa.ErrBadRequest("Invalid email or password"))
	}

	token, err := jwt.CreateJWTToken(*u.Email)
	if err != nil {
		c.Service.LogError("Login user", "err", err)
		return ctx.InternalServerError()
	}

	res := &app.Token {
		Token: &token,
	}
	return ctx.OK(res)
}

// Register runs the register action.
func(c *AuthenticationController) Register(ctx *app.RegisterAuthenticationContext) error {
	payload := ctx.Payload
	exists, err := models.CheckEmailExists(c.DB, payload.Email)
	if err != nil {
		c.Service.LogError("Register user", "err", err)
		return ctx.InternalServerError()
	}

	if exists {
		return ctx.BadRequest(goa.ErrBadRequest("Email already exists"))
	}

	err = models.AddUserToDatabase(c.DB, payload.FirstName, payload.LastName, payload.Email, payload.Password)
	if err != nil {
		c.Service.LogError("Register user", "err", err)
		return ctx.InternalServerError()
	}

	token, err := jwt.CreateJWTToken(payload.Email)
	if err != nil {
		c.Service.LogError("Failed to create token", "err", err)
		return ctx.InternalServerError()
	}

	res := &app.Token{Token: &token}
	return ctx.OK(res)
}