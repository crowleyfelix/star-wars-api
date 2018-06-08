use starWars;

//Collections
db.createCollection('planets');
db.createCollection('counters');

//Counters
db.counters.insert({
    "_id": "planet_id",
    "sequence_value": NumberInt(-1)
});

//Indexes
db.planets.createIndex({ name: "text" }, { unique: true })

