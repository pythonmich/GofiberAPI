package api

import (
	"FiberFinanceAPI/auth"
	model "FiberFinanceAPI/database/models"
	"github.com/gofiber/fiber/v2"
)

type permissionType string

const (
	// User has admin role
	admin permissionType = "admin"
	// User is logged in (we have their user id)
	member permissionType = "member"
	// User is logged in and user id passed to api is the same
	memberIsTarget permissionType = "memberIsTarget"

	//	anonymous prospects can access the resource allowed for viewing in our server
	prospect permissionType = "prospect"
)

// Admin
var adminOnly = func(role model.UserRole) bool {
	switch role.Role {
	case model.RoleAdmin:
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
var memberIsTargetOnly = func(ctx *fiber.Ctx, userID model.UserID, payload *auth.AccessPayload) bool {
	if userID == "" || payload.SUB == "" {
		return false
	}
	if userID != model.UserID(payload.SUB) {
		return false
	}
	// we are storing the userID in a context after it matches the above conditions and will be available to the request that require the value
	ctx.Locals("userID", userID)
	return true
}

var prospects = func() bool {
	return true
}
