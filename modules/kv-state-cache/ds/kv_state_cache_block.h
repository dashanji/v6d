/** Copyright 2020-2023 Alibaba Group Holding Limited.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

#ifndef MODULES_KV_STATE_CACHE_BLOCK_H_
#define MODULES_KV_STATE_CACHE_BLOCK_H_

#include <array>
#include <iostream>
#include <map>
#include <vector>

#include "basic/ds/tensor.h"
#include "client/ds/blob.h"
#include "client/ds/i_object.h"

typedef std::map<int, std::pair<std::vector<double>, std::vector<double>>>
    KV_STATE_WITH_LAYER;
typedef std::vector<
    std::map<int, std::pair<std::vector<double>, std::vector<double>>>>
    LIST_KV_STATE_WITH_LAYER;

// Set the bit to 1, which means the resource is not being used
#define FREE_BIT_RESOURCE(value, bit) ((value) |= (((uint64_t) 1) << (bit)))

// Set the bit to 0, which means the resource is being used
#define ACQUIRE_BIT_RESOURCE(value, bit) \
  ((value) &= (~(((uint64_t) 1) << (bit))))

struct offset_data {
  short offset;
};

namespace vineyard {
/**
 * If the kv-state cache block is full, the kv-state cache block will be
 * splited into two kv-state cache blocks.
*/
#define LIST_SIZE 5

/**
 * @brief KVStateCacheBlock is a cache for kv-cache of LLM. When a new prompt
 * comes, LLM can query KVStateCacheBlock to get the state of the kv-cache to
 * avoid caclulating the kv-cache again if the new prompt is similar to the
 * previous one.
 *
 * KVStateCacheBlock is stored in vineyard as a vineyard object which contains a
 * radix tree. The token sequence is the key of the radix tree and the value
 * point out the offset of the kv-cache in the tensor list.
 *
 * KVStateCacheBlock can be shared by multiple machines.
 */

class KVStateCacheBlock : public vineyard::Registered<KVStateCacheBlock> {
 private:
  std::shared_ptr<Tensor<double>> k_tensor;
  std::shared_ptr<Tensor<double>> v_tensor;
  std::vector<std::shared_ptr<KVStateCacheBlock>>
      child_kv_state_cache_block_list;
  uint64_t bitmap;
  ObjectID id;
  int dimension;

 public:
  static std::unique_ptr<Object> Create() __attribute__((used)) {
    return std::static_pointer_cast<Object>(
        std::unique_ptr<KVStateCacheBlock>{new KVStateCacheBlock()});
  }

  void Construct(const ObjectMeta& meta) override;

  std::string GetBitmapStr();

  uint64_t GetDimension() { return this->dimension; }

  uint64_t GetBitmap() { return this->bitmap; }

  std::shared_ptr<const Tensor<double>> GetKTensor() { return this->k_tensor; }

  std::shared_ptr<const Tensor<double>> GetVTensor() { return this->v_tensor; }

  friend class KVStateCacheBlockBuilder;
};

class KVStateCacheBlockBuilder : public ObjectBuilder {
 private:
  std::shared_ptr<TensorBuilder<double>> k_builder;
  std::shared_ptr<TensorBuilder<double>> v_builder;
  std::vector<KVStateCacheBlockBuilder*> child_kv_state_cache_builder_list;
  // TBD
  // support more than 64 kv-state cache slots
  uint64_t bitmap;
  pthread_spinlock_t spin_lock;
  int dimension;

  int FindEmptySlot();

 public:
  KVStateCacheBlockBuilder(Client& client, int dimension);

  KVStateCacheBlockBuilder(
      Client& client, std::shared_ptr<KVStateCacheBlock> kv_state_cache_block);

  /**
   * @brief Update the kv-state using next token.
   *
   * @param client The vineyard client.
   * @param kv_state The kv-state of the prompt. A LLM inference can contain
   * multiple kv-states for each layer.
   */
  void Update(const KV_STATE_WITH_LAYER& kv_state, offset_data* data);

  void Update(double* k_data, double* v_data, unsigned long data_length,
              offset_data* data);

  /**
   * @brief Query the kv-state using the whole token list.
   *
   * @param token_list The token list of the prompt.
   * @param token The token of the prompt.
   * @param kv_state The kv-state of the prompt returned by radix-tree. If the
   * kv-state is not found, the data of kv-state is invalid.
   */
  Status Query(Client& client, int index, KV_STATE_WITH_LAYER& kv_state);

  bool IsFull();

  Status Build(Client& client) override;

  std::shared_ptr<Object> _Seal(Client& client) override;

  void Lock() { pthread_spin_lock(&(this->spin_lock)); }

  void UnLock() { pthread_spin_unlock(&(this->spin_lock)); }

  const std::shared_ptr<TensorBuilder<double>> getKBuilder() {
    return k_builder;
  }

  const std::shared_ptr<TensorBuilder<double>> getVBuilder() {
    return v_builder;
  }

  void DeleteKVCache(int bit) { FREE_BIT_RESOURCE(this->bitmap, bit); }

  void SetChildKVStateCacheBlockBuilder(
      KVStateCacheBlockBuilder* child_kv_state_cache_builder);

  std::string GetBitmapStr();

  uint64_t GetBitmap() { return this->bitmap; }

  uint64_t GetDimension() { return this->dimension; }
};

}  // namespace vineyard

#endif