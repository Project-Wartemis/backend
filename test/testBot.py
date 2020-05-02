#!/usr/bin/env python

# WS client example

import json
import asyncio
import websockets

async def register():
    uri = "ws://localhost:8080/register"
    async with websockets.connect(uri) as websocket:
        name = {
            "botName": "QBot"
            }

        name = json.dumps(name)
        await websocket.send(name)
        print(f"sending > {name}")

        resp = await websocket.recv()
        print(f"Server response: {resp}")

        access_key = json.loads(resp)["accessKey"]

        game_endpoint = await websocket.recv()
        print(f"Server response: {game_endpoint}")

        await play(game_endpoint, access_key)

async def play(game_endpoint, access_key):
    uri = "ws://localhost:8080" + game_endpoint
    async with websockets.connect(uri) as websocket:
        next_move = json.dumps(
            {
            "accessKey": access_key,
            "move": "{}"
            }
        )

        await websocket.send(next_move)
        print(f"> {next_move}")

        game_endpoint = await websocket.recv()
        print(f"< {game_endpoint}")


asyncio.get_event_loop().run_until_complete(register())