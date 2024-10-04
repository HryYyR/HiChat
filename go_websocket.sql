/*
 Navicat Premium Data Transfer

 Source Server         : index
 Source Server Type    : MySQL
 Source Server Version : 80017
 Source Host           : localhost:3306
 Source Schema         : go_websocket

 Target Server Type    : MySQL
 Target Server Version : 80017
 File Encoding         : 65001

 Date: 23/06/2024 14:19:09
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;


-- ----------------------------
-- Table structure for group
-- ----------------------------
DROP TABLE IF EXISTS `group`;
CREATE TABLE `group`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uuid` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `creater_id` int(11) NOT NULL,
  `creater_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `group_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `grade` int(11) NULL DEFAULT 1,
  `member_count` int(11) NULL DEFAULT NULL,
  `unread_message` int(11) NULL DEFAULT NULL,
  `created_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `UQE_group_uuid`(`uuid` ASC) USING BTREE,
  UNIQUE INDEX `UQE_group_group_name`(`group_name` ASC) USING BTREE,
  INDEX `IDX_group_id`(`id` ASC) USING BTREE,
  INDEX `group_ibfk_1`(`unread_message` ASC) USING BTREE,
  INDEX `creater_id`(`creater_id` ASC, `creater_name` ASC) USING BTREE,
  CONSTRAINT `group_ibfk_2` FOREIGN KEY (`creater_id`, `creater_name`) REFERENCES `users` (`id`, `user_name`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 234 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;



-- ----------------------------
-- Table structure for apply_add_user
-- ----------------------------
DROP TABLE IF EXISTS `apply_add_user`;
CREATE TABLE `apply_add_user`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pre_apply_user_id` int(11) NOT NULL,
  `pre_apply_user_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `apply_user_id` int(11) NOT NULL,
  `apply_user_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `apply_msg` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `apply_way` int(11) NULL DEFAULT 1,
  `handle_status` int(11) NULL DEFAULT 0,
  `created_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `IDX_apply_add_user_id`(`id` ASC) USING BTREE,
  INDEX `pre_apply_user_id`(`pre_apply_user_id` ASC, `pre_apply_user_name` ASC) USING BTREE,
  INDEX `apply_user_id`(`apply_user_id` ASC, `apply_user_name` ASC) USING BTREE,
  CONSTRAINT `apply_add_user_ibfk_1` FOREIGN KEY (`pre_apply_user_id`, `pre_apply_user_name`) REFERENCES `users` (`id`, `user_name`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `apply_add_user_ibfk_2` FOREIGN KEY (`apply_user_id`, `apply_user_name`) REFERENCES `users` (`id`, `user_name`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 66 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for apply_join_group
-- ----------------------------
DROP TABLE IF EXISTS `apply_join_group`;
CREATE TABLE `apply_join_group`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `apply_user_id` int(11) NOT NULL,
  `apply_user_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `group_id` int(11) NOT NULL,
  `apply_msg` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `apply_way` int(11) NULL DEFAULT 1,
  `handle_status` int(11) NULL DEFAULT 0,
  `created_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `IDX_apply_join_group_id`(`id` ASC) USING BTREE,
  INDEX `apply_user_id`(`apply_user_id` ASC, `apply_user_name` ASC) USING BTREE,
  INDEX `group_id`(`group_id` ASC) USING BTREE,
  CONSTRAINT `apply_join_group_ibfk_1` FOREIGN KEY (`apply_user_id`, `apply_user_name`) REFERENCES `users` (`id`, `user_name`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `apply_join_group_ibfk_2` FOREIGN KEY (`group_id`) REFERENCES `group` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 439 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for group_message
-- ----------------------------
DROP TABLE IF EXISTS `group_message`;
CREATE TABLE `group_message`  (
  `id` int(128) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `user_uuid` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `user_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `user_avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `user_city` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `user_age` int(11) NULL DEFAULT NULL,
  `group_id` int(11) NOT NULL,
  `msg` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `msg_type` int(11) NULL DEFAULT NULL,
  `is_reply` tinyint(1) NULL DEFAULT NULL,
  `reply_user_id` int(11) NULL DEFAULT NULL,
  `context` blob NULL,
  `created_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `IDX_group_message_user_id`(`user_id` ASC) USING BTREE,
  INDEX `IDX_group_message_group_id`(`group_id` ASC) USING BTREE,
  INDEX `user_id`(`user_id` ASC, `user_name` ASC, `user_avatar` ASC, `user_city` ASC, `user_age` ASC) USING BTREE,
  CONSTRAINT `group_message_ibfk_1` FOREIGN KEY (`user_id`, `user_name`, `user_avatar`, `user_city`, `user_age`) REFERENCES `users` (`id`, `user_name`, `avatar`, `city`, `age`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `group_message_ibfk_2` FOREIGN KEY (`group_id`) REFERENCES `group` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 1697 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for group_unread_message
-- ----------------------------
DROP TABLE IF EXISTS `group_unread_message`;
CREATE TABLE `group_unread_message`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `user_id` int(11) NOT NULL,
  `group_id` int(11) NOT NULL,
  `unread_number` int(11) NOT NULL DEFAULT 0,
  `created_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `IDX_group_unread_message_user_id`(`user_id` ASC) USING BTREE,
  INDEX `unread_number`(`unread_number` ASC) USING BTREE,
  INDEX `group_unread_message_ibfk_2`(`user_id` ASC, `user_name` ASC) USING BTREE,
  INDEX `group_unread_message_ibfk_1`(`group_id` ASC) USING BTREE,
  CONSTRAINT `group_unread_message_ibfk_1` FOREIGN KEY (`group_id`) REFERENCES `group` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 173 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for group_user_relative
-- ----------------------------
DROP TABLE IF EXISTS `group_user_relative`;
CREATE TABLE `group_user_relative`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `group_id` int(11) NOT NULL,
  `group_uuid` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `IDX_group_user_relative_id`(`id` ASC) USING BTREE,
  INDEX `group_user_relative_ibfk_2`(`user_id` ASC) USING BTREE,
  INDEX `group_id`(`group_id` ASC) USING BTREE,
  CONSTRAINT `group_user_relative_ibfk_2` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `group_user_relative_ibfk_3` FOREIGN KEY (`group_id`) REFERENCES `group` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 742 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_message
-- ----------------------------
DROP TABLE IF EXISTS `user_message`;
CREATE TABLE `user_message`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uuid` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `user_id` int(11) NOT NULL,
  `user_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `user_avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `receive_user_id` int(11) NOT NULL,
  `receive_user_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `receive_user_avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `msg` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `msg_type` int(11) NULL DEFAULT NULL,
  `is_reply` tinyint(1) NULL DEFAULT NULL,
  `reply_user_id` int(11) NULL DEFAULT NULL,
  `context` blob NULL,
  `created_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `IDX_user_message_id`(`id` ASC) USING BTREE,
  INDEX `user_message_ibfk_1`(`user_id` ASC, `user_name` ASC, `user_avatar` ASC) USING BTREE,
  INDEX `receive_user_id`(`receive_user_id` ASC, `receive_user_name` ASC, `receive_user_avatar` ASC) USING BTREE,
  CONSTRAINT `user_message_ibfk_1` FOREIGN KEY (`user_id`, `user_name`, `user_avatar`) REFERENCES `users` (`id`, `user_name`, `avatar`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `user_message_ibfk_2` FOREIGN KEY (`receive_user_id`, `receive_user_name`, `receive_user_avatar`) REFERENCES `users` (`id`, `user_name`, `avatar`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 13569 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_unread_message
-- ----------------------------
DROP TABLE IF EXISTS `user_unread_message`;
CREATE TABLE `user_unread_message`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `user_id` int(11) NOT NULL,
  `friend_id` int(11) NOT NULL,
  `unread_number` int(11) NOT NULL DEFAULT 0,
  `created_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `IDX_user_unread_message_uid_fid_unique`(`user_id` ASC, `friend_id` ASC) USING BTREE,
  INDEX `user_id`(`user_id` ASC, `user_name` ASC) USING BTREE,
  INDEX `friend_id`(`friend_id` ASC) USING BTREE,
  CONSTRAINT `user_unread_message_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `user_unread_message_ibfk_2` FOREIGN KEY (`friend_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 41 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_user_relative
-- ----------------------------
DROP TABLE IF EXISTS `user_user_relative`;
CREATE TABLE `user_user_relative`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pre_user_id` int(11) NOT NULL,
  `pre_user_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `back_user_id` int(11) NOT NULL,
  `back_user_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `IDX_user_user_relative_id`(`id` ASC) USING BTREE,
  INDEX `b`(`pre_user_id` ASC, `pre_user_name` ASC) USING BTREE,
  INDEX `c`(`back_user_id` ASC, `back_user_name` ASC) USING BTREE,
  CONSTRAINT `user_user_relative_ibfk_1` FOREIGN KEY (`pre_user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `user_user_relative_ibfk_2` FOREIGN KEY (`back_user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 24 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for users_file
-- ----------------------------
DROP TABLE IF EXISTS `users_file`;
CREATE TABLE `users_file`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `identity` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `hash` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `ext` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `size` bigint(20) NULL DEFAULT NULL,
  `path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `created_at` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `UQE_users_file_identity`(`identity` ASC) USING BTREE,
  UNIQUE INDEX `UQE_users_file_hash`(`hash` ASC) USING BTREE,
  INDEX `IDX_users_file_id`(`id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 75 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
