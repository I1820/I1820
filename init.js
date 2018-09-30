/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 28-04-2018
 * |
 * | File Name:     init.js
 * +===============================================
 */
/* eslint-env mongo */

var thingID = "5ba3f1a287a142b0a840fae1"
var projectID = "5ba3f19c87a142b0a840fae0"
var collection = "data." + projectID +  "." + thingID

db[collection].insert([{
  raw: 7000,
  value: {
    number: 7000,
  },
  at: new ISODate("2018-09-26T21:52:06.443+03:30"),
  project: projectID,
  thing_id: thingID,
  asset: "100"
}, {
  raw: 6500,
  value: {
    number: 6500,
  },
  at: new ISODate("2018-09-26T21:52:07.443+03:30"),
  project: projectID,
  thing_id: thingID,
  asset: "100"
}, {
  raw: 7000,
  value: {
    number: 7000,
  },
  at: new ISODate("2018-09-26T21:52:08.443+03:30"),
  project: projectID,
  thing_id: thingID,
  asset: "100"
}, {
  raw: 6500,
  value: {
    number: 6500,
  },
  at: new ISODate("2018-09-26T21:52:09.443+03:30"),
  project: projectID,
  thing_id: thingID,
  asset: "100"
}, {
  raw: 6500,
  value: {
    number: 6500,
  },
  at: new ISODate("2018-09-26T21:52:09.443+03:30"),
  project: projectID,
  thing_id: thingID,
  asset: "101"
}, {
  raw: "hello",
  value: {
    string: "hello",
  },
  at: new ISODate("2018-09-26T21:52:09.443+03:30"),
  project: projectID,
  thing_id: thingID,
  asset: "101"
}]);
