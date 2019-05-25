package db

const createBlogTable = `
CREATE TABLE IF NOT EXISTS blog(
	id CHAR(40) NOT NULL PRIMARY KEY,  
	title TEXT,
	content TEXT,
	userid CHAR(40) NOT NULL, 
	username TEXT NOT NULL,
	tags CHAR(40) NOT NULL,
	create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS comment(
	id CHAR(40) NOT NULL PRIMARY KEY,
	content TEXT NOT NULL,
	userid CHAR(40) NOT NULL,
	username TEXT NOT NULL,
	blogid CHAR(40) NOT NULL,
	create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS tags(
	blog_id CHAR(40) NOT NULL,
	tag_id CHAR(40) NOT NULL
);
CREATE TABLE IF NOT EXISTS tag(
	id CHAR(40) NOT NULL PRIMARY KEY,
	name CHAR(40) NOT NULL UNIQUE
);
`

const findBlogsByTitle = `
SELECT 
blog.id AS id,
blog.title AS title,
blog.content AS content,
blog.userid AS userid,
blog.username AS username,
blog.create_at AS create_at,
blog.update_at AS update_at,
comment.id AS commentid,
comment.content AS commentcontent,
comment.userid AS commentuserid,
comment.username AS commentusername,
comment.blogid AS commentblogid,
comment.create_at AS commentcreate_at,
comment.update_at As commentupdate_at
FROM comment LEFT JOIN blog
ON blog.id = comment.blogid
WHERE blog.username = $1
ORDER BY create_at DESC
`

const findBlogByID = `
SELECT 
blog.id AS id,
blog.title AS title,
blog.content AS content,
blog.userid AS userid,
blog.username AS username,
blog.create_at AS create_at,
blog.update_at AS update_at,
comment.id AS commentid,
comment.content AS commentcontent,
comment.userid AS commentuserid,
comment.username AS commentusername,
comment.blogid AS commentblogid,
comment.create_at AS commentcreate_at,
comment.update_at As commentupdate_at,
t.tagid AS tagid,
t.tagname AS tagname
FROM comment LEFT JOIN blog
ON blog.id = comment.blogid
LEFT JOIN (
	SELECT
	tags.blog_id as blogid,
	tag.name as tagname,
	tag.id AS tagid
	FROM tags INNER JOIN tag 
	ON tags.tag_id = tag.id
	WHERE tags.blog_id = $1
) t ON blog.id = t.blogid
WHERE blog.id = $1
`

const findBlogsByUserID = `
SELECT 
blog.id AS id,
blog.title AS title,
blog.content AS content,
blog.userid AS userid,
blog.username AS username,
blog.create_at AS create_at,
blog.update_at AS update_at,
comment.id AS commentid,
comment.content AS commentcontent,
comment.userid AS commentuserid,
comment.username AS commentusername,
comment.blogid AS commentblogid,
comment.create_at AS commentcreate_at,
comment.update_at As commentupdate_at
FROM comment LEFT JOIN blog
ON blog.id = comment.blogid
WHERE blog.userid = $1
ORDER BY create_at DESC
`

const findBlogsByUserName = `
SELECT 
blog.id AS id,
blog.title AS title,
blog.content AS content,
blog.userid AS userid,
blog.username AS username,
blog.create_at AS create_at,
blog.update_at AS update_at,
comment.id AS commentid,
comment.content AS commentcontent,
comment.userid AS commentuserid,
comment.username AS commentusername, 
comment.blogid AS commentblogid,
comment.create_at AS commentcreate_at,
comment.update_at As commentupdate_at 
FROM comment LEFT JOIN blog
ON blog.id = comment.blogid
WHERE blog.username = $1
ORDER BY create_at DESC
`

const searchBlog = `
SELECT 
id,
title,
content,
userid,
username,
create_at,
update_at
WHERE blog.title LIKE  $1
OR blog.content LIKE $2
ORDER BY create_at DESC
`

const findBlogsByTagID = `
SELECT 
blog.id AS id,
blog.title AS title,
blog.content AS content,
blog.userid AS userid,
blog.username AS username,
blog.create_at AS create_at,
blog.update_at AS update_at,
comment.id AS commentid,
comment.content AS commentcontent,
comment.userid AS commentuserid,
comment.username AS commentusername, 
comment.blogid AS commentblogid,
comment.create_at AS commentcreate_at,
comment.update_at As commentupdate_at 
FROM comment LEFT JOIN blog
ON blog.id = comment.blogid
WHERE blog.id IN (
	SELECT 
	blog_id 
	FROM tags 
	WHERE tag_id =$1
)
ORDER BY create_at DESC
`
