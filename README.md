## Follow Multiple Handles on bsky.app

### Endpoint URL

https://bsky-migrate.onrender.com/follow

### Request Payload Properties

| Property     | Description                              | Example                |
|--------------|------------------------------------------|------------------------|
| `handle`     | Your bsky handle                     | `abdnahid`           |
| `password`   | Your bsky password                   | `P@ssw0rd123`        |
| `follow`     | Accounts to follow (comma-separated)    | `stephaniehicks,anshulkundaje,jlsteenwyk"`  |

[!NOTE]
Accepted handle formats: abdnahid.bsky.social or abdnahid

[!WARNING]
Unaccepted handle format: @abdnahid.bsky.social

### Using Python

```python
import requests

url = "https://bsky-migrate.onrender.com/follow"
payload = {
    "handle": "abdnahid",
    "password": "XXXXXX",
    "follow": "stephaniehicks, anshulkundaje, jlsteenwyk"
}
response = requests.post(url, data=payload)
print(response.text)
```

### Using R
```R
library(httr)

url <- "https://bsky-migrate.onrender.com/follow"
payload <- list(
  handle = "your_basky_handle",
  password = "your_basky_password",
  follow = "account1, account2, account3"
)
response <- POST(url, body = payload)
content(response, "text")
```