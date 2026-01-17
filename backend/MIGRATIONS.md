
## Migration Guide

Migrations are a way we can manage our database. It serves as a log of how our database schema has changed over time.
Think of it as git logs but for our database! It makes it really easy to re-create our database, and rollback changes when
needed

### What tools our Migration workflow uses:
- [Atlas](https://atlasgo.io/): A really great tool we can use for migration generation and managment. It can also be used to store
migrations, but not completely necessary for a project of our scope. This is used to generate migrations automatically from our GORM models
- GORM: works with Atlas to make the migration generation process really easy. Atlas scans the repo for GORM models, identifies any changes in the schema, and creates a migration file for them

### Generating a Migration

All of our migrations live in the `./backend/internal/migrations` directory. These hold the information needed apply the changes
to the database (or to roll them back)

To generate a migration run: `make gen-migration`. This will create a new file in your migrations directory

### Testing a Migration

NEVER apply a migration to prod. This will be done by the TL's after your code has been reviewed. In order to test your
migration run it on your local supabase database. This can be run with `make migration-dev`.



## Reverting Migrations

You can revert a migration by running: atlas migrate down --env [env]. To preview what will happen by running this command
you can add the `--dry-run` flag. This will generate a plan so you can make sure you are reverting what you actually want.

In the makefile I have targets for reverting a single migration for both prod and dev. If you want to use 
these other flags run them yourself. 

After reverting if you want to completely get rid of this change you must:

1) Delete the Migration: `rm migrations/20240101120000_bad_change.sql`

2) Regenerate the hash file: `atlas migrate hash`

Here are some useful flags for reverting migrations: 

| Goal                  | Command          |
|-----------------------|------------------|
| Revert last migration | `--steps 1`      |
| Revert multiple       | `--steps N`      |
| Revert to version     | `--to <version>` |
| Preview changes       | `--dry-run`      |
