-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS `pre_go_acc_user_9999`;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE `pre_go_acc_user_9999` (
    `user_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT 'User ID',
    `user_account` VARCHAR(255) NOT NULL COMMENT 'User account',
    `user_nickname` VARCHAR(255) DEFAULT NULL COMMENT 'User nickname',
    `user_avatar` VARCHAR(255) DEFAULT NULL COMMENT 'User avatar',
    `user_state` TINYINT UNSIGNED NOT NULL COMMENT 'User state: 0-Locked ,  1-Activated ,  2-Not Activated',
    `user_mobile` VARCHAR(20) DEFAULT NULL COMMENT 'Mobile phone number',
    `user_gender` TINYINT UNSIGNED DEFAULT NULL COMMENT 'User gender: 0-Secret ,  1-Male ,  2-Female',
    `user_birthday` DATE DEFAULT NULL COMMENT 'User birthday',
    `user_email` VARCHAR(255) DEFAULT NULL COMMENT 'User email address',
    `user_is_authentication` TINYINT UNSIGNED NOT NULL COMMENT 'Authentication status: 0-Not Authenticated ,  1-Pending ,  2-Authenticated ,  3-Failed',
    `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record creation time',
    `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record update time'
);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE `pre_go_acc_user_9999` 
    ADD UNIQUE `pre_go_acc_user_9999_user_account_unique`(`user_account`),
    ADD INDEX `pre_go_acc_user_9999_user_state_index`(`user_state`),
    ADD INDEX `pre_go_acc_user_9999_user_mobile_index`(`user_mobile`),
    ADD INDEX `pre_go_acc_user_9999_user_email_index`(`user_email`),
    ADD INDEX `pre_go_acc_user_9999_user_is_authentication_index`(`user_is_authentication`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `pre_go_acc_user_9999`;
-- +goose StatementEnd
