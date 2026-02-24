import asyncio
import grpc
from app.core.simulator import BoilerSimulator
from app.server.protos import boiler_pb2
from app.server.protos import boiler_pb2_grpc
class BoilerRPCServicer(boiler_pb2_grpc.BoilerPhysicsServicer):
    def __init__(self, smulator: BoilerSimulator):
        self._sim = smulator
    async def GetStatus(self, request, context):
        state = self._sim.get_state()
        return boiler_pb2.BoilerStatus(
            timestamp=state.timestamp,
            furnace_temp=state.outputs.furnace_temp,
            steam_pressure=state.outputs.steam_pressure,
            drum_level=state.outputs.drum_level,
            steam_flow=state.outputs.steam_flow,
        )
    async def SetControls(self, request, context):
        self._sim.set_controls(
            fuel=request.fuel_valve,
            water=request.feedwater_valve,
            steam=request.steam_valve,
        )
        return await self.GetStatus(request, context)
async def start_grpc_server(simulator: BoilerSimulator, port: int = 50051):
    server = grpc.aio.server()
    boiler_pb2_grpc.add_BoilerPhysicsServicer_to_server(
        BoilerRPCServicer(simulator), server
    )
    server.add_insecure_port(f"[::]:{port}")
    print(f"gRPC Server starting on port {port}...")
    await server.start()
    try:
        await server.wait_for_termination()
    except asyncio.CancelledError:
        print("gRPC Server stopping...")
        await server.stop(grace=None)
