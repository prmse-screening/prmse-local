from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class InferenceRequest(_message.Message):
    __slots__ = ("series", "path", "model", "cpu")
    SERIES_FIELD_NUMBER: _ClassVar[int]
    PATH_FIELD_NUMBER: _ClassVar[int]
    MODEL_FIELD_NUMBER: _ClassVar[int]
    CPU_FIELD_NUMBER: _ClassVar[int]
    series: str
    path: str
    model: str
    cpu: bool
    def __init__(self, series: _Optional[str] = ..., path: _Optional[str] = ..., model: _Optional[str] = ..., cpu: bool = ...) -> None: ...

class InferenceResponse(_message.Message):
    __slots__ = ("series", "result")
    SERIES_FIELD_NUMBER: _ClassVar[int]
    RESULT_FIELD_NUMBER: _ClassVar[int]
    series: str
    result: str
    def __init__(self, series: _Optional[str] = ..., result: _Optional[str] = ...) -> None: ...
