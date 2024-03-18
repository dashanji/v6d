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

#include "llm-cache/storage/file_storage.h"

namespace vineyard {

FileStorage::FileStorage(int chunkSize){
    this->chunkSize = chunkSize;
}

FileStorage::~FileStorage() {
  // TBD
}

Status FileStorage::Update(const std::vector<int>& tokenList, int nextToken,
                   const std::map<int, std::pair<LLMKV, LLMKV>>& kvState) {
    // not implemented
    return Status::NotImplemented();
}

Status FileStorage::Update(const std::vector<int>& tokenList,
                   const std::vector<std::map<int, std::pair<LLMKV, LLMKV>>>&
                       kvStateList) {
    
    return Status::OK();
}

Status FileStorage::Query(const std::vector<int>& tokenList, int token,
                  std::map<int, std::pair<LLMKV, LLMKV>>& kvState) {
    // not implemented
    return Status::NotImplemented();
}

Status FileStorage::Query(const std::vector<int>& tokenList,
                  std::vector<std::map<int, std::pair<LLMKV, LLMKV>>>& kvStateList) {
    return Status::OK();
}

}  // namespace vineyard