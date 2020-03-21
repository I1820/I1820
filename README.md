[![Drone (cloud)](https://img.shields.io/drone/build/I1820/I1820.svg?style=flat-square)](https://cloud.drone.io/I1820/I1820)

## Link
Link component of I1820 platfrom. This service collects
raw data from bottom layer (protocols), stores them into mongo database
and decodes them using user's selected decoder.
This service also sends data into bottom layer (protocols) after
encoding them using user's selected encoder.

Link uses MQTT for communicating with the bottom layer and this communication can be customized
using Protocol's interface which is defined in `protocols/protocol.go`.

## Thing Manager
Thing manager manages I1820 Things and their properties.
Things belong to the projects, but this component doesn't validate this relationship so other services
must verify project identification and existence before calls this project APIs.

## Data Manager
DM is a Data Manager component of the I1820 platform.
It has some useful built-in queries that can returns data from the database (MongoDB) to the API backend.

### Profiler
Enable MongoDB built-in profiler:

```
use i1820
db.setProfileLevel(2)
```

The profiling results will be in a special capped collection called `system.profile`
which is located in the database where you executed the `setProfileLevel` command.

```
db.system.profile.find().pretty()
```

