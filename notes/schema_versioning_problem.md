# Schema versioning problem

When developing in a collaborative environment, it is posible that relying only in a timestamp based schema migration versioning might cause some problems.

An example would be, if a developer A makes a migration file on 2024-01-01 and a developer B makes a migration file on 2024-01-02, but developer B merges it's changes to main before developer A, then the order of the migrations won't match what really was tested.

To solve this, it is recommended to follow the next steps:

1. Down your local migration changes with goose:
```
goose <driver-ex.postgres> "<dbString-ex.host=localhost port=5432 user=baloo password=junglebook dbname=lenslocked sslmode=disable>" down
```
Use `down` as many times equal to the migrations you need to revert. You could also use `down-to VERSION` or `reset` depending your needs.
Note: to use goose properly, you need to be in the migrations folder or tell goose where is it with `-dir`.

2. In case you have work in progress, stash it:
```
git stash
```

3. Pull all the latest changes from main or the branch you want to push to:
```
git pull origin main --rebase
```

4. In the case that you stashed changes previously, you can now re-apply them:
```
git stash apply
```

5. Fix the schema versions with goose:
```
goose fix
```
This will substitute the timestamp versioning used by default with a fixed number as follows:
```
2024/04/13 16:35:11 RENAMED 20240413214125_users.sql => 00001_users.sql
2024/04/13 16:35:11 RENAMED 20240413214247_sessions.sql => 00002_sessions.sql
```

6. Finally you can Up your migrations and test the final state of your db before merging:
```
goose <driver> "<dbString>" up
goose <driver> "<dbString>" status
```
