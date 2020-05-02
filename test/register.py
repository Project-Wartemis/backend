#!/usr/bin/env python

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

        response = await websocket.recv()
        print(f"Server response: {response}")

asyncio.get_event_loop().run_until_complete(hello())