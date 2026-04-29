-- Add Stripe customer ID to users table
ALTER TABLE "public"."users"
ADD COLUMN IF NOT EXISTS "stripe_customer_id" varchar(255) CONSTRAINT "uni_users_stripe_customer_id" UNIQUE;

-- Create user_subscriptions table
CREATE TABLE IF NOT EXISTS "public"."user_subscriptions" (
    "id"                     uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    "created_at"             timestamptz NOT NULL DEFAULT now(),
    "updated_at"             timestamptz NOT NULL DEFAULT now(),
    "deleted_at"             timestamptz,
    "user_id"                uuid        NOT NULL UNIQUE REFERENCES "public"."users"("id") ON DELETE CASCADE,
    "stripe_subscription_id" varchar(255) NOT NULL,
    "stripe_price_id"        varchar(255) NOT NULL,
    "status"                 varchar(50)  NOT NULL,
    "current_period_start"   timestamptz NOT NULL,
    "current_period_end"     timestamptz NOT NULL,
    "canceled_at"            timestamptz
);

CREATE INDEX IF NOT EXISTS "idx_user_subscriptions_deleted_at" ON "public"."user_subscriptions"("deleted_at");
