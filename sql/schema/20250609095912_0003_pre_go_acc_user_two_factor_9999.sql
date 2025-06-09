-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS `pre_go_acc_user_two_factor_9999`;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE `pre_go_acc_user_two_factor_9999` (
  `two_factor_id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  `user_id` INT NOT NULL,
  `two_factor_auth_type` ENUM('SMS', 'EMAIL', 'APP') NOT NULL,
  `two_factor_auth_secret` VARCHAR(200) NOT NULL,
  `two_factor_phone` VARCHAR(20) DEFAULT NULL,
  `two_factor_email` VARCHAR(255) DEFAULT NULL,
  `two_factor_is_active` BOOLEAN NOT NULL DEFAULT TRUE,
  `two_factor_created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `two_factor_updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
--   CONSTRAINT `fk_user_id_two_factor` FOREIGN KEY (`user_id`) 
--     REFERENCES `pre_go_acc_user_base_9999` (`user_id`) 
--     ON DELETE CASCADE,
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_auth_type` (`two_factor_auth_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 
  COMMENT='pre_go_acc_user_two_factor_9999';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `pre_go_acc_user_two_factor_9999`;
-- +goose StatementEnd
