-- MySQL dump 10.13  Distrib 5.7.16, for Linux (x86_64)
--
-- Host: localhost    Database: maintain
-- ------------------------------------------------------
-- Server version	5.7.16-0ubuntu0.16.04.1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `cluster_minion`
--

DROP TABLE IF EXISTS `cluster_minion`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cluster_minion` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `cluster_id` int(11) NOT NULL DEFAULT '0' COMMENT '服务器组id',
  `minion_id` int(11) NOT NULL DEFAULT '0' COMMENT 'minion id',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cluster_minion`
--

LOCK TABLES `cluster_minion` WRITE;
/*!40000 ALTER TABLE `cluster_minion` DISABLE KEYS */;
INSERT INTO `cluster_minion` VALUES (5,1,2,'2017-06-23 09:25:33'),(7,5,2,'2017-06-30 03:36:20'),(8,6,2,'2017-06-30 08:14:18');
/*!40000 ALTER TABLE `cluster_minion` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `clusters`
--

DROP TABLE IF EXISTS `clusters`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `clusters` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '组名称',
  `online` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否在线，0-在线; 1-下线',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `clusters`
--

LOCK TABLES `clusters` WRITE;
/*!40000 ALTER TABLE `clusters` DISABLE KEYS */;
INSERT INTO `clusters` VALUES (1,'web-immortal',0,'2017-06-20 06:20:53','2017-06-23 09:25:33',NULL),(2,'web-immortal',0,'2017-06-23 07:25:25','2017-06-23 08:34:45','2017-06-23 08:34:45'),(3,'web-immortal',0,'2017-06-23 07:29:51','2017-06-23 08:34:48','2017-06-23 08:34:48'),(4,'web-immortal',0,'2017-06-23 07:30:31','2017-06-23 08:34:58','2017-06-23 08:34:58'),(5,'web-api',0,'2017-06-23 07:31:42','2017-06-30 03:36:20',NULL),(6,'web-v1.5',0,'2017-06-30 08:14:18','2017-06-30 08:14:18',NULL);
/*!40000 ALTER TABLE `clusters` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `minions`
--

DROP TABLE IF EXISTS `minions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `minions` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '机器名称',
  `ipv4_internal` varchar(15) COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '内网ip地址',
  `ipv4_external` varchar(15) COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '外网ip地址',
  `online` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否在线，0-在线; 1-下线',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `minions`
--

LOCK TABLES `minions` WRITE;
/*!40000 ALTER TABLE `minions` DISABLE KEYS */;
INSERT INTO `minions` VALUES (1,'vm-centos7.1-2','192.168.56.67','124.205.90.68',0,'2017-06-23 03:32:19','2017-06-23 04:12:22','2017-06-23 04:12:22'),(2,'vm-centos7.1-2','192.168.56.67','124.205.90.68',0,'2017-06-23 04:12:45','2017-06-23 04:12:45',NULL),(3,'vm-centos7.1-1','192.168.56.66','124.205.90.68',0,'2017-06-23 06:38:07','2017-06-23 06:38:07',NULL),(4,'123','','',0,'2017-06-23 06:44:29','2017-06-23 06:46:11','2017-06-23 06:46:11'),(5,'123','','',0,'2017-06-23 06:46:22','2017-06-23 06:47:13','2017-06-23 06:47:13'),(6,'123','','',0,'2017-06-23 06:46:43','2017-06-23 06:47:16','2017-06-23 06:47:16');
/*!40000 ALTER TABLE `minions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `projects`
--

DROP TABLE IF EXISTS `projects`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `projects` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '项目名称',
  `root_dir` varchar(500) COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '项目根目录',
  `cmd_script` varchar(500) COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT 'shell文件，用于git部署',
  `cluster_id` int(11) NOT NULL DEFAULT '0' COMMENT '服务器组id',
  `current_tag` varchar(50) COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '当前git tag',
  `previous_tag` varchar(50) COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '上一个git tag',
  `online` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否在线，0-在线; 1-下线',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `projects`
--

LOCK TABLES `projects` WRITE;
/*!40000 ALTER TABLE `projects` DISABLE KEYS */;
INSERT INTO `projects` VALUES (3,'v2','/opt/adcloud-v2/','',1,'v0.1.24','v0.1.25',1,'2017-06-17 07:30:14','2017-06-20 02:20:28',NULL),(4,'immortal','/opt/immortal/','/srv/salt/scripts/immortal-upgrade.sh',1,'v0.1.3','v0.1.3',0,'2017-06-17 07:31:49','2017-07-05 08:31:48',NULL),(5,'123','234','',1,'','',0,'2017-06-20 09:42:59','2017-06-20 11:01:48','2017-06-20 11:01:48'),(6,'123','234','',1,'','',0,'2017-06-20 09:51:05','2017-06-20 11:01:53','2017-06-20 11:01:53'),(7,'dfsfsdfsdf','234','',1,'','',0,'2017-06-20 09:59:14','2017-06-22 09:19:16','2017-06-22 09:19:16'),(8,'adcloud-api','sdfsf','',1,'','',0,'2017-06-20 10:29:02','2017-06-22 09:19:20','2017-06-22 09:19:20'),(9,'test','ttee','',1,'','',0,'2017-06-20 11:01:43','2017-06-21 02:44:34','2017-06-21 02:44:34'),(10,'monitor-lua','/opt/monitor-lua/','/srv/salt/scripts/monitor-lua-upgrade.sh',1,'v1.16','v1.16',0,'2017-08-15 07:21:00','2017-08-15 07:22:20',NULL);
/*!40000 ALTER TABLE `projects` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `upgrades`
--

DROP TABLE IF EXISTS `upgrades`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `upgrades` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `project_id` int(11) NOT NULL DEFAULT '0' COMMENT '项目id',
  `tag_from` varchar(50) COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '升级前tag',
  `tag_to` varchar(50) COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '升级后tag',
  `minions_succ` text COLLATE utf8_unicode_ci COMMENT '升级成功的minions',
  `minions_fail` text COLLATE utf8_unicode_ci COMMENT '升级失败的minions',
  `duration` int(11) NOT NULL DEFAULT '0' COMMENT '升级耗时，单位毫秒',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态，0-成功; 1-失败',
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '用户id',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=49 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `upgrades`
--

LOCK TABLES `upgrades` WRITE;
/*!40000 ALTER TABLE `upgrades` DISABLE KEYS */;
INSERT INTO `upgrades` VALUES (1,4,'v0.1.2','v0.1.3','vm-centos7.1-2','',477,0,0,'2017-06-27 03:58:30'),(2,4,'v0.1.3','v0.1.2','vm-centos7.1-2','',481,0,0,'2017-06-27 04:17:23'),(3,4,'v0.1.2','v0.1.1','vm-centos7.1-2','',440,0,0,'2017-06-27 04:23:52'),(4,4,'v0.1.1','v0.1.2','vm-centos7.1-2','',5485,0,0,'2017-06-27 04:37:08'),(5,4,'v0.1.2','v0.1.3','vm-centos7.1-2','',1120,0,0,'2017-06-27 06:29:59'),(6,4,'v0.1.3','v0.1.0','vm-centos7.1-2','',568,0,0,'2017-06-27 06:44:14'),(7,4,'v0.1.0','v0.1.1','vm-centos7.1-2','',477,0,0,'2017-06-27 06:45:17'),(8,4,'v0.1.1','v0.1.3','vm-centos7.1-2','',536,0,0,'2017-06-27 08:25:15'),(9,4,'v0.1.3','v0.1.2','vm-centos7.1-2','',477,0,0,'2017-06-27 08:29:32'),(10,4,'v0.1.2','v0.1.3','vm-centos7.1-2','',472,0,0,'2017-06-27 08:30:10'),(11,4,'v0.1.3','v0.1.2','vm-centos7.1-2','',478,0,0,'2017-06-27 08:32:03'),(12,4,'v0.1.2','v0.1.3','vm-centos7.1-2','',1326,0,0,'2017-06-28 03:42:00'),(13,4,'v0.1.3','v0.1.2','vm-centos7.1-2','',649,0,0,'2017-06-28 03:42:23'),(14,4,'v0.1.2','v0.1.3','vm-centos7.1-2','',474,0,0,'2017-06-28 03:44:54'),(15,4,'v0.1.3','v0.1.2','vm-centos7.1-2','',458,0,0,'2017-06-28 03:46:19'),(16,4,'v0.1.2','v0.1.3','vm-centos7.1-2','',655,0,0,'2017-06-28 03:48:19'),(17,4,'v0.1.3','v0.1.2','vm-centos7.1-2','',659,0,0,'2017-06-28 03:52:34'),(18,4,'v0.1.2','v0.1.3','vm-centos7.1-2','',597,0,0,'2017-06-28 03:54:33'),(19,4,'v0.1.3','v0.1.3','vm-centos7.1-2','',6360,0,0,'2017-07-03 09:13:41'),(20,4,'v0.1.3','v0.1.3','vm-centos7.1-2','',416,0,0,'2017-07-03 09:17:02'),(21,4,'v0.1.3','v0.1.3','vm-centos7.1-2','',429,0,0,'2017-07-03 09:20:33'),(22,4,'v0.1.3','v0.1.3','vm-centos7.1-2','',5468,0,0,'2017-07-03 09:21:43'),(23,4,'v0.1.3','v0.1.3','vm-centos7.1-2','',1396,0,0,'2017-07-03 09:37:23'),(24,4,'v0.1.3','v0.1.2','vm-centos7.1-2','',1575,0,0,'2017-07-03 09:44:47'),(25,4,'v0.1.2','v0.1.0','vm-centos7.1-2','',488,0,0,'2017-07-03 09:47:31'),(26,4,'v0.1.0','v0.1.3','vm-centos7.1-2','',477,0,0,'2017-07-03 09:47:56'),(27,4,'v0.1.3','v0.1.3','vm-centos7.1-2','',410,0,0,'2017-07-03 09:51:29'),(28,4,'v0.1.3','v0.1.2','vm-centos7.1-2','',520,0,0,'2017-07-03 09:52:51'),(29,4,'v0.1.2','v0.1.3','','',86,1,0,'2017-07-05 08:13:23'),(30,4,'v0.1.2','v0.1.3','','',58,1,0,'2017-07-05 08:15:42'),(31,4,'v0.1.2','v0.1.3','','',62,1,0,'2017-07-05 08:16:24'),(32,4,'v0.1.2','v0.1.3','','',105,1,0,'2017-07-05 08:22:34'),(33,4,'v0.1.2','v0.1.3','vm-centos7.1-2','',1279,0,0,'2017-07-05 08:30:38'),(34,4,'v0.1.3','v0.1.3','vm-centos7.1-2','',457,0,0,'2017-07-05 08:31:48'),(35,4,'v0.1.3','v0.1.3','vm-centos7.1-2','',439,0,0,'2017-07-05 09:05:29'),(36,4,'v0.1.3','v0.1.3','vm-centos7.1-2','',1080,0,0,'2017-08-15 06:06:54'),(37,10,'v1.16','v1.16','vm-centos7.1-2','',1856,0,0,'2017-08-15 07:22:20'),(38,10,'v1.16','v1.16','vm-centos7.1-2','',1695,0,0,'2017-08-15 07:23:53'),(39,10,'v1.16','v1.16','vm-centos7.1-2','',1867,0,0,'2017-08-15 07:24:49'),(40,10,'v1.16','v1.16','vm-centos7.1-2','',1778,0,0,'2017-08-15 07:26:19'),(41,10,'v1.16','v1.16','vm-centos7.1-2','',1906,0,0,'2017-08-15 07:27:15'),(42,10,'v1.16','v1.16','vm-centos7.1-2','',6694,0,0,'2017-08-15 07:28:37'),(43,10,'v1.16','v1.16','vm-centos7.1-2','',1785,0,0,'2017-08-15 07:33:56'),(44,10,'v1.16','v1.16','vm-centos7.1-2','',3548,0,0,'2017-08-15 07:35:03'),(45,10,'v1.16','v1.16','vm-centos7.1-2','',3712,0,0,'2017-08-15 07:36:19'),(46,10,'v1.16','v1.16','vm-centos7.1-2','',1780,0,0,'2017-08-15 07:43:59'),(47,10,'v1.16','v1.16','vm-centos7.1-2','',2670,0,0,'2017-08-15 07:46:45'),(48,10,'v1.16','v1.16','vm-centos7.1-2','',2750,0,0,'2017-08-15 07:57:10');
/*!40000 ALTER TABLE `upgrades` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `users` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '用户名',
  `password` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '密码',
  `realname` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '姓名',
  `phone` varchar(25) COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '手机号码',
  `email` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '邮箱',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态0-正常;1-禁用',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'admin','77a40784aeb6546fe406493337f0664b','administrator','15010240697','tangshouqiang@tvmining.com',0,'2017-06-29 07:55:28','2017-06-29 07:55:28');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2017-10-10 11:11:17
