# Eventually consistent

Using hydra with RethinkDB implies eventual consistency on all endpoints, except `/oauth2/auth` and `/oauth2/token`.
Eventual consistent data is usually not immediately available. This is dependent on the network latency between Hydra
and RethinkDB.