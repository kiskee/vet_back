INSERT INTO vets (user_id, status, max_concurrent)
SELECT id, 'available', 1
FROM users
WHERE role = 'veterinarian'
  AND id NOT IN (SELECT user_id FROM vets)
ON CONFLICT (user_id) DO NOTHING;
