package api

import (
	"fmt"

	db "github.com/RahmounFares2001/simple-bank-backend/db/sqlc"
	"github.com/RahmounFares2001/simple-bank-backend/token"
	"github.com/RahmounFares2001/simple-bank-backend/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config util.Config
	store  *db.Store
	token  token.Maker
	router *gin.Engine
}

// creates new HTTP server + setup routing
func NewServer(config util.Config, store *db.Store) (*Server, error) {
	// new paseto maker
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)

	if err != nil {
		return nil, fmt.Errorf("cannot create token maker %w", err)
	}

	server := &Server{
		config: config,
		store:  store,
		token:  tokenMaker,
	}

	// currency validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

// setup router
func (server *Server) setupRouter() {
	router := gin.Default()

	// user
	router.POST("/users", server.createUser)
	// login
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddlware(server.token))

	// create account
	authRoutes.POST("/accounts", server.createAccount)
	// get account
	authRoutes.GET("/accounts/:id", server.getAccount)
	// get accounts from query
	authRoutes.GET("/accounts", server.listAccounts)
	// delete account
	authRoutes.DELETE("/accounts/:id", server.deleteAccount)

	// transfer
	authRoutes.POST("/transfers", server.createTransfer)
	// get transfer
	authRoutes.GET("/transfers/:id", server.getTransfer)
	// get list transfers
	authRoutes.GET("/transfers", server.getListTransfers)

	// create entry
	authRoutes.POST("/entries", server.createEntry)
	// get entry
	authRoutes.GET("/entries/:id", server.getEntry)
	// get list entries
	authRoutes.GET("/entries", server.getListEntries)

	server.router = router
}

// runs the http server on a specific address
func (server *Server) Start(adress string) error {
	return server.router.Run(adress)
}

// gin.H : shortcut of map[]string interface{}
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
