-- name: EnableTwoFactorTypeEmail :exec
INSERT INTO pre_go_acc_user_two_factor_9999 
  (user_id, two_factor_auth_type, two_factor_email, two_factor_auth_secret, two_factor_is_active, two_factor_created_at, two_factor_updated_at)
VALUES (?, ?, ?, "OTP", FALSE, NOW(), NOW());

-- name: EnableTwoFactorTypeSMS :exec
INSERT INTO pre_go_acc_user_two_factor_9999 
  (user_id, two_factor_auth_type, two_factor_phone, two_factor_auth_secret, two_factor_is_active, two_factor_created_at, two_factor_updated_at)
VALUES (?, ?, ?, "OTP", FALSE, NOW(), NOW());

-- name: DisableTwoFactor :exec
UPDATE pre_go_acc_user_two_factor_9999
SET two_factor_is_active = 0,
    two_factor_updated_at = NOW()
WHERE user_id = ? AND two_factor_auth_type = ?;

-- name: UpdateTwoFactorStatusVerification :exec
UPDATE pre_go_acc_user_two_factor_9999
SET two_factor_is_active = 1, 
    two_factor_updated_at = NOW()
WHERE user_id = ? AND two_factor_auth_type = ? AND two_factor_is_active = 0;

-- name: VerifyTwoFactor :one
SELECT COUNT(*)
FROM pre_go_acc_user_two_factor_9999
WHERE user_id = ? AND two_factor_auth_type = ? AND two_factor_is_active = 1
LIMIT 1;

-- name: GetTwoFactorStatus :one
SELECT two_factor_is_active
FROM pre_go_acc_user_two_factor_9999
WHERE user_id = ? AND two_factor_auth_type = ?
LIMIT 1;

-- name: IsTwoFactorEnabled :one
SELECT COUNT(*)
FROM pre_go_acc_user_two_factor_9999
WHERE user_id = ? AND two_factor_is_active = 1;

-- name: AddOrUpdateTwoFactor :exec
INSERT INTO pre_go_acc_user_two_factor_9999 (user_id, two_factor_phone, two_factor_is_active)
VALUES (?, ?, TRUE)
ON DUPLICATE KEY UPDATE
    two_factor_phone = VALUES(two_factor_phone),
    two_factor_updated_at = NOW();

-- name: AddOrUpdateEmail :exec
INSERT INTO pre_go_acc_user_two_factor_9999 (user_id, two_factor_email, two_factor_is_active)
VALUES (?, ?, TRUE)
ON DUPLICATE KEY UPDATE
    two_factor_email = VALUES(two_factor_email),
    two_factor_updated_at = NOW();

-- name: GetUserFactorMethods :many
SELECT two_factor_id, user_id, two_factor_auth_type, two_factor_auth_secret,
       two_factor_phone, two_factor_email, two_factor_is_active,
       two_factor_created_at, two_factor_updated_at
FROM pre_go_acc_user_two_factor_9999
WHERE user_id = ?;

-- name: ReactivateTwoFactor :exec
UPDATE pre_go_acc_user_two_factor_9999
SET two_factor_is_active = 1, 
    two_factor_updated_at = NOW()
WHERE user_id = ? AND two_factor_auth_type = ?;

-- name: RemoveTwoFactor :exec
DELETE FROM pre_go_acc_user_two_factor_9999
WHERE user_id = ? AND two_factor_auth_type = ?;

-- name: CountActiveTwoFactorMethods :one
SELECT COUNT(*)
FROM pre_go_acc_user_two_factor_9999
WHERE user_id = ? AND two_factor_is_active = 1;

-- name: GetTwoFactorMethodsByID :one
SELECT two_factor_id, user_id, two_factor_auth_type, two_factor_auth_secret,
       two_factor_phone, two_factor_email, two_factor_is_active,
       two_factor_created_at, two_factor_updated_at
FROM pre_go_acc_user_two_factor_9999
WHERE two_factor_id = ?
LIMIT 1;

-- name: GetTwoFactorMethodsByIDAndType :one
SELECT two_factor_id, user_id, two_factor_auth_type, two_factor_auth_secret,
       two_factor_phone, two_factor_email, two_factor_is_active,
       two_factor_created_at, two_factor_updated_at
FROM pre_go_acc_user_two_factor_9999
WHERE two_factor_id = ? AND two_factor_auth_type = ?
LIMIT 1;
