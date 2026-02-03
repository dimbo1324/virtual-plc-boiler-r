import asyncio
import sys
import os

sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))

from app.core.simulator import BoilerSimulator
from app.runners.console import run_demo_console_simulation


async def main():
    simulator = BoilerSimulator()
    await run_demo_console_simulation(simulator)


if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        print("\nSIMULATION STOPPED")
    except Exception as e:
        print(f"Critical error: {e}", file=sys.stderr)
        sys.exit(1)
