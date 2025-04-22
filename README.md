Ensure local ~/.ssh/config file contains an entry for "deploy.target"

Deploy:
```./deploy.sh```

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
    
Config
    Upload with:
```aws s3 cp <local_file_path> s3://<bucket_name>/<object_key>```
    Download with:
```aws s3 cp s3://<bucket_name>/<object_key> <local_file_path>```

Config locations:
    Config is grabbed from s3 but make sure you emplace s3 files at:
```/root/.aws/config```
```/root/.aws/credentials```

See status:
```sudo systemctl status jobby.service```

See logs:
```journalctl -u jobby.service```

