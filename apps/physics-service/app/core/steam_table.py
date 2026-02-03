from typing import List, Tuple


class SteamTable:
    """
    Static lookup table for approximate steam pressure based on temperature.

    This is a simplified linear interpolation table for saturated steam pressure
    (approximation of real steam tables for simulation purposes).
    Data points are in °C and bar.
    """

    _DATA: List[Tuple[float, float]] = [
        (0.0, 0.0061),
        (20.0, 0.0234),
        (50.0, 0.1235),
        (99.0, 0.973),
        (100.0, 1.01325),
        (105.0, 1.208),
        (110.0, 1.433),
        (115.0, 1.691),
        (120.0, 1.985),
        (125.0, 2.321),
        (130.0, 2.701),
        (135.0, 3.132),
        (140.0, 3.614),
        (145.0, 4.155),
        (150.0, 4.760),
        (155.0, 5.432),
        (160.0, 6.180),
        (165.0, 7.008),
        (170.0, 7.920),
        (175.0, 8.924),
        (180.0, 10.027),
        (190.0, 12.55),
        (200.0, 15.549),
        (210.0, 19.077),
        (220.0, 23.198),
        (230.0, 28.005),
        (240.0, 33.480),
        (250.0, 39.776),
        (260.0, 46.940),
        (270.0, 55.09),
        (280.0, 64.21),
        (290.0, 74.53),
        (300.0, 85.93),
        (310.0, 99.63),
        (320.0, 112.9),
        (330.0, 128.2),
        (340.0, 145.5),
        (350.0, 165.4),
        (360.0, 187.8),
        (370.0, 213.0),
        (374.0, 221.2),
        (400.0, 250.0),
        (500.0, 300.0),
        (1000.0, 400.0),
        (1500.0, 500.0),
    ]

    @staticmethod
    def get_pressure(temp_c: float) -> float:
        """
        Returns approximate saturated steam pressure for given temperature using linear interpolation.

        Args:
            temp_c: Temperature in Celsius

        Returns:
            Pressure in bar (clamped to table range)

        Note: This is a linear approximation between data points.
              For temperatures below min or above max, returns edge values.
        """
        data = SteamTable._DATA

        if temp_c <= data[0][0]:
            return data[0][1]

        for i in range(len(data) - 1):
            t1, p1 = data[i]
            t2, p2 = data[i + 1]

            if t1 <= temp_c <= t2:
                ratio = (temp_c - t1) / (t2 - t1)
                pressure = p1 + (p2 - p1) * ratio
                return pressure

        return data[-1][1]
