# pgpac

Database deployment tool for PostgreSQL. This project is still a work in progress so there are a lot of core features that are unimplemented.

## Goal

There are a lot of tools for database deployments. I have used a few of these tools and some were really nice and some were not so nice. My goal is to implement the stuff that I like and I need and avoid the things I don't like and don't need.

### Things I like

- Tools that auto detect changes in the schema. This allows for a greatly simplified set of schema definition files making the database schema easier to maintain (IMO)

### Things I need

- Backup every bit of change so we don't accidentally lose data.
- Smart detection of seeding data changes
- Configuration of pre deploy and post deploy scripts for the purpose of oneshot type scripts

### Things I don't like

- Migration scripts for the schema
  - Arbitrary sql defining the schema easily becomes difficult to maintain
  - Makes it tough bringing up new databases since you have to run every migration script even though some of them won't have an effect on the ending database.
  - Requires ordering of scripts and this ordering is tedious to maintain when working with multiple people making changes to the schema. It isn't super difficult but it is annoying to resolve the issues.

### Things I don't need

- Reversible migrations
  - Never would reverse production databases because of data loss without extra tools for backing up database
- Change history
  - Already use git which tracks the change history

## Ideas

These are ideas that I have had for this tool

### Overall

Output the plan file for backup purposes. This could even include exporting data that will be deleted so the user can restore it later (add this as an option in the plan file cause this could cause the deployment time to be pretty slow without fancy optimizations (we probably could parallelize a lot of the deployment process pretty easily)). Maybe the backup can be stored on the database and then at the end of the deploy it will be backed up. I kinda really like that idea.

Control if a step should rollback everything and to what point it should rollback to. Config:

- RollbackStep [ALL, CURRENTSTEP, PREDEPLOY, DEPLOY, POSTDEPLOY, SEED] Default: ALL
  - ALL: rollsback everything that had been done during the deployment process
  - CURRENTSTEP: rollsback only the current step and halts the process
  - PREDEPLOY: rollback to before the predeploy step (same as ALL)
  - DEPLOY: rollback to before the deploy step
  - POSTDEPLOY: rollback to before the postdeploy step
  - SEED: rollback to before the seed step
  - If a step fails and it is before the rollback step then we fall back to rolling back the current step
    - i.e. the setting is "POSTDEPLOY" but the predeploy scripts fail. Then all the predeploy scripts will be rolled back

### Predeploy / Postdeploy scripts

Allow using some kind of macro system so you can modify the functions of the pre/post deploy script based on the resources and the different kinds of actions that will/have happen(ed). We actually know the whole plan at this point so predeploy scripts should be able to use almost all of the same macros as postdeploy.

Use cases:

- I want to be able to define an update trigger on every table without requiring me to write the same code over and over again.
- I want to modify some custom logging tables based on the tables and columns I am creating, altering, deleting, renaming, etc.during the predeploy and postdeploy steps.

Allow for running on specific databases (the macros actually solve this if we add a macro for the database these scripts are running on!!!).

Mark that a script should run before or after another particular script. (Need to make sure we disallow dependency cycles)

### Seeds

Allow for using both csvs or sql files of insert statements. Checks for adding, deleting, and updating (probably need per table configuration on this [OnlyAdd, NoDelete, All]). If there is an id in the table but not in the seed then we could automatically move it and update all the references to that id. Then we can maintain an exact match between seed data and the deployed table data. Maybe have a config for it (ConflictAction [AutoFix, Skip, Error] default Error or AutoFix)

### Schema

I think this is pretty straight forward. Just need to make sure we allow for organizing in folders but the schema plan still is one unit.
