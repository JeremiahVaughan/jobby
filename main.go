package main

import (
    "time"
    "flag"
)


func main() {
    configPath := flag.String("config/config.json", "location of config file")
    flag.Parse()
    
    c := config.New(*configPath) 
    
}
