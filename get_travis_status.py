import requests

response = requests.get('https://api.travis-ci.org/repos/t11230/ramenbot/cc.xml?branch=master')

print response.text
