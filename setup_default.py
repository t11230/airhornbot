#!/usr/bin/env python2
import json
with open('config.json.default') as data:
    d = json.load(data)
    data.close()
token = raw_input("Please enter your Discord Token: ")
d['Token']=token
mongodb = raw_input("Please enter your MongoDB Connect String: ")
d['MongoDB']=mongodb
config = json.dumps(d)
config_file = open("config.json", "w")
config_file.write(config)
