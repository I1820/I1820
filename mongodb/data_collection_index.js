/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 04-06-2019
 * |
 * | File Name:     data_collection_index.js
 * +===============================================
 */
/* eslint-env mongo */

var collection = "data";

db[collection].createIndex({
  timestamp: -1,
});

db[collection].createIndex({
  thing_id: 1,
});

db[collection].createIndex({
  project: 1,
});
