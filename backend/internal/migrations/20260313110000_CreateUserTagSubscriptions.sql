-- Create "user_tag_subscriptions" table
CREATE TABLE "public"."user_tag_subscriptions" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL,
  "user_id" uuid NOT NULL,
  "tag_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_user_tag_subscriptions_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_user_tag_subscriptions_tag" FOREIGN KEY ("tag_id") REFERENCES "public"."tags" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_user_tag_subscriptions_user_tag" to table: "user_tag_subscriptions"
CREATE UNIQUE INDEX "idx_user_tag_subscriptions_user_tag" ON "public"."user_tag_subscriptions" ("user_id", "tag_id");
