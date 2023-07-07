package http

import (
	"net/http"
	"simple_crud/domain"
	"simple_crud/user/delivery/http/middleware"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	validator "gopkg.in/go-playground/validator.v9"
)

type ResponseError struct {
	Message string `json:"message"`
}

type jwtCustomClaims struct {
	ID int64 `json:"id"`
	jwt.RegisteredClaims
}

type UserHandler struct {
	UUsecase domain.UserUsecase
}

func NewArticleHandler(e *echo.Echo, us domain.UserUsecase) {
	handler := &UserHandler{
		UUsecase: us,
	}
	e.GET("/user", handler.FetchUser, middleware.IsAuthenticated)
	e.POST("/user", handler.Register)
	e.POST("/login", handler.Login)
	e.PUT("/change-password", handler.ChangePassword, middleware.IsAuthenticated)
	e.DELETE("/delete", handler.Delete, middleware.IsAuthenticated)
}

func (u *UserHandler) FetchUser(c echo.Context) error {
	numS := c.QueryParam("num")
	num, _ := strconv.Atoi(numS)
	cursor := c.QueryParam("cursor")
	ctx := c.Request().Context()

	listAr, nextCursor, err := u.UUsecase.Fetch(ctx, cursor, int64(num))
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	c.Response().Header().Set(`X-Cursor`, nextCursor)
	return c.JSON(http.StatusOK, listAr)
}

func (u *UserHandler) Register(c echo.Context) (err error) {
	var user domain.User
	now := time.Now()
	err = c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	user.CreatedAt = now
	user.UpdatedAt = now

	var ok bool
	if ok, err = isRequestValid(&user); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	err = u.UUsecase.Register(ctx, &user)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, http.StatusText(http.StatusCreated))

}

func (u *UserHandler) Login(c echo.Context) (err error) {
	var user domain.User
	err = c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&user); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	res, err := u.UUsecase.Login(ctx, user.Username, user.Password)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(user.Password))
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: domain.ErrLogin.Error()})
	}

	// Set custom claims
	claims := &jwtCustomClaims{
		res.ID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("B0mb45Tic"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func (u *UserHandler) ChangePassword(c echo.Context) (err error) {
	var user domain.User
	err = c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&user); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	err = u.UUsecase.ChangePassword(ctx, user.ID, user.Password)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, http.StatusText(http.StatusOK))

}

func (u *UserHandler) Delete(c echo.Context) (err error) {
	var user domain.User
	err = c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&user); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	err = u.UUsecase.Delete(ctx, user.ID)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, http.StatusText(http.StatusOK))

}

func isRequestValid(u *domain.User) (bool, error) {
	validate := validator.New()
	err := validate.Struct(u)
	if err != nil {
		return false, err
	}
	return true, nil
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
