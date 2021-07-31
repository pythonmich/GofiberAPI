package api

import "github.com/gofiber/fiber/v2"

// registerRoutes manages our routes and middleware
func (s *Server) registerRoutes() {
	s.logs.WithField("func", "router.go -> registerRoutes()").Debug()
	s.routes = fiber.New()
	s.routes.Get("/version", s.version)
	permissions := newPermissions(s.repo, s.logs)

	v1 := s.routes.Group("/api/v1")
	// --------USER--------
	v1.Use(permissions.wrap(prospect))
	v1.Post("/users", s.createUser)
	v1.Post("/login", s.loginUser)

	// ------VERIFICATION REQUIRED ROUTES -----
	v1auth := v1.Use(authTokenMiddleWare(s.token, s.logs))
	v1auth.Get("/users/:userID", permissions.wrap(memberIsTarget, admin), s.getUserByID)
	v1auth.Get("/users", permissions.wrap(admin), s.listUsers)
	v1auth.Put("/users/:userID", permissions.wrap(memberIsTarget, admin), s.changePassword)
	v1auth.Delete("/users/:userID", permissions.wrap(memberIsTarget, admin), s.deleteUser)

	// -------TOKENS--------
	v1auth.Post("/refresh", permissions.wrap(prospect), s.refreshToken)

	//	-------ACCOUNTS-------
	//TODO: Remove admin in future from accounts
	v1auth.Post("/users/:userID/accounts", permissions.wrap(memberIsTarget), s.createAccount)
	v1auth.Get("/users/:userID/accounts/:accountID", permissions.wrap(memberIsTarget), s.getAccount)
	v1auth.Get("/users/:userID/accounts/:accountID/balance", permissions.wrap(memberIsTarget), s.accountBalance)
	v1auth.Get("/users/:userID/accounts", permissions.wrap(memberIsTarget), s.listAccounts)
	v1auth.Delete("/users/:userID/accounts/:accountID", permissions.wrap(memberIsTarget), s.deleteAccount)

	// -----CATEGORY-----
	v1auth.Post("/users/:userID/categories", permissions.wrap(memberIsTarget), s.createCategory)
	v1auth.Get("/users/:userID/categories/:categoryID", permissions.wrap(memberIsTarget), s.getCategory)
	v1auth.Get("/users/:userID/categories", permissions.wrap(memberIsTarget), s.listCategories)
	v1auth.Put("/users/:userID/categories/:categoryID", permissions.wrap(memberIsTarget), s.updateCategory)
	v1auth.Delete("/users/:userID/categories/:categoryID", permissions.wrap(memberIsTarget), s.deleteCategory)

	// -----MERCHANTS-----
	v1auth.Post("/users/:userID/merchants", permissions.wrap(memberIsTarget), s.createMerchant)
	v1auth.Get("/users/:userID/merchants/:merchantID", permissions.wrap(memberIsTarget), s.getMerchant)
	v1auth.Get("/users/:userID/merchants", permissions.wrap(memberIsTarget), s.listMerchants)
	v1auth.Put("/users/:userID/merchants/:merchantID", permissions.wrap(memberIsTarget), s.updateMerchant)
	v1auth.Delete("/users/:userID/merchants/:merchantID", permissions.wrap(memberIsTarget), s.deleteMerchant)

	// -----TRANSACTIONS-----
	v1auth.Post("/users/:userID/transactions", permissions.wrap(memberIsTarget), s.createTransaction)
	v1auth.Get("/users/:userID/transactions/:transactionID", permissions.wrap(memberIsTarget), s.getTransaction)
	v1auth.Get("/users/:userID/transactions", permissions.wrap(memberIsTarget), s.listTransactionsByUserID)
	v1auth.Get("/accounts/:accountID/transactions", permissions.wrap(memberIsTarget), s.listTransactionsByAccountID)
	v1auth.Get("/categories/:categoryID/transactions", permissions.wrap(memberIsTarget), s.listTransactionsByCategoryID)
	v1auth.Put("/users/:userID/transactions/:transactionID", permissions.wrap(memberIsTarget), s.updateTransaction)
	v1auth.Delete("/users/:userID/transactions/:transactionID", permissions.wrap(memberIsTarget), s.deleteTransaction)

	//  ----ADMIN ROLES----
	v1Admin := v1auth.Use(permissions.wrap(admin))
	v1Admin.Post("/users/:userID/role", s.grantRole)
	v1Admin.Delete("/users/:userID/role", s.revokeRole)
	v1Admin.Get("/users/:userID/role", s.getUserRole)
	v1Admin.Get("/users/:userID/roles", s.listRoles)
}
