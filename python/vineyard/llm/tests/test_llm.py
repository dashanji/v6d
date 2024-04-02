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

import torch

from vineyard.llm import KVCache
from vineyard.llm import KVTensor
from vineyard.llm.config import FileCacheConfig
from vineyard.llm.config import VineyardCacheConfig


def test_kv_cache_update_and_query_on_blob(vineyard_ipc_sockets):
    vineyard_cache_config = VineyardCacheConfig(
        socket=vineyard_ipc_sockets[0],
        block_size=5,
        sync_interval=3,
        llm_cache_sync_lock="llmCacheSyncLock",
        llm_cache_object_name="llm_cache_object",
        llm_ref_cnt_object_name="llm_refcnt_object",
    )
    cache = KVCache(
        cache_config=vineyard_cache_config,
        tensor_bytes=16,  # should be the same as the nbytes of the tensor
        cache_capacity=10,
        layer=2,
    )

    tokens = [1, 2, 3, 4]

    kv_tensors_to_update = []
    kv_tensors = []
    for _ in range(len(tokens)):
        k_tensor = torch.rand(2, 2, dtype=torch.float32)
        v_tensor = torch.rand(2, 2, dtype=torch.float32)
        kv_tensors.append([(k_tensor, v_tensor) for _ in range(cache.layer)])
        kv_tensors_to_update.append(
            [
                (
                    KVTensor(k_tensor.data_ptr(), k_tensor.nbytes),
                    KVTensor(v_tensor.data_ptr(), v_tensor.nbytes),
                )
                for _ in range(cache.layer)
            ]
        )

    # insert the token list and the related kv cache list
    updated = cache.update(None, tokens, kv_tensors_to_update)
    assert updated == len(tokens)

    kv_tensors_to_query = []
    kv_tensors_from_cache = []
    for _ in range(len(tokens)):
        kv_tensors_to_query.append(
            [
                (
                    KVTensor(0, 0),
                    KVTensor(0, 0),
                )
                for _ in range(cache.layer)
            ]
        )

    matched = cache.query(tokens, kv_tensors_to_query)
    kv_tensors_from_cache = kv_tensors_to_query[:matched]
    assert matched == len(tokens)

    assert len(kv_tensors) == len(kv_tensors_from_cache)
    for kv, kv_from_cache in zip(kv_tensors, kv_tensors_from_cache):
        assert len(kv) == len(kv_from_cache)
        for (k_tensor, v_tensor), (queried_k_tensor, queried_v_tensor) in zip(
            kv, kv_from_cache
        ):
            queried_k_tensor = torch.frombuffer(
                memoryview(queried_k_tensor),
                dtype=k_tensor.dtype,
                count=k_tensor.numel(),
            ).reshape(k_tensor.shape)
            queried_v_tensor = torch.frombuffer(
                memoryview(queried_v_tensor),
                dtype=v_tensor.dtype,
                count=v_tensor.numel(),
            ).reshape(v_tensor.shape)
            assert torch.equal(k_tensor, queried_k_tensor)
            assert torch.equal(v_tensor, queried_v_tensor)


def test_kv_cache_update_and_query_on_fs():
    file_cache_config = FileCacheConfig(
        chunk_size=2,
        split_number=2,
        root="/tmp/vineyard/llm_cache",
    )
    cache = KVCache(
        cache_config=file_cache_config,
        tensor_bytes=16,  # should be the same as the nbytes of the tensor
        cache_capacity=10,
        layer=2,
    )

    tokens = [1, 2, 3, 4]
    original_kv_tensors = []
    for i in range(0, len(tokens), file_cache_config.chunk_size):
        kv_tensors_to_update = []
        k_tensor = torch.rand(2, 2, dtype=torch.float32)
        v_tensor = torch.rand(2, 2, dtype=torch.float32)
        for _ in range(file_cache_config.chunk_size):
            original_kv_tensors.append(
                [(k_tensor, v_tensor) for _ in range(cache.layer)]
            )
            kv_tensors_to_update.append(
                [
                    (
                        KVTensor(k_tensor.numpy().ctypes.data, k_tensor.nbytes),
                        KVTensor(v_tensor.numpy().ctypes.data, v_tensor.nbytes),
                    )
                    for _ in range(cache.layer)
                ]
            )
        updated = cache.update(
            tokens[:i],
            tokens[i : i + file_cache_config.chunk_size],
            kv_tensors_to_update,
        )
        assert updated == file_cache_config.chunk_size

    kv_tensors_from_cache = []
    kv_tensors = []
    for _ in range(len(tokens)):
        k_tensor = torch.empty(2, 2, dtype=torch.float32)
        v_tensor = torch.empty(2, 2, dtype=torch.float32)
        kv_tensors_from_cache.append([(k_tensor, v_tensor) for _ in range(cache.layer)])
        kv_tensors.append(
            [
                (
                    KVTensor(k_tensor.numpy().ctypes.data, k_tensor.nbytes),
                    KVTensor(v_tensor.numpy().ctypes.data, v_tensor.nbytes),
                )
                for _ in range(cache.layer)
            ]
        )
    matched = cache.query(tokens, kv_tensors)
    assert matched == len(tokens)

    assert len(kv_tensors) == len(kv_tensors_from_cache)
    for kv, kv_from_cache in zip(original_kv_tensors, kv_tensors_from_cache):
        assert len(kv) == len(kv_from_cache)
        for (k_tensor, v_tensor), (queried_k_tensor, queried_v_tensor) in zip(
            kv, kv_from_cache
        ):
            assert torch.equal(k_tensor, queried_k_tensor)
            assert torch.equal(v_tensor, queried_v_tensor)
