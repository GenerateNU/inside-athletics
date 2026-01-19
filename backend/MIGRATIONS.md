
## Migration Guide

Migrations are a way we can manage our database. It serves as a log of how our database schema has changed over time.
Think of it as git logs but for our database! It makes it really easy to re-create our database, and rollback changes when
needed

### What tools our Migration workflow uses:
- [Atlas](https://atlasgo.io/): A really great tool we can use for migration generation and managment. It can also be used to store
migrations, but not completely necessary for a project of our scope. This is used to generate migrations automatically from our GORM models
- [GORM](https://gorm.io/index.html): Works with Atlas to make the migration generation process really easy. Atlas scans the repo for GORM models, identifies any changes in the schema using a temporary db, and creates a migration file for them

### Generating a Migration

All of our migrations live in the `./backend/internal/migrations` directory. These hold the information needed apply the changes to the database (or to roll them back)

To generate a migration run: `make gen-migration MIGRATION_NAME=[migration-name]`. This will create a new file in your migrations directory

Please use a meaningful name for the [migration-name] like "CreateUserTable". This is useful for easily being able to tell what the migration was for. The make target will not run if you do not give it a migration name. 

### Testing a Migration

NEVER apply a migration to prod. This will be done by the TL's after your code has been reviewed. In order to test your
migration run it on your local supabase database. This can be run with `make migration-dev`.


## Reverting Migrations

You can revert a migration by running: `make revert-single-prod` or `make revert-single-dev`. This will generate a plan so you can make sure you are reverting what you actually want.

In the makefile I have targets for reverting a single migration for both prod and dev. If you want to use 
any other optional flags like --dry-run (give you a preview of the changes) run them yourself. 

After reverting you must:

1) Delete the Migration: `atlas migrate rm [migration file name] --env dev`

example: `atlas migration rm 20260118001236_CreateUserTable --env dev`

2) Regenerate the hash file: `atlas migrate hash --env dev`
