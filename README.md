# KI
Simple CRUD service for storing and retrieving configuration files. Currently meant as a reference project. Not fit for any use as is.

Service used to store config data. Each config has to belong tou a group. The group object consists of an id and an array of
config ids. A config consists of some metadata fields and a property field containing the config data. The properties can be
arbitrary json data.

Both resources support PUT and GET.

Group URL: /config/{groupId}  
Config URL: /config/{groupId}/{configId}

Group example:
```
{
    "id": "someGroup",
    "configs": [
        "someId"
    ]
}
```

Config example:
```
{
  "id": "someId",
  "name": "someName",
  "lastModified": "2018-09-02T14:53:56.281992009+02:00",
  "version": 0,
  "group": "someGroup",
  "properties": {
    "property1": 12,
    "property2": "12",
    "property3": "someString",
    "property4": "someOtherString",
    "property5": 12.1
  }
}
```
