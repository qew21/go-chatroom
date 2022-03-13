-- ----------------------------
-- Table structure for message
-- ----------------------------
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message`
(
    `id`            bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `user_id`       bigint(20) unsigned NOT NULL COMMENT '所属类型的id',
    `request_id`    bigint(20) NOT NULL COMMENT '请求id',
    `sender_type`   tinyint(3) NOT NULL COMMENT '发送者类型',
    `sender_id`     bigint(20) unsigned NOT NULL COMMENT '发送者id',
    `receiver_type` tinyint(3) NOT NULL COMMENT '接收者类型,1:个人；2：群组',
    `receiver_id`   bigint(20) unsigned NOT NULL COMMENT '接收者id,如果是单聊信息，则为user_id，如果是群组消息，则为group_id',
    `to_user_ids`   varchar(255) NOT NULL COMMENT '需要@的用户id列表，多个用户用，隔开',
    `type`          tinyint(4) NOT NULL COMMENT '消息类型',
    `content`       blob         NOT NULL COMMENT '消息内容',
    `seq`           bigint(20) unsigned NOT NULL unique COMMENT '消息序列号',
    `send_time`     datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP (3) COMMENT '消息发送时间',
    `status`        tinyint(255) NOT NULL DEFAULT '0' COMMENT '消息状态，0：未处理1：消息撤回',
    `create_time`   datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`   datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT ='消息';

  -- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`
(
    `id`           bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `phone_number` varchar(20)   NOT NULL COMMENT '手机号',
    `nickname`     varchar(20)   NOT NULL COMMENT '昵称',
    `sex`          tinyint(4) NOT NULL COMMENT '性别，0:未知；1:男；2:女',
    `password`     varchar(32)  NOT NULL COMMENT '哈希密码',
    `create_time`  datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`  datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_phone_number` (`phone_number`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT ='用户';

-- ----------------------------
-- Table structure for group
-- ----------------------------
DROP TABLE IF EXISTS `group`;
CREATE TABLE `group`
(
    `id`           bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `name`         varchar(50)   NOT NULL COMMENT '群组名称',
    `introduction` varchar(255)  NOT NULL COMMENT '群组简介',
    `create_time`  datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`  datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT ='群组';

  -- ----------------------------
-- Table structure for group_user
-- ----------------------------
DROP TABLE IF EXISTS `group_user`;
CREATE TABLE `group_user`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `group_id`    bigint(20) unsigned NOT NULL COMMENT '组id',
    `user_id`     bigint(20) unsigned NOT NULL COMMENT '用户id',
    `member_type` tinyint(4) NOT NULL COMMENT '成员类型，1：管理员；2：普通成员',
    `remarks`     varchar(20)   NOT NULL COMMENT '备注',
    `create_time` datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_group_id_user_id` (`group_id`, `user_id`) USING BTREE,
    KEY           `idx_user_id` (`user_id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT ='群组成员';