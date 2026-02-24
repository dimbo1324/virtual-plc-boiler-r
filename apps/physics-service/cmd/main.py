import asyncio
import sys
import os
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))
from app.core.simulator import BoilerSimulator
from app.server.grpc_server import start_grpc_server
async def physics_loop(sim: BoilerSimulator):
    print("Physics Engine started")
    try:
        while True:
            sim.tick()
            await asyncio.sleep(0.1)
    except asyncio.CancelledError:
        print("Physics Engine stopped")
async def main():
    sim = BoilerSimulator()
    physics_task = asyncio.create_task(physics_loop(sim))
    try:
        await start_grpc_server(sim, port=50051)
    except KeyboardInterrupt:
        pass
    finally:
        physics_task.cancel()
        await physics_task
if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        print("\nSee you next time!")
