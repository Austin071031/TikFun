CREATE DATABASE IF NOT EXISTS `tiktok` /*!40100 DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci */;
USE `tiktok`;

-- users
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`
(
    `id`                bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `name`              varchar(128)        NOT NULL DEFAULT '' COMMENT '用户昵称',
    `signature`         varchar(128)        NOT NULL DEFAULT '' COMMENT '个人简介',
    `follow_count`       int(10)             NOT NULL DEFAULT 1 COMMENT '关注数量',
    `follower_count`     int(10)             NOT NULL DEFAULT 0 COMMENT '粉丝数量',
    `isfollow`          boolean             NOT NULL DEFAULT FALSE COMMENT '是否被关注',
    `avatar`            varchar(128)        NOT NULL DEFAULT '' COMMENT '个人头像图片地址',
    `background_image`   varchar(128)        NOT NULL DEFAULT '' COMMENT '背景图片地址',
    `favorite_count`     int(10)             NOT NULL DEFAULT 0 COMMENT '点赞视频数',
    `total_favorited`    varchar(128)        NOT NULL DEFAULT '0' COMMENT '视频获赞数',
    `work_count`         int(10)             NOT NULL DEFAULT 0 COMMENT '作品数',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='用户表';

-- loginuserdata
DROP TABLE IF EXISTS `loginuserdata`;
CREATE TABLE `loginuserdata`
(
    `name`          varchar(128)        NOT NULL COMMENT '主键用户昵称',
    `password`      varchar(128)        NOT NULL DEFAULT '' COMMENT '哈希加密密码，验证用户',
    `token`         varchar(512)        NOT NULL DEFAULT '' COMMENT 'Token',
    PRIMARY KEY (`name`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='Token表';

-- like
DROP TABLE IF EXISTS `likes`;
CREATE TABLE `likes`
(
    `id`                bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `name`              varchar(128)        NOT NULL DEFAULT '' COMMENT '主键用户昵称',
    `videoId`           int(128)            NOT NULL COMMENT '视频编号',
    `liked`             int(10)             NOT NULL DEFAULT 0 COMMENT '是否点赞',
    `create_time`      	datetime            NOT NULL DEFAULT '1970-01-01' COMMENT '创建时间',
  
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='点赞表';

-- videos
DROP TABLE IF EXISTS `videos`;
CREATE TABLE `videos`
(
    `id`                bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `name`              varchar(128)        NOT NULL DEFAULT '' COMMENT '作者',
    `playurl`           varchar(128)        NOT NULL DEFAULT '' COMMENT '视频路径',
    `coverurl`          varchar(128)        NOT NULL DEFAULT 1 COMMENT '封面路径',
    `favoritecount`     int(10)             NOT NULL DEFAULT 0 COMMENT '点赞数量',
    `commentcount`      int(10)             NOT NULL DEFAULT 0 COMMENT '评论数量',
    `isfavorite`        boolean             NOT NULL DEFAULT FALSE COMMENT '视频是否被点赞',
    `title`             varchar(128)        NOT NULL DEFAULT '' COMMENT '视频标题',
    `uploadtime`        datetime            NOT NULL DEFAULT '1970-01-01' COMMENT '视频上传时间',
  
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='视频表';

-- Comment
DROP TABLE IF EXISTS `comments`;
CREATE TABLE `comments`
(
    `id`                bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `video_id`          bigint(20) unsigned	NOT NULL COMMENT '视频id',
    `user_id`           bigint(20) unsigned	NOT NULL COMMENT '用户id',
    `content`           varchar(128)        NOT NULL DEFAULT '' COMMENT '内容',
    `create_date`      	varchar(128)        NOT NULL DEFAULT '' COMMENT '创建时间',
    PRIMARY KEY (`id`)
		-- FOREIGN KEY (`user_id`)   REFERENCES user(`id`)
    -- FOREIGN KEY (`video_id`)  REFERENCES videos(`id`) 
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='评论表';
