import asyncio
import time
from app.core.simulator import BoilerSimulator
from app.scenarios.demo import get_demo_controls


async def run_demo_console_simulation(
    simulator: BoilerSimulator,
    update_interval: float = 1.0,
    step_duration: float = 10.0,
):

    print("╔════════════════════════════════════════════╗")
    print("║      STARTING DEMO SIMULATION MODE         ║")
    print("║           Press Ctrl+C to stop             ║")
    print("╚════════════════════════════════════════════╝\n")

    simulator.set_controls(fuel=0.0, water=30.0, steam=0.0)

    start_time = time.time()
    last_print_time = start_time

    try:
        while True:
            simulator.tick()
            state = simulator.get_state()

            now = time.time()
            elapsed = now - start_time

            status, fuel, water, steam = get_demo_controls(elapsed, step_duration)

            if (fuel, water, steam) != (
                state.inputs.fuel_valve,
                state.inputs.feedwater_valve,
                state.inputs.steam_valve,
            ):
                simulator.set_controls(fuel=fuel, water=water, steam=steam)

            if now - last_print_time >= update_interval:
                print(
                    f"[{time.strftime('%H:%M:%S')}] "
                    f"{status:<12} | "
                    f"t={elapsed:>5.1f}s | "
                    f"Fuel {state.inputs.fuel_valve:>3.0f}% | "
                    f"Water {state.inputs.feedwater_valve:>3.0f}% | "
                    f"Steam {state.inputs.steam_valve:>3.0f}% | "
                    f"Temp {state.outputs.furnace_temp:>6.1f}°C | "
                    f"Press {state.outputs.steam_pressure:>5.1f} bar | "
                    f"Level {state.outputs.drum_level:>6.1f} mm | "
                    f"Flow {state.outputs.steam_flow:>5.1f} "
                )
                last_print_time = now

            await asyncio.sleep(0.1)

    except KeyboardInterrupt:
        raise
