INSERT INTO roles (name) VALUES
  ('user'),
  ('moderator'),
  ('admin')
ON CONFLICT (name) DO NOTHING;

INSERT INTO permissions (name) VALUES
  ('read_profile'),
  ('edit_profile'),
  ('ban_user'),
  ('delete_user'),
  ('manage_system')
ON CONFLICT (name) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON
  (r.name = 'user'      AND p.name IN ('read_profile', 'edit_profile')) OR
  (r.name = 'moderator' AND p.name IN ('read_profile', 'edit_profile', 'ban_user')) OR
  (r.name = 'admin'     AND p.name IN ('read_profile', 'edit_profile', 'ban_user', 'delete_user', 'manage_system'))
ON CONFLICT DO NOTHING;