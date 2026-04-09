-- Create "surveys" table
CREATE TABLE "public"."surveys" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" uuid NOT NULL,
  "college_id" uuid NOT NULL,
  "sport_id" uuid NOT NULL,
  "player_dev" smallint NOT NULL,
  "academics_athletics_priority" smallint NOT NULL,
  "academic_career_resources" smallint NOT NULL,
  "mental_health_priority" smallint NOT NULL,
  "environment" smallint NOT NULL,
  "culture" smallint NOT NULL,
  "transparency" smallint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_surveys_college" FOREIGN KEY ("college_id") REFERENCES "public"."colleges" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_surveys_sport" FOREIGN KEY ("sport_id") REFERENCES "public"."sports" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_surveys_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_surveys_college_id" to table: "surveys"
CREATE INDEX "idx_surveys_college_id" ON "public"."surveys" ("college_id");
-- Create index "idx_surveys_deleted_at" to table: "surveys"
CREATE INDEX "idx_surveys_deleted_at" ON "public"."surveys" ("deleted_at");
-- Create index "idx_surveys_sport_id" to table: "surveys"
CREATE INDEX "idx_surveys_sport_id" ON "public"."surveys" ("sport_id");
-- Create index "idx_surveys_user_id" to table: "surveys"
CREATE INDEX "idx_surveys_user_id" ON "public"."surveys" ("user_id");
