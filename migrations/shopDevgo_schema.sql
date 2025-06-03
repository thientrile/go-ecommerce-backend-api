-- MySQL dump 10.13  Distrib 8.0.42, for Linux (x86_64)
--
-- Host: localhost    Database: shopDevgo
-- ------------------------------------------------------
-- Server version	8.0.42

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Current Database: `shopDevgo`
--

/*!40000 DROP DATABASE IF EXISTS `shopDevgo`*/;

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `shopDevgo` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;

USE `shopDevgo`;

--
-- Table structure for table `go_crm_user`
--

DROP TABLE IF EXISTS `go_crm_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `go_crm_user` (
  `usr_id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'Account ID',
  `usr_email` varchar(30) NOT NULL DEFAULT '' COMMENT 'Email',
  `usr_phone` varchar(15) NOT NULL DEFAULT '' COMMENT 'Phone Number',
  `usr_username` varchar(30) NOT NULL DEFAULT '' COMMENT 'Username',
  `usr_password` varchar(32) NOT NULL DEFAULT '' COMMENT 'Password',
  `usr_created_at` int NOT NULL DEFAULT '0' COMMENT 'Creation Time',
  `usr_updated_at` int NOT NULL DEFAULT '0' COMMENT 'Update Time',
  `usr_create_ip_at` varchar(12) NOT NULL DEFAULT '' COMMENT 'Creation IP',
  `usr_last_login_at` int NOT NULL DEFAULT '0' COMMENT 'Last Login Time',
  `usr_last_login_ip_at` varchar(12) NOT NULL DEFAULT '' COMMENT 'Last Login IP',
  `usr_login_times` int NOT NULL DEFAULT '0' COMMENT 'Login Times',
  `usr_status` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'Status 1:enable, 0:disable, -1:deleted',
  PRIMARY KEY (`usr_id`),
  KEY `idx_email` (`usr_email`),
  KEY `idx_phone` (`usr_phone`),
  KEY `idx_username` (`usr_username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Account';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `go_crm_userV2`
--

DROP TABLE IF EXISTS `go_crm_userV2`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `go_crm_userV2` (
  `usr_id` int NOT NULL AUTO_INCREMENT COMMENT 'Account ID',
  `usr_email` longtext COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Email',
  `usr_phone` longtext COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Phone Number',
  `usr_username` longtext COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Username',
  `usr_password` longtext COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Password',
  `usr_created_at` int NOT NULL COMMENT 'Creation Time',
  `usr_updated_at` int NOT NULL COMMENT 'Update Time',
  `usr_create_ip_at` longtext COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Creation IP',
  `usr_last_login_at` int NOT NULL COMMENT 'Last Login Time',
  `usr_last_login_ip_at` longtext COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Last Login IP',
  `usr_login_times` int NOT NULL COMMENT 'Login Times',
  `usr_status` tinyint(1) NOT NULL COMMENT 'Status 1:enable, 0:disable, -1:deleted',
  PRIMARY KEY (`usr_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `go_db_role`
--

DROP TABLE IF EXISTS `go_db_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `go_db_role` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT ' ''Role ID''',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `role_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT ' ''Role Name''',
  `role_note:type:text` longtext COLLATE utf8mb4_unicode_ci,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_go_db_role_role_name` (`role_name`),
  KEY `idx_go_db_role_deleted_at` (`deleted_at`),
  KEY `idx_role_name` (`role_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `go_db_user`
--

DROP TABLE IF EXISTS `go_db_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `go_db_user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `uuid` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `username` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `is_active` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_go_db_user_uuid` (`uuid`),
  UNIQUE KEY `uni_go_db_user_username` (`username`),
  KEY `idx_go_db_user_deleted_at` (`deleted_at`),
  KEY `idx_uuid` (`uuid`),
  KEY `idx_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `go_user_roles`
--

DROP TABLE IF EXISTS `go_user_roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `go_user_roles` (
  `user_id` bigint unsigned NOT NULL,
  `role_id` int NOT NULL COMMENT ' ''Role ID''',
  PRIMARY KEY (`user_id`,`role_id`),
  KEY `fk_go_user_roles_role` (`role_id`),
  CONSTRAINT `fk_go_user_roles_role` FOREIGN KEY (`role_id`) REFERENCES `go_db_role` (`id`),
  CONSTRAINT `fk_go_user_roles_user` FOREIGN KEY (`user_id`) REFERENCES `go_db_user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-06-03  6:33:48
