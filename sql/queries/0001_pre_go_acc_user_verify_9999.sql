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

