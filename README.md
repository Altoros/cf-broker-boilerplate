# cf-broker-boilerplate

A boilerplate for Cloud Foundry Service Broker.

### How to run

```
go get github.com/Altoros/cf-broker-boilerplate
export DATABASE_URL=mysql://yser:passvord@localhost/dbname
cf-broker-boilerplate -c plans.yml
```

The app also can be deployed to Cloud Foundry (use [manifest.example.yml](https://github.com/Altoros/cf-broker-boilerplate/blob/master/manifest.example.yml) as an example of CF deployment manifest).

### TODO

[ ] add tests
[ ] extract service manager and binding managers to a separate packages
[ ] create a persister package that will encapsulate working with Mysql DB
