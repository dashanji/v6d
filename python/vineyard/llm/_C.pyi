#! /usr/bin/env python
# -*- coding: utf-8 -*-
#
# Copyright 2020-2023 Alibaba Group Holding Limited.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

from collections.abc import Iterable
from typing import Any
from typing import Dict
from typing import List
from typing import Tuple
from typing import Union
from typing import overload

from vineyard.core.client import Client

# Define type aliases to match the C++ typedefs
KVState = Tuple[List[float], List[float]]
KVStateWithLayer = Dict[int, KVState]
ListKVStateWithLayer = List[KVStateWithLayer]

class KVStateCache:
    def __init__(
        self,
        dimension: int = 10,
        cacheCapacity: int = 10,
        layer: int = 10,
        blockSize: int = 5,
    ) -> None: ...
    @property
    def update(
        self,
        tokenList: List[int] = None,
        nextToken: int = None,
        kvState: KVStateWithLayer = None,
    ) -> None: ...

class KVStateCacheBuilder:
    def __init__(
        self,
        client: Client = None,
        dimension: int = 10,
        cacheCapacity: int = 10,
        layer: int = 10,
        blockSize: int = 5,
    ) -> None: ...
