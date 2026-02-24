from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional
DESCRIPTOR: _descriptor.FileDescriptor
class Empty(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...
class BoilerStatus(_message.Message):
    __slots__ = ("timestamp", "furnace_temp", "steam_pressure", "drum_level", "steam_flow")
    TIMESTAMP_FIELD_NUMBER: _ClassVar[int]
    FURNACE_TEMP_FIELD_NUMBER: _ClassVar[int]
    STEAM_PRESSURE_FIELD_NUMBER: _ClassVar[int]
    DRUM_LEVEL_FIELD_NUMBER: _ClassVar[int]
    STEAM_FLOW_FIELD_NUMBER: _ClassVar[int]
    timestamp: float
    furnace_temp: float
    steam_pressure: float
    drum_level: float
    steam_flow: float
    def __init__(self, timestamp: _Optional[float] = ..., furnace_temp: _Optional[float] = ..., steam_pressure: _Optional[float] = ..., drum_level: _Optional[float] = ..., steam_flow: _Optional[float] = ...) -> None: ...
class ControlInput(_message.Message):
    __slots__ = ("fuel_valve", "feedwater_valve", "steam_valve")
    FUEL_VALVE_FIELD_NUMBER: _ClassVar[int]
    FEEDWATER_VALVE_FIELD_NUMBER: _ClassVar[int]
    STEAM_VALVE_FIELD_NUMBER: _ClassVar[int]
    fuel_valve: float
    feedwater_valve: float
    steam_valve: float
    def __init__(self, fuel_valve: _Optional[float] = ..., feedwater_valve: _Optional[float] = ..., steam_valve: _Optional[float] = ...) -> None: ...
