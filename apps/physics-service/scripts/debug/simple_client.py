import asyncio
import os
import sys
import grpc


sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "../..")))

from app.server.protos import boiler_pb2
from app.server.protos import boiler_pb2_grpc


async def run():
    async with grpc.aio.insecure_channel("localhost:50051") as channel:
        stub = boiler_pb2_grpc.BoilerPhysicsStub(channel)
        print("--- 1. Read state (GetStatus) ---")
        response = await stub.GetStatus(boiler_pb2.Empty())
        print(
            f"Temp: {response.furnace_temp:.2f}°C, Pressure: {response.steam_pressure:.2f} Bar"
        )
        print("\n--- 2. Send Commands (SetControls) ---")

        new_state = await stub.SetControls(
            boiler_pb2.ControlInput(
                fuel_valve=100.0, feedwater_valve=50.0, steam_valve=0.0
            )
        )

        print(f"Command sent. System responding...")
        print("\n--- 3. Observe (3 sec) ---")

        for _ in range(3):
            await asyncio.sleep(1.0)
            response = await stub.GetStatus(boiler_pb2.Empty())
            print(f"Temp: {response.furnace_temp:.2f}°C (Calculate!)")


if __name__ == "__main__":
    asyncio.run(run())
