#!/usr/bin/env python2
import json
with open('config.json.default') as data:
    d = json.load(data)
    data.close()

token = raw_input("Please enter your Discord Token: ")
d['Token']=token
mongodb = raw_input("Please enter your MongoDB Connect String: ")
d['MongoDB']=mongodb
adminyn = raw_input("Do you want admin enabled? (y/n): ")
if(adminyn=="Y" or adminyn=="y"):
    d['Modules'][0]['enable']= True
elif(adminyn=="N" or adminyn=="n"):
    d['Modules'][0]['enable']= False
gamblingyn = raw_input("Do you want gambling enabled? (y/n): ")
if(gamblingyn=="Y" or gamblingyn=="y"):
    d['Modules'][1]['enable']= True
elif(gamblingyn=="N" or gamblingyn=="n"):
    d['Modules'][1]['enable']= False
greeteryn = raw_input("Do you want greeter enabled? (y/n): ")
if(greeteryn=="Y" or greeteryn=="y"):
    d['Modules'][2]['enable']= True
elif(greeteryn=="N" or greeteryn=="n"):
    d['Modules'][2]['enable']= False
helpyn = raw_input("Do you want help enabled? (y/n): ")
if(helpyn=="Y" or helpyn=="y"):
    d['Modules'][3]['enable']= True
elif(helpyn=="N" or helpyn=="n"):
    d['Modules'][3]['enable']= False
rolemodyn = raw_input("Do you want role modification enabled? (y/n): ")
if(rolemodyn=="Y" or rolemodyn=="y"):
    d['Modules'][4]['enable']= True
elif(rolemodyn=="N" or rolemodyn=="n"):
    d['Modules'][4]['enable']= False
soundboardyn = raw_input("Do you want soundboard enabled? (y/n): ")
if(soundboardyn=="Y" or soundboardyn=="y"):
    d['Modules'][5]['enable']= True
elif(soundboardyn=="N" or soundboardyn=="n"):
    d['Modules'][5]['enable']= False
voicebonusyn = raw_input("Do you want voicebonus enabled? (y/n): ")
if(voicebonusyn=="Y" or voicebonusyn=="y"):
    d['Modules'][6]['enable']= True
elif(voicebonusyn=="N" or voicebonusyn=="n"):
    d['Modules'][6]['enable']= False

config = json.dumps(d)
config_file = open("config.json", "w")
config_file.write(config)
