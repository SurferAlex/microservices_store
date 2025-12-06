package roles

import (
	_ "github.com/lib/pq"
)

type PermissionSet map[string]struct{}

var rolePermissions = map[string]PermissionSet{
	"user":      {"read_profile": {}, "edit_profile": {}},
	"moderator": {"read_profile": {}, "edit_profile": {}, "ban_user": {}},
	"admin":     {"read_profile": {}, "edit_profile": {}, "ban_user": {}, "delete_user": {}, "manage_system": {}},
}

func RolePermissions() map[string]PermissionSet {
	return rolePermissions
}

func HasPermission(role, permission string) bool {
	perms, ok := rolePermissions[role]
	if !ok {
		return false
	}
	_, ok = perms[permission]
	return ok
}
