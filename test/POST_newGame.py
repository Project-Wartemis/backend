import requests

d = {"bots": ["QBot"], "gameEngineName": "demoEngine"}

resp = requests.post('http://localhost:80/game/new', json=d)
print(resp)
print(resp.text)