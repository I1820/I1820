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


db.data.insert({
  raw: new BinData(0, "omMxMDAZF3BjMTAxGRnO"),
  data: {
    101: 6606,
    100: 6000
  },
  timestamp: new ISODate("2018-05-07T05:49:54.415Z"),
  thingid: "0000000000000001",
  rxinfo: [
    {
      mac: "b827ebffff633260",
      name: "isrc-gateway",
      time: new ISODate("2018-05-07T05:49:53.874Z"),
      rssi: -57,
      lorasnr: 10
    }
  ],
  txinfo: {
    frequency: 868100000,
    adr: true,
    coderate: "4/6"
  },
  project: "hello"

}, {
  raw: new BinData(0, "omMxMDAZF3BjMTAxGRnO"),
  data: null,
  timestamp: new ISODate("2018-05-06T18:40:03.156Z"),
  thingid: "0000000000000010",
  rxinfo: [
    {
      mac: "b827ebffff633260",
      name: "isrc-gateway",
      time: new ISODate("2018-05-06T18:40:03.151Z"),
      rssi: -57,
      lorasnr: 10
    }
  ],
  txinfo: {
    frequency : 868100000,
    adr : true,
    coderate : "4/6"
  },
  project : "hello"
});
