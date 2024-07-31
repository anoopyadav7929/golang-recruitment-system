CREATE TABLE `job_applications` (
   `job_id` bigint NOT NULL,
   `user_id` bigint NOT NULL,
   PRIMARY KEY (`job_id`,`user_id`)
 ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci


CREATE TABLE `jobs` (
   `id` int NOT NULL AUTO_INCREMENT,
   `title` varchar(100) NOT NULL,
   `description` text,
   `posted_on` timestamp NOT NULL,
   `total_applications` bigint DEFAULT '0',
   `company_name` varchar(255) NOT NULL,
   `posted_by` int DEFAULT NULL,
   PRIMARY KEY (`id`)
 ) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci


CREATE TABLE `profiles` (
   `id` int NOT NULL AUTO_INCREMENT,
   `user_id` int DEFAULT NULL,
   `resume_id` int DEFAULT NULL,
   `skills` text,
   `education` text,
   `experience` text,
   `name` varchar(100) DEFAULT NULL,
   `email` varchar(100) DEFAULT NULL,
   `phone` varchar(20) DEFAULT NULL,
   PRIMARY KEY (`id`)
 ) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci


CREATE TABLE `resumes` (
   `id` int NOT NULL AUTO_INCREMENT,
   `user_id` int NOT NULL,
   `doc_content` longblob NOT NULL,
   `doc_type` enum('pdf','docx') NOT NULL,
   PRIMARY KEY (`id`),
   KEY `user_id` (`user_id`),
   CONSTRAINT `resumes_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
 ) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci


CREATE TABLE `users` (
   `id` int NOT NULL AUTO_INCREMENT,
   `name` varchar(100) NOT NULL,
   `email` varchar(100) NOT NULL,
   `address` varchar(255) DEFAULT NULL,
   `user_type` varchar(10) NOT NULL,
   `password_hash` varchar(255) NOT NULL,
   `profile_headline` varchar(255) DEFAULT NULL,
   PRIMARY KEY (`id`),
   UNIQUE KEY `email` (`email`),
   CONSTRAINT `users_chk_1` CHECK ((`user_type` in (_utf8mb4'applicant',_utf8mb4'admin')))
 ) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
