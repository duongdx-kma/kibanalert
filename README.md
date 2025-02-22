# I. BUILD NEW LIBRARY FILES: 
### 1. install go and make sure go1.19

```bash
export GO_HOME=/usr/local/go1.19
export PATH=$GO_HOME/bin:$PATH

go version
go version go1.19.13 linux/amd64
```
### 2. install dependence and build binary file:
```bash
# install dependence
cd kibanaler/src
go mod tidy

# build binary file
cd kibanaler/src
go build -o ../bin/run .
```

# II. CONFIGURE FOR ELASTICSEARCH

### 1. create elastic index for save the log:
```bash
# I assume you use the following index mappings:

PUT /kibanalert/_mapping
{
    "properties": {
      "alert_id": {
        "type": "keyword"
      },
      "date": {
        "type": "date"
      },
      "reason": {
        "type": "text"
      },
      "rule_id": {
        "type": "keyword"
      },
      "service_name": {
        "type": "text"
      }
    }
}
```

### 2. Document to index:
```bash
# And that alerts populate index with this document template: 

{
  "alert_id": "{{alert.id}}",
  "rule_id": "{{rule.id}}",
  "reason": "The latency of service {{context.serviceName}} has reached {{context.threshold}} for {{context.interval}}. Current Value: {{context.triggerValue}}",
  "service_name": "{{context.serviceName}}",
  "date": "{{date}}"
}
```

### 3. Get Kibana active rules:
```bash
curl -X GET -u user:password "http://ip:5601/api/alerting/rules/_find?per_page=100"
```


### 4. get latest alert by rule_id:
```bash
GET kibanalert/_search
{
  "query": {
    "term": {
      "rule_id": "a96d2a00-edb2-11ef-980d-fd730a580b72"
    }
  },
  "size": 1,
  "sort": [{"date": {"order": "desc"}}]
}
```

### 5. create API-KEY for app:
```bash
POST /_security/api_key/grant
{
  "grant_type": "password",
  "username": "username",
  "password": "12345678",
  "api_key": {
    "name": "user-user-key"
  }
}
```

# III. GETTING STARTED
You only need the `./bin/run` and `./bin/.env` to use kibanalert.  Here's how:

### 1. create your own .env
```bash
cd ./bin
cp env.sample .env
```

### 2. fill your own .env
Fill out the details of `.env`.

### 3. run the app
```bash
Type `./run` to start kibanalert.
```