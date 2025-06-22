-- name: GetValidOTP :one
SELECT verify_otp,verify_key_hash,verify_key,verify_id 
FROM `pre_go_acc_user_verify_9999`
WHERE `verify_key_hash` = ? AND `is_verified` = 0 
LIMIT 1;

-- update lại
-- name: UpdateUserVerificationStatus :exec
UPDATE `pre_go_acc_user_verify_9999` 
SET `is_verified` = 1, `verify_updated_at` = NOW() 
WHERE `verify_key_hash` = ? AND `is_verified` = 0;


-- name: InsertOTPVerify :execresult   
-- Insert a new OTP verification record
INSERT INTO `pre_go_acc_user_verify_9999` (
    verify_key_hash,
    verify_otp,
    verify_key,
    verify_type,
    verify_updated_at
) VALUES (?, ?, ?, ?, NOW())
ON DUPLICATE KEY UPDATE
    verify_otp = VALUES(verify_otp),
    verify_key = VALUES(verify_key),
    verify_type = VALUES(verify_type),
    is_verified = CAST(0 AS UNSIGNED),
    is_deleted = CAST(0 AS UNSIGNED),
    verify_updated_at = NOW();


-- name: GetInfoOTP :one
SELECT verify_id, verify_otp,verify_key,verify_key_hash,verify_type,is_verified,is_deleted,verify_created_at 
FROM `pre_go_acc_user_verify_9999` 
WHERE `verify_key_hash` = ? AND `is_deleted` = 0
LIMIT 1;

-- name: DeleteVerifyByID :exec
-- Delete OTP verification record by verify_id (soft delete)
UPDATE `pre_go_acc_user_verify_9999` 
SET `is_deleted` = 1, `verify_updated_at` = NOW() 
WHERE `verify_id` = ?;

-- name: DeleteVerifyByIDHard :exec
-- Hard delete OTP verification record by verify_id
DELETE FROM `pre_go_acc_user_verify_9999` 
WHERE `verify_id` = ?;

-- name: DeleteExpiredVerifications :exec
-- Delete all expired verification records (older than 24 hours)
UPDATE `pre_go_acc_user_verify_9999` 
SET `is_deleted` = 1, `verify_updated_at` = NOW() 
WHERE `verify_created_at` < DATE_SUB(NOW(), INTERVAL 24 HOUR) 
AND `is_deleted` = 0;

-- name: DeleteVerifiedRecords :exec
-- Clean up verified records older than 7 days
DELETE FROM `pre_go_acc_user_verify_9999` 
WHERE `is_verified` = 1 
AND `verify_updated_at` < DATE_SUB(NOW(), INTERVAL 7 DAY);

-- name: GetVerifyByID :one
-- Get verification record by verify_id
SELECT verify_id, verify_otp, verify_key, verify_key_hash, verify_type, is_verified, is_deleted, verify_created_at, verify_updated_at
FROM `pre_go_acc_user_verify_9999` 
WHERE `verify_id` = ? AND `is_deleted` = 0
LIMIT 1;

-- name: DeleteVerifyByKeyHash :exec
-- Delete OTP verification record by verify_key_hash (soft delete)
UPDATE `pre_go_acc_user_verify_9999` 
SET `is_deleted` = 1, `verify_updated_at` = NOW() 
WHERE `verify_key_hash` = ? AND `is_deleted` = 0;

-- name: DeleteVerifyByKeyHashHard :exec
-- Hard delete OTP verification record by verify_key_hash
DELETE FROM `pre_go_acc_user_verify_9999` 
WHERE `verify_key_hash` = ?;

-- name: DeleteUnverifiedByKeyHash :exec
-- Delete only unverified records by verify_key_hash (soft delete)
UPDATE `pre_go_acc_user_verify_9999` 
SET `is_deleted` = 1, `verify_updated_at` = NOW() 
WHERE `verify_key_hash` = ? AND `is_verified` = 0 AND `is_deleted` = 0;

-- name: DeleteExpiredByKeyHash :exec
-- Delete expired verification by verify_key_hash (older than specified minutes)
UPDATE `pre_go_acc_user_verify_9999` 
SET `is_deleted` = 1, `verify_updated_at` = NOW() 
WHERE `verify_key_hash` = ? 
AND `verify_created_at` < DATE_SUB(NOW(), INTERVAL ? MINUTE) 
AND `is_deleted` = 0;


