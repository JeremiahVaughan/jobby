Dependencies:
    - Install pg_dump
```sudo apt update && sudo apt install postgresql-client -y```
    - Install sqlite3_rsync
        - on: 
            - Host of where the DB lives
            - Host preforming the backup operations
        - Place the binary at: ```/usr/local/bin```
        - references: 
            - Documentation: https://www.sqlite.org/rsync.html
            - Install From: https://www.sqlite.org/download.html
                - Unless Raspberry pi you will need to compile from source: https://til.simonwillison.net/sqlite/compile-sqlite3-rsync
    

Config locations:
    - For local:
```./config.json```
    - For deployment:
```./deploy/config.json```

See status:
```sudo systemctl status jobby.service```

See logs:
```journalctl -u jobby.service```

