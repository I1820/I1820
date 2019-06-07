/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 04-06-2019
 * |
 * | File Name:     things_collection_index.js
 * +===============================================
 */
/* eslint-env mongo */

var collection = "things";

db[collection].createIndex({
  project: 1,
});

db[collection].createIndex({
  name: 1,
}, {
  unique: true,
});
