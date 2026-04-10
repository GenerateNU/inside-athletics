-- Seed test survey data for the college ratings component.
-- This only inserts rows when the target college exists.

INSERT INTO "public"."users" (
  "id",
  "created_at",
  "updated_at",
  "deleted_at",
  "first_name",
  "last_name",
  "email",
  "username",
  "bio",
  "account_type",
  "expected_grad_year",
  "verified_athlete_status",
  "college_id",
  "sport_id",
  "division"
)
VALUES
  (
    '6a23f697-17f1-46dc-a04b-9556f6ad6d01',
    NOW(),
    NOW(),
    NULL,
    'Survey',
    'Tester One',
    'college-rating-seed-1@insideathletics.test',
    'college-rating-seed-1',
    NULL,
    FALSE,
    NULL,
    'pending',
    NULL,
    NULL,
    NULL
  ),
  (
    '6a23f697-17f1-46dc-a04b-9556f6ad6d02',
    NOW(),
    NOW(),
    NULL,
    'Survey',
    'Tester Two',
    'college-rating-seed-2@insideathletics.test',
    'college-rating-seed-2',
    NULL,
    FALSE,
    NULL,
    'pending',
    NULL,
    NULL,
    NULL
  ),
  (
    '6a23f697-17f1-46dc-a04b-9556f6ad6d03',
    NOW(),
    NOW(),
    NULL,
    'Survey',
    'Tester Three',
    'college-rating-seed-3@insideathletics.test',
    'college-rating-seed-3',
    NULL,
    FALSE,
    NULL,
    'pending',
    NULL,
    NULL,
    NULL
  ),
  (
    '6a23f697-17f1-46dc-a04b-9556f6ad6d04',
    NOW(),
    NOW(),
    NULL,
    'Survey',
    'Tester Four',
    'college-rating-seed-4@insideathletics.test',
    'college-rating-seed-4',
    NULL,
    FALSE,
    NULL,
    'pending',
    NULL,
    NULL,
    NULL
  )
ON CONFLICT ("id") DO NOTHING;

WITH target_college AS (
  SELECT "id"
  FROM "public"."colleges"
  WHERE "id" = '014d2c09-4023-445d-9779-66aff4824245'::uuid
),
selected_sports AS (
  SELECT
    "id" AS sport_id,
    ROW_NUMBER() OVER (
      ORDER BY
        CASE LOWER("name")
          WHEN 'soccer' THEN 1
          WHEN 'football' THEN 2
          WHEN 'basketball' THEN 3
          WHEN 'volleyball' THEN 4
          WHEN 'lacrosse' THEN 5
          WHEN 'baseball' THEN 6
          ELSE 100
        END,
        "name"
    ) AS sport_slot
  FROM "public"."sports"
  WHERE "deleted_at" IS NULL
  LIMIT 3
),
response_templates AS (
  SELECT *
  FROM (
    VALUES
      (
        1,
        '6a23f697-17f1-46dc-a04b-9556f6ad6d01'::uuid,
        5::smallint,
        4::smallint,
        5::smallint,
        4::smallint,
        5::smallint,
        5::smallint,
        4::smallint
      ),
      (
        2,
        '6a23f697-17f1-46dc-a04b-9556f6ad6d02'::uuid,
        4::smallint,
        4::smallint,
        4::smallint,
        3::smallint,
        4::smallint,
        4::smallint,
        4::smallint
      ),
      (
        3,
        '6a23f697-17f1-46dc-a04b-9556f6ad6d03'::uuid,
        4::smallint,
        3::smallint,
        4::smallint,
        4::smallint,
        4::smallint,
        5::smallint,
        3::smallint
      ),
      (
        4,
        '6a23f697-17f1-46dc-a04b-9556f6ad6d04'::uuid,
        5::smallint,
        5::smallint,
        4::smallint,
        4::smallint,
        5::smallint,
        4::smallint,
        5::smallint
      )
  ) AS template (
    response_slot,
    user_id,
    player_dev,
    academics_athletics_priority,
    academic_career_resources,
    mental_health_priority,
    environment,
    culture,
    transparency
  )
),
survey_ids AS (
  SELECT *
  FROM (
    VALUES
      (1, 1, '55fe6fe4-ef81-4631-8e22-8a7f0fa8f101'::uuid),
      (1, 2, '55fe6fe4-ef81-4631-8e22-8a7f0fa8f102'::uuid),
      (1, 3, '55fe6fe4-ef81-4631-8e22-8a7f0fa8f103'::uuid),
      (1, 4, '55fe6fe4-ef81-4631-8e22-8a7f0fa8f104'::uuid),
      (2, 1, '55fe6fe4-ef81-4631-8e22-8a7f0fa8f201'::uuid),
      (2, 2, '55fe6fe4-ef81-4631-8e22-8a7f0fa8f202'::uuid),
      (2, 3, '55fe6fe4-ef81-4631-8e22-8a7f0fa8f203'::uuid),
      (2, 4, '55fe6fe4-ef81-4631-8e22-8a7f0fa8f204'::uuid),
      (3, 1, '55fe6fe4-ef81-4631-8e22-8a7f0fa8f301'::uuid),
      (3, 2, '55fe6fe4-ef81-4631-8e22-8a7f0fa8f302'::uuid),
      (3, 3, '55fe6fe4-ef81-4631-8e22-8a7f0fa8f303'::uuid),
      (3, 4, '55fe6fe4-ef81-4631-8e22-8a7f0fa8f304'::uuid)
  ) AS ids (sport_slot, response_slot, survey_id)
)
INSERT INTO "public"."surveys" (
  "id",
  "created_at",
  "updated_at",
  "deleted_at",
  "user_id",
  "college_id",
  "sport_id",
  "player_dev",
  "academics_athletics_priority",
  "academic_career_resources",
  "mental_health_priority",
  "environment",
  "culture",
  "transparency"
)
SELECT
  survey_ids.survey_id,
  NOW(),
  NOW(),
  NULL,
  response_templates.user_id,
  target_college.id,
  selected_sports.sport_id,
  response_templates.player_dev,
  response_templates.academics_athletics_priority,
  response_templates.academic_career_resources,
  response_templates.mental_health_priority,
  response_templates.environment,
  response_templates.culture,
  response_templates.transparency
FROM target_college
JOIN selected_sports ON TRUE
JOIN response_templates ON TRUE
JOIN survey_ids
  ON survey_ids.sport_slot = selected_sports.sport_slot
 AND survey_ids.response_slot = response_templates.response_slot
ON CONFLICT ("id") DO NOTHING;
