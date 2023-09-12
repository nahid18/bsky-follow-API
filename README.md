## Follow Multiple Handles on bsky.app API

### Endpoint URL

https://bsky-migrate.onrender.com/follow

### Request Payload Properties

| Property   | Description                          | Example                                     |
| ---------- | ------------------------------------ | ------------------------------------------- |
| `handle`   | Your bsky handle                     | `abdnahid`                                  |
| `password` | Your bsky password                   | `P@ssw0rd123`                               |
| `follow`   | Accounts to follow (comma-separated) | `stephaniehicks, anshulkundaje, jlsteenwyk` |

> [!NOTE]
> **Accepted `handle` formats**: `abdnahid.bsky.social` or `abdnahid`
> 
> **Accepted `follow` formats**: `account1, account2, account3` or `account1,account2,account3`

> [!WARNING]
> **Unaccepted `handle` format**: `@abdnahid.bsky.social`

### Using Python

```python
import requests
import json

url = "https://bsky-migrate.onrender.com/follow"
payload = {
    "handle": "abdnahid",
    "password": "your_password",
    "follow": "stephaniehicks,anshulkundaje,jlsteenwyk"
}
headers = {
    "Content-Type": "application/json"
}

response = requests.post(url, headers=headers, data=json.dumps(payload))
print(response.text)
```

### Using R
```R
library(httr)
library(jsonlite)

url <- "https://bsky-migrate.onrender.com/follow"
data <- list(
  handle = "your_handle",
  password = "your_password",
  follow = "stephaniehicks, anshulkundaje, jlsteenwyk"
)
json_data <- toJSON(data, auto_unbox = TRUE)
headers <- c(
  "Content-Type" = "application/json"
)
response <- POST(url, body = json_data, encode = "json", add_headers(headers))
content(response, "text")

```

### Acknowledgement

The `bsky` CLI used in this project is from here: 
https://github.com/mattn/bsky
