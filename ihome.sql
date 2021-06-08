/*
 Navicat Premium Data Transfer

 Source Server         : 127.0.0.1
 Source Server Type    : MySQL
 Source Server Version : 50733
 Source Host           : localhost:3306
 Source Schema         : ihome

 Target Server Type    : MySQL
 Target Server Version : 50733
 File Encoding         : 65001

 Date: 08/06/2021 15:24:27
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for area
-- ----------------------------
DROP TABLE IF EXISTS `area`;
CREATE TABLE `area`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 31 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of area
-- ----------------------------
INSERT INTO `area` VALUES (1, '东城区');
INSERT INTO `area` VALUES (2, '西城区');
INSERT INTO `area` VALUES (3, '朝阳区');
INSERT INTO `area` VALUES (4, '海淀区');
INSERT INTO `area` VALUES (5, '昌平区');
INSERT INTO `area` VALUES (6, '丰台区');
INSERT INTO `area` VALUES (7, '房山区');
INSERT INTO `area` VALUES (8, '通州区');
INSERT INTO `area` VALUES (9, '顺义区');
INSERT INTO `area` VALUES (10, '大兴区');
INSERT INTO `area` VALUES (11, '怀柔区');
INSERT INTO `area` VALUES (12, '平谷区');
INSERT INTO `area` VALUES (13, '密云区');
INSERT INTO `area` VALUES (14, '延庆区');
INSERT INTO `area` VALUES (15, '石景山区');
INSERT INTO `area` VALUES (16, '东城区');
INSERT INTO `area` VALUES (17, '西城区');
INSERT INTO `area` VALUES (18, '朝阳区');
INSERT INTO `area` VALUES (19, '海淀区');
INSERT INTO `area` VALUES (20, '昌平区');
INSERT INTO `area` VALUES (21, '丰台区');
INSERT INTO `area` VALUES (22, '房山区');
INSERT INTO `area` VALUES (23, '通州区');
INSERT INTO `area` VALUES (24, '顺义区');
INSERT INTO `area` VALUES (25, '大兴区');
INSERT INTO `area` VALUES (26, '怀柔区');
INSERT INTO `area` VALUES (27, '平谷区');
INSERT INTO `area` VALUES (28, '密云区');
INSERT INTO `area` VALUES (29, '延庆区');
INSERT INTO `area` VALUES (30, '石景山区');

-- ----------------------------
-- Table structure for facility
-- ----------------------------
DROP TABLE IF EXISTS `facility`;
CREATE TABLE `facility`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 51 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of facility
-- ----------------------------
INSERT INTO `facility` VALUES (1, '无线网络');
INSERT INTO `facility` VALUES (2, '热水淋浴');
INSERT INTO `facility` VALUES (3, '空调');
INSERT INTO `facility` VALUES (4, '暖气');
INSERT INTO `facility` VALUES (5, '允许吸烟');
INSERT INTO `facility` VALUES (6, '饮水设备');
INSERT INTO `facility` VALUES (7, '牙具');
INSERT INTO `facility` VALUES (8, '香皂');
INSERT INTO `facility` VALUES (9, '拖鞋');
INSERT INTO `facility` VALUES (10, '手纸');
INSERT INTO `facility` VALUES (11, '毛巾');
INSERT INTO `facility` VALUES (12, '沐浴露、洗发露');
INSERT INTO `facility` VALUES (13, '冰箱');
INSERT INTO `facility` VALUES (14, '洗衣机');
INSERT INTO `facility` VALUES (15, '电梯');
INSERT INTO `facility` VALUES (16, '允许做饭');
INSERT INTO `facility` VALUES (17, '允许带宠物');
INSERT INTO `facility` VALUES (18, '允许聚会');
INSERT INTO `facility` VALUES (19, '门禁系统');
INSERT INTO `facility` VALUES (20, '停车位');
INSERT INTO `facility` VALUES (21, '有线网络');
INSERT INTO `facility` VALUES (22, '电视');
INSERT INTO `facility` VALUES (23, '浴缸');
INSERT INTO `facility` VALUES (24, '吃鸡');
INSERT INTO `facility` VALUES (25, '打台球');
INSERT INTO `facility` VALUES (26, '无线网络');
INSERT INTO `facility` VALUES (27, '热水淋浴');
INSERT INTO `facility` VALUES (28, '空调');
INSERT INTO `facility` VALUES (29, '暖气');
INSERT INTO `facility` VALUES (30, '允许吸烟');
INSERT INTO `facility` VALUES (31, '饮水设备');
INSERT INTO `facility` VALUES (32, '牙具');
INSERT INTO `facility` VALUES (33, '香皂');
INSERT INTO `facility` VALUES (34, '拖鞋');
INSERT INTO `facility` VALUES (35, '手纸');
INSERT INTO `facility` VALUES (36, '毛巾');
INSERT INTO `facility` VALUES (37, '沐浴露、洗发露');
INSERT INTO `facility` VALUES (38, '冰箱');
INSERT INTO `facility` VALUES (39, '洗衣机');
INSERT INTO `facility` VALUES (40, '电梯');
INSERT INTO `facility` VALUES (41, '允许做饭');
INSERT INTO `facility` VALUES (42, '允许带宠物');
INSERT INTO `facility` VALUES (43, '允许聚会');
INSERT INTO `facility` VALUES (44, '门禁系统');
INSERT INTO `facility` VALUES (45, '停车位');
INSERT INTO `facility` VALUES (46, '有线网络');
INSERT INTO `facility` VALUES (47, '电视');
INSERT INTO `facility` VALUES (48, '浴缸');
INSERT INTO `facility` VALUES (49, '吃鸡');
INSERT INTO `facility` VALUES (50, '打台球');

-- ----------------------------
-- Table structure for facility_houses
-- ----------------------------
DROP TABLE IF EXISTS `facility_houses`;
CREATE TABLE `facility_houses`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `facility_id` int(11) NOT NULL,
  `house_id` int(11) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 23 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of facility_houses
-- ----------------------------
INSERT INTO `facility_houses` VALUES (1, 1, 1);
INSERT INTO `facility_houses` VALUES (2, 2, 1);
INSERT INTO `facility_houses` VALUES (3, 3, 1);
INSERT INTO `facility_houses` VALUES (4, 7, 1);
INSERT INTO `facility_houses` VALUES (5, 12, 1);
INSERT INTO `facility_houses` VALUES (6, 14, 1);
INSERT INTO `facility_houses` VALUES (7, 16, 1);
INSERT INTO `facility_houses` VALUES (8, 17, 1);
INSERT INTO `facility_houses` VALUES (9, 18, 1);
INSERT INTO `facility_houses` VALUES (10, 21, 1);
INSERT INTO `facility_houses` VALUES (11, 22, 1);
INSERT INTO `facility_houses` VALUES (12, 1, 2);
INSERT INTO `facility_houses` VALUES (13, 2, 2);
INSERT INTO `facility_houses` VALUES (14, 3, 2);
INSERT INTO `facility_houses` VALUES (15, 7, 2);
INSERT INTO `facility_houses` VALUES (16, 12, 2);
INSERT INTO `facility_houses` VALUES (17, 14, 2);
INSERT INTO `facility_houses` VALUES (18, 16, 2);
INSERT INTO `facility_houses` VALUES (19, 17, 2);
INSERT INTO `facility_houses` VALUES (20, 18, 2);
INSERT INTO `facility_houses` VALUES (21, 21, 2);
INSERT INTO `facility_houses` VALUES (22, 22, 2);

-- ----------------------------
-- Table structure for house
-- ----------------------------
DROP TABLE IF EXISTS `house`;
CREATE TABLE `house`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `area_id` int(11) NOT NULL,
  `title` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `price` int(11) NOT NULL DEFAULT 0,
  `address` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `room_count` int(11) NOT NULL DEFAULT 1,
  `acreage` int(11) NOT NULL DEFAULT 0,
  `unit` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `capacity` int(11) NOT NULL DEFAULT 1,
  `beds` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `deposit` int(11) NOT NULL DEFAULT 0,
  `min_days` int(11) NOT NULL DEFAULT 1,
  `max_days` int(11) NOT NULL DEFAULT 0,
  `order_count` int(11) NOT NULL DEFAULT 0,
  `index_image_url` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `ctime` datetime NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of house
-- ----------------------------
INSERT INTO `house` VALUES (1, 5, 1, '上奥世纪中心', 66600, '西三旗桥东建材城1号', 2, 60, '2室1厅', 3, '双人床2张', 20000, 3, 0, 1, '/group1/default/20210421/14/07/5/d5025964.jpg', '2021-04-23 06:24:11');
INSERT INTO `house` VALUES (2, 5, 1, '上奥世纪中心', 66600, '西三旗桥东建材城1号', 2, 60, '2室1厅', 3, '双人床2张', 20000, 3, 0, 0, '', '2021-04-23 06:28:28');

-- ----------------------------
-- Table structure for house_image
-- ----------------------------
DROP TABLE IF EXISTS `house_image`;
CREATE TABLE `house_image`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `url` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `house_id` int(11) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of house_image
-- ----------------------------
INSERT INTO `house_image` VALUES (1, '/group1/default/20210421/14/07/5/d5025964.jpg', 1);
INSERT INTO `house_image` VALUES (2, '/group1/default/20210421/14/07/5/d5025964.jpg', 1);

-- ----------------------------
-- Table structure for order_house
-- ----------------------------
DROP TABLE IF EXISTS `order_house`;
CREATE TABLE `order_house`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `house_id` int(11) NOT NULL,
  `begin_date` datetime NOT NULL,
  `end_date` datetime NOT NULL,
  `days` int(11) NOT NULL DEFAULT 0,
  `house_price` int(11) NOT NULL DEFAULT 0,
  `amount` int(11) NOT NULL DEFAULT 0,
  `status` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'WAIT_ACCEPT',
  `comment` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `ctime` datetime NOT NULL,
  `credit` tinyint(1) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of order_house
-- ----------------------------
INSERT INTO `order_house` VALUES (1, 5, 1, '2017-11-11 13:23:49', '2017-11-12 13:23:49', 2, 66600, 133200, 'COMPLETE', '烂房子！', '2021-04-25 09:46:54', 0);
INSERT INTO `order_house` VALUES (2, 5, 1, '2017-11-11 21:23:49', '2017-11-12 21:23:49', 2, 66600, 133200, 'WAIT_ACCEPT', '', '2021-04-25 06:56:04', 0);

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `password_hash` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `mobile` varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `real_name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `id_card` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `avatar_url` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `mobile`(`mobile`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 18 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES (17, '18148789239', 'e10adc3949ba59abbe56e057f20f883e', '18148789239', '', '', 'http://192.168.0.65:3666/group1/default/20210510/15/42/4/8430311709_f16e717e14_k-e1553251390297.jpg');

SET FOREIGN_KEY_CHECKS = 1;
