import asyncio
import sys
import os

sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))

from app.core.simulator import BoilerSimulator


async def run_simulation():
    sim = BoilerSimulator()
    print("--- ЗАПУСК СИМУЛЯЦИИ КОТЛА ---")
    print("Нажмите Ctrl+C для выхода")

    sim.set_controls(fuel=0.0, water=30.0, steam=0.0)

    try:
        start_time = asyncio.get_running_loop().time()

        while True:
            sim.tick()
            state = sim.get_state()

            current_loop_time = asyncio.get_running_loop().time() - start_time

            if 0 < current_loop_time < 10:
                sim.set_controls(fuel=50.0, water=10.0, steam=0.0)
                status = "РАЗОГРЕВ"
            elif 10 < current_loop_time < 20:
                sim.set_controls(fuel=100.0, water=20.0, steam=0.0)
                status = "ФОРСАЖ"
            elif 20 < current_loop_time < 30:
                sim.set_controls(fuel=100.0, water=25.0, steam=80.0)
                status = "ОТДАЧА ПАРА"
            else:
                sim.set_controls(fuel=0.0, water=50.0, steam=100.0)
                status = "ОСТАНОВ"

            output = (
                f"STATUS: {status:<10} | "
                f"Time: {current_loop_time:.1f}s | "
                f"Fuel: {state.inputs.fuel_valve:>3.0f}% | "
                f"Temp: {state.outputs.furnace_temp:>6.1f}°C | "
                f"Press: {state.outputs.steam_pressure:>5.1f} Bar | "
                f"Level: {state.outputs.drum_level:>5.1f} mm | "
                f"Flow: {state.outputs.steam_flow:>4.1f}"
            )
            print(output)

            await asyncio.sleep(0.5)

    except KeyboardInterrupt:
        print("\n--- СИМУЛЯЦИЯ ОСТАНОВЛЕНА ---")


if __name__ == "__main__":
    try:
        asyncio.run(run_simulation())
    except KeyboardInterrupt:
        pass
