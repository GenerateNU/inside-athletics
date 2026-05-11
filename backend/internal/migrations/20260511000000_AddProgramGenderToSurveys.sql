-- Separate men's vs women's program ratings for the same (college, sport).
ALTER TABLE "public"."surveys"
  ADD COLUMN "program_gender" varchar(16) NOT NULL DEFAULT '';

CREATE INDEX "idx_surveys_program_gender" ON "public"."surveys" ("program_gender");
