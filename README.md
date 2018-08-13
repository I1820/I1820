# Data Manager

## Introduction
DM queries and returns data from database (mongodb).
it has grafana plugin for better data management.

## Profiler
Enable MongoDB buit-in profiler:

```
use isrc
db.setProfileLevel(2)
```

The profiling results in a special capped collection called `system.profile`
which is located in the database where you executed the `setProfileLevel` command.

```
db.system.profile.find().pretty()
```
