-- Create index "idx_user_college" to table: "college_follows"
CREATE UNIQUE INDEX "idx_user_college" ON "public"."college_follows" ("college_id", "user_id");
-- Create index "idx_user_sport" to table: "sport_follows"
CREATE UNIQUE INDEX "idx_user_sport" ON "public"."sport_follows" ("sport_id", "user_id");
-- Create index "idx_user_tag" to table: "tag_follows"
CREATE UNIQUE INDEX "idx_user_tag" ON "public"."tag_follows" ("tag_id", "user_id");
