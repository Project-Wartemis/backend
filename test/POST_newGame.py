import requests

d = {"bots": ["QBot"], "gameEngineName": "demoEngine"}

resp = requests.post('http://localhost:8080/game/new', json=d)
print(resp)
print(resp.text)