SELECT posts.*,
            (SELECT COUNT(*) FROM post_likes WHERE post_likes.post_id = posts.id) AS like_count,
            (SELECT COUNT(*) FROM comments WHERE comments.post_id = posts.id) AS comment_count,
            (SELECT COUNT(*) > 0 FROM post_likes WHERE post_likes.post_id = posts.id AND post_likes.user_id = '0214ce3c-773f-4e40-8b7b-8ebba293a6c4') AS is_liked FROM "posts" JOIN tag_posts ON posts.id = tag_posts.id WHERE posts.college_id = '136f7843-384a-4c04-95e8-890dfed57715' AND "posts"."deleted_at" IS NULL LIMIT 20