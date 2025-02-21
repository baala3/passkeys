package controller

import (
	"context"
	"net/http"
	"os"
	"strings"
	"testing"

	"net/http/httptest"

	"github.com/alexedwards/argon2id"
	"github.com/alicebob/miniredis/v2"
	"github.com/baala3/passkeys/db"
	"github.com/baala3/passkeys/model"
	"github.com/baala3/passkeys/pkg"
	"github.com/baala3/passkeys/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
)

var (
	database *bun.DB
	e *echo.Echo
	miniRedis *miniredis.Miniredis
	userRepository repository.UserRepository
	sessionRepository repository.SessionRepository
	passwordController *PasswordController
)

func setup() {
	database = db.GetTestDB()
	e = echo.New()
	miniRedis = miniredis.NewMiniRedis()
	err := miniRedis.StartAddr(":16379")
	if err != nil {
		panic(err)
	}
	userRepository = repository.UserRepository{DB: database}
	sessionRepository = repository.SessionRepository{RedisClient: pkg.GetRedisClient()}
	passwordController = &PasswordController{
		UserRepository: userRepository,
		SessionRepository: sessionRepository,
	}
	loadFixtures()
}

func loadFixtures() {
	password := "password123"
	passwordHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		panic(err)
	}
	_, err = database.NewInsert().Model(&model.User{
		ID: uuid.New(),
		Email: "existing@email.com",
		PasswordHash: passwordHash,
	}).Exec(context.Background())
	if err != nil {
		panic(err)
	}
}

func teardown() {
	miniRedis.Close()
	_ =database.Close()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestPasswordController_SignUp(t *testing.T) {
	t.Run("invalid email", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/signup", strings.NewReader(`{"email":"invalid","password":"password123"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		assert.NoError(t, passwordController.SignUp()(ctx))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"status": "error", "errorMessage":"Invalid email"}`, rec.Body.String())
	})

	t.Run("password too short", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/signup", strings.NewReader(`{"email":"test@email.com", "password":"short"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		assert.NoError(t, passwordController.SignUp()(ctx))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"status": "error", "errorMessage":"Password must be at least 8 characters"}`, rec.Body.String())
	})

	t.Run("account already exists", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/signup", strings.NewReader(`{"email":"existing@email.com", "password":"password123"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		assert.NoError(t, passwordController.SignUp()(ctx))
		assert.Equal(t, http.StatusConflict, rec.Code)
		assert.JSONEq(t, `{"status": "error", "errorMessage":"An account with that email already exists."}`, rec.Body.String())
	})

	t.Run("successful sign up", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/signup", strings.NewReader(`{"email":"new@email.com", "password":"password123"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		assert.NoError(t, passwordController.SignUp()(ctx))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"status": "ok", "errorMessage":""}`, rec.Body.String())

		// Creates the user in the user repository
		user, err := userRepository.FindUserByEmail(ctx.Request().Context(), "new@email.com")
		if err != nil {
			t.Error(err)
		}
		assert.NotNil(t, user)

		// Creates a session
		cookie := rec.Result().Cookies()[0]
		assert.Equal(t, "auth", cookie.Name)
		value, err := miniRedis.Get(cookie.Value)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, user.ID.String(), value)
	})
}

func TestPasswordController_Login(t *testing.T) {
	t.Run("invalid email", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/login", strings.NewReader(`{"email":"invalid","password":"password123"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		assert.NoError(t, passwordController.Login()(ctx))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"status": "error", "errorMessage":"Invalid email"}`, rec.Body.String())
	})

	t.Run("account does not exist", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/login", strings.NewReader(`{"email":"notexisting@email.com", "password":"password123"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		assert.NoError(t, passwordController.Login()(ctx))
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.JSONEq(t, `{"status": "error", "errorMessage":"An account with that email does not exist."}`, rec.Body.String())
	})

	t.Run("incorrect password", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/login", strings.NewReader(`{"email":"existing@email.com", "password":"wrongPassword"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		assert.NoError(t, passwordController.Login()(ctx))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.JSONEq(t, `{"status": "error", "errorMessage":"Invalid password."}`, rec.Body.String())
	})

	t.Run("successful login", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/login", strings.NewReader(`{"email":"existing@email.com", "password":"password123"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		assert.NoError(t, passwordController.Login()(ctx))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"status": "ok", "errorMessage":""}`, rec.Body.String())

		// Creates a session
		cookie := rec.Result().Cookies()[0]
		assert.Equal(t, "auth", cookie.Name)
		value, err := miniRedis.Get(cookie.Value)
		if err != nil {
			t.Error(err)
		}
		user, err := userRepository.FindUserByEmail(ctx.Request().Context(), "existing@email.com")
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, user.ID.String(), value)
	})
}

func TestPasswordController_Logout(t *testing.T) {
	t.Run("not logged in", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/logout", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		assert.NoError(t, passwordController.Logout()(ctx))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.JSONEq(t, `{"status": "error", "errorMessage":"Not logged in."}`, rec.Body.String())
	})

	t.Run("successful logout", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/login", strings.NewReader(`{"email":"existing@email.com", "password":"password123"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		assert.NoError(t, passwordController.Login()(ctx))
		cookie := rec.Result().Cookies()[0]

		req = httptest.NewRequest(echo.POST, "/logout", nil)
		req.AddCookie(cookie)
		rec = httptest.NewRecorder()
		ctx = e.NewContext(req, rec)
		assert.NoError(t, passwordController.Logout()(ctx))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"status": "ok", "errorMessage":""}`, rec.Body.String())

		// It removes session information from the session store
		_, err := miniRedis.Get(cookie.Value)
		assert.Error(t, err)
	})
}
