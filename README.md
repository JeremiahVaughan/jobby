Dependencies:
    - Install pg_dump
```sudo apt update && sudo apt install postgresql-client -y```

Config locations:
    - For local:
```./config.json```
    - For deployment:
```./deploy/config.json```

See status:
```sudo systemctl status jobby.service```

See logs:
```journalctl -u jobby.service```

