-- name: InsertOTPVerify :one
INSERT INTO user_verify (
    verify_key,
    verify_hash_key,
    verify_otp,
    verify_type
)
VALUES (
    $1, $2, $3, $4
)
RETURNING id;

-- name: CheckOTPVerifyExist :one
SELECT EXISTS (
    SELECT 1 FROM user_verify
    WHERE verify_key = $1 AND is_deleted = false
);


-- name: UpdateUserVerificationStatus :exec
UPDATE user_verify
SET is_verified = true,
    updated_at = now()
WHERE verify_hash_key = $1;

-- name: GetInfoOTP :one
SELECT verify_otp, verify_key, verify_hash_key, verify_type, is_verified, is_deleted
FROM user_verify
WHERE verify_hash_key = $1;
