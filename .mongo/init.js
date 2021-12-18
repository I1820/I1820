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

const projectID = "el-project";
const collection = "data";

db = db.getSiblingDB("i1820");

db[collection].insert([
  {
    raw: new BinData(0, "o2MxMDAZG1hjMTAxGRnOZWNvdW50YjE3"),
    data: {
      100: 7000,
      101: 6606,
      count: "17",
    },
    timestamp: new ISODate("2018-11-11T09:47:20.902Z"),
    thingid: "0000000000000073",
    rxinfo: [
      {
        mac: "b827ebffff70c80a",
        name: "5cf0f70f7064c500094b5e31",
        time: new ISODate("0001-01-01T00:00:00Z"),
        rssi: new NumberLong(-57),
        lorasnr: 7,
      },
    ],
    txinfo: {
      frequency: new NumberLong(868300000),
      adr: false,
      coderate: "4/5",
    },
    project: projectID,
  },
  {
    raw: new BinData(0, "o2MxMDAZG1hjMTAxGRnOZWNvdW50YjE3"),
    data: {
      100: 7000,
      101: 6606,
      count: "18",
    },
    timestamp: new ISODate("2018-11-11T09:47:21.902Z"),
    thingid: "0000000000000073",
    rxinfo: [
      {
        mac: "b827ebffff70c80a",
        name: "5cf0f70f7064c500094b5e31",
        time: new ISODate("0001-01-01T00:00:00Z"),
        rssi: new NumberLong(-57),
        lorasnr: 7,
      },
    ],
    txinfo: {
      frequency: new NumberLong(868300000),
      adr: false,
      coderate: "4/5",
    },
    project: projectID,
  },
  {
    raw: new BinData(0, "o2MxMDAZG1hjMTAxGRnOZWNvdW50YjE3"),
    data: {
      100: 7000,
      101: 6606,
      count: "19",
    },
    timestamp: new ISODate("2018-11-11T09:47:22.902Z"),
    thingid: "0000000000000073",
    rxinfo: [
      {
        mac: "b827ebffff70c80a",
        name: "5cf0f70f7064c500094b5e31",
        time: new ISODate("0001-01-01T00:00:00Z"),
        rssi: new NumberLong(-57),
        lorasnr: 7,
      },
    ],
    txinfo: {
      frequency: new NumberLong(868300000),
      adr: false,
      coderate: "4/5",
    },
    project: projectID,
  },
  {
    raw: new BinData(0, "o2MxMDAZG1hjMTAxGRnOZWNvdW50YjE3"),
    data: null,
    timestamp: new ISODate("2018-11-11T09:47:23.902Z"),
    thingid: "0000000000000073",
    rxinfo: [
      {
        mac: "b827ebffff70c80a",
        name: "5cf0f70f7064c500094b5e31",
        time: new ISODate("0001-01-01T00:00:00Z"),
        rssi: new NumberLong(-57),
        lorasnr: 7,
      },
    ],
    txinfo: {
      frequency: new NumberLong(868300000),
      adr: false,
      coderate: "4/5",
    },
    project: projectID,
  },
]);
