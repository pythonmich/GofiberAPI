package api

import (
	"FiberFinanceAPI/auth"
	db "FiberFinanceAPI/database/sqlc"
)

type permissionType string

const (
	// User has admin role
	admin permissionType = "admin"
	// User is logged in (we have their user id)
	member permissionType = "member"
	// User is logged in and user id passed to api is the same
	memberIsTarget permissionType = "memberIsTarget"
)

// Admin
var adminOnly = func(role *db.UserRole) bool {
	switch role.Role {
	case db.RoleAdmin:
		return true
	}
	return false
}

// logged in User
// auth.AccessPayload value shall be extracted from context.Locals
var memberOnly = func(payload *auth.AccessPayload) bool {
	return payload.SUB != ""
}

// logged in user == target user
var memberIsTargetOnly = func(userID db.UserID, payload *auth.AccessPayload) bool {
	if userID == "" || payload.SUB == "" {
		return false
	}
	if userID != db.UserID(payload.SUB) {
		return false
	}
	return true
}
