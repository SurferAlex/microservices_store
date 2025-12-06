DELETE FROM role_permissions
USING roles r, permissions p
WHERE role_permissions.role_id = r.id
  AND role_permissions.permission_id = p.id
  AND r.name IN ('user', 'moderator', 'admin')
  AND p.name IN ('read_profile', 'edit_profile', 'ban_user', 'delete_user', 'manage_system');

DELETE FROM permissions
WHERE name IN ('read_profile', 'edit_profile', 'ban_user', 'delete_user', 'manage_system');

DELETE FROM roles
WHERE name IN ('user', 'moderator', 'admin');