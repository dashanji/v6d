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

#include "client/rpc_client.h"

#include <iostream>
#include <map>
#include <mutex>
#include <set>
#include <string>
#include <unordered_map>
#include <unordered_set>
#include <utility>
#include <vector>

#include "client/ds/blob.h"
#include "client/ds/object_factory.h"
#include "client/ds/remote_blob.h"
#include "client/io.h"
#include "client/utils.h"
#include "common/compression/compressor.h"
#include "common/util/env.h"
#include "common/util/logging.h"
#include "common/util/protocols.h"
#include "common/util/version.h"

namespace vineyard {

RPCClient::~RPCClient() {
#ifdef VINEYARD_WITH_RDMA
  Status status = StopRDMA();
  if (!status.ok()) {
    std::cout << "Failed to stop RDMA client: " << status.ToString()
              << ". May cause resource leak." << std::endl;
  }
#endif
  Disconnect();
}

void RPCClient::Disconnect() {
#ifdef VINEYARD_WITH_RDMA
  Status status = StopRDMA();
  if (!status.ok()) {
    std::cout << "Failed to stop RDMA client: " << status.ToString()
              << ". May cause resource leak." << std::endl;
  }
#endif
  std::lock_guard<std::recursive_mutex> __guard(this->client_mutex_);
  if (!this->connected_) {
    return;
  }
  std::string message_out;
  WriteExitRequest(message_out);
  VINEYARD_SUPPRESS(doWrite(message_out));
  close(vineyard_conn_);
  connected_ = false;
}

Status RPCClient::Connect() {
  auto ep = read_env("VINEYARD_RPC_ENDPOINT");
  auto rdma_endpoint = read_env("VINEYARD_RDMA_ENDPOINT");
  if (!ep.empty()) {
    return Connect(ep, rdma_endpoint);
  }
  return Status::ConnectionError(
      "Environment variable VINEYARD_RPC_ENDPOINT does't exists");
}

Status RPCClient::Connect(std::string const& username,
                          std::string const& password) {
  auto ep = read_env("VINEYARD_RPC_ENDPOINT");
  auto rdma_ep = read_env("VINEYARD_RDMA_ENDPOINT");
  if (!ep.empty()) {
    return Connect(ep, username, password, rdma_ep);
  }
  return Status::ConnectionError(
      "Environment variable VINEYARD_RPC_ENDPOINT does't exists");
}

Status RPCClient::Connect(const std::string& rpc_endpoint) {
  auto rdma_ep = read_env("VINEYARD_RDMA_ENDPOINT");
  return this->Connect(rpc_endpoint, RootSessionID(), "", "", rdma_ep);
}

Status RPCClient::Connect(const std::string& rpc_endpoint,
                          std::string const& username,
                          std::string const& password,
                          const std::string& rdma_endpoint) {
  return this->Connect(rpc_endpoint, RootSessionID(), username, password,
                       rdma_endpoint);
}

Status RPCClient::Connect(const std::string& rpc_endpoint,
                          const SessionID session_id,
                          std::string const& username,
                          std::string const& password,
                          const std::string& rdma_endpoint) {
  size_t pos = rpc_endpoint.find(":");
  std::string host, port;
  if (pos == std::string::npos) {
    host = rpc_endpoint;
    port = "9600";
  } else {
    host = rpc_endpoint.substr(0, pos);
    port = rpc_endpoint.substr(pos + 1);
  }

  std::string rdma_host = "", rdma_port = "-1";
  pos = rdma_endpoint.find(":");
  if (pos == std::string::npos) {
    std::cout << "No RDMA endpoint provided. Fall back to TCP." << std::endl;
  } else {
    rdma_host = rdma_endpoint.substr(0, pos);
    rdma_port = rdma_endpoint.substr(pos + 1);
  }

  return this->Connect(host, static_cast<uint32_t>(std::stoul(port)),
                       session_id, username, password, rdma_host,
                       static_cast<uint32_t>(std::stoul(rdma_port)));
}

Status RPCClient::Connect(const std::string& host, uint32_t port,
                          const std::string& rdma_host, uint32_t rdma_port) {
  return this->Connect(host, port, RootSessionID(), "", "", rdma_host,
                       rdma_port);
}

Status RPCClient::Connect(const std::string& host, uint32_t port,
                          std::string const& username,
                          std::string const& password,
                          const std::string& rdma_host, uint32_t rdma_port) {
  return this->Connect(host, port, RootSessionID(), username, password,
                       rdma_host, rdma_port);
}

Status RPCClient::Connect(const std::string& host, uint32_t port,
                          const SessionID session_id,
                          std::string const& username,
                          std::string const& password,
                          const std::string& rdma_host, uint32_t rdma_port) {
  std::lock_guard<std::recursive_mutex> guard(client_mutex_);
  std::string rpc_endpoint = host + ":" + std::to_string(port);
  RETURN_ON_ASSERT(!connected_ || rpc_endpoint == rpc_endpoint_);
  if (connected_) {
    return Status::OK();
  }
  rpc_endpoint_ = rpc_endpoint;
  RETURN_ON_ERROR(connect_rpc_socket_retry(host, port, vineyard_conn_));
  std::string message_out;
  WriteRegisterRequest(message_out, StoreType::kDefault, session_id, username,
                       password);
  RETURN_ON_ERROR(doWrite(message_out));
  json message_in;
  RETURN_ON_ERROR(doRead(message_in));
  std::string ipc_socket_value, rpc_endpoint_value;
  bool store_match;
  RETURN_ON_ERROR(ReadRegisterReply(
      message_in, ipc_socket_value, rpc_endpoint_value, remote_instance_id_,
      session_id_, server_version_, store_match, support_rpc_compression_));
  ipc_socket_ = ipc_socket_value;
  connected_ = true;
  set_compression_enabled(support_rpc_compression_);

  if (!compatible_server(server_version_)) {
    std::clog << "[warn] Warning: this version of vineyard client may be "
                 "incompatible with connected server: "
              << "client's version is " << vineyard_version()
              << ", while the server's version is " << server_version_
              << std::endl;
  }

  // RPC client doesn't have a concrete instance id, even the unspecified
  // instance id.
  instance_id_ = UnspecifiedInstanceID() - 1;

  if (rdma_host.length() > 0) {
    Status status = ConnectRDMA(rdma_host, rdma_port);
    if (status.ok()) {
      rdma_endpoint_ = rdma_host + ":" + std::to_string(rdma_port);
      std::cout << "Connected to RPC server: " << rpc_endpoint
                << ", RDMA server: " << rdma_host << ":" << rdma_port
                << std::endl;
    } else {
      std::cout << "Connect RDMA server failed! Fall back to RPC mode. Error:"
                << status.message() << std::endl;
    }
  }

  return Status::OK();
}

Status RPCClient::ConnectRDMA(const std::string& rdma_host,
                              uint32_t rdma_port) {
#ifdef VINEYARD_WITH_RDMA
  if (this->rdma_connected_) {
    return Status::OK();
  }
  RETURN_ON_ERROR(RDMAClientCreator::Create(this->rdma_client_, rdma_host,
                                            static_cast<int>(rdma_port)));
  std::cout << "Try to connect to RDMA server " << rdma_host << ":" << rdma_port
            << "..." << std::endl;

  RETURN_ON_ERROR(this->rdma_client_->Connect());
  RETURN_ON_ERROR(RDMAExchangeMemInfo());
  this->rdma_connected_ = true;
  return Status::OK();
#else
  return Status::NotImplemented("RDMA is not supported in this build.");
#endif
}

Status RPCClient::RDMAExchangeMemInfo() {
#ifdef VINEYARD_WITH_RDMA
  void* buffer;
  this->rdma_client_->GetTXFreeMsgBuffer(buffer);
  VineyardMsg* msg = reinterpret_cast<VineyardMsg*>(buffer);
  msg->type = VINEYARD_MSG_EXCHANGE_KEY;
  msg->remoteMemInfo.remote_address = -1;
  msg->remoteMemInfo.key = -1;
  msg->remoteMemInfo.len = -1;
  std::cout << "Send remote addr: "
            << reinterpret_cast<void*>(msg->remoteMemInfo.remote_address)
            << ", rkey: " << msg->remoteMemInfo.key << std::endl;
  void* remoteMsg;
  this->rdma_client_->GetRXFreeMsgBuffer(remoteMsg);
  memset(remoteMsg, 0, 64);
  VINEYARD_CHECK_OK(
      this->rdma_client_->Recv(remoteMsg, sizeof(VineyardMsg), nullptr));
  this->rdma_client_->Send(buffer, sizeof(VineyardMsg), nullptr);
  this->rdma_client_->GetTXCompletion(-1, nullptr);

  VINEYARD_CHECK_OK(rdma_client_->GetRXCompletion(-1, nullptr));

  VineyardMsg* vmsg = reinterpret_cast<VineyardMsg*>(remoteMsg);
  if (vmsg->type == VINEYARD_MSG_EXCHANGE_KEY) {
    remote_info_.address = vmsg->remoteMemInfo.remote_address;
    remote_info_.rkey = vmsg->remoteMemInfo.key;
    std::cout << "Get remote address: " << remote_info_.address
              << ", rkey: " << remote_info_.rkey << std::endl;
  } else {
    std::cout << "Unknown message type: " << vmsg->type << std::endl;
  }
#endif
  return Status::OK();
}

Status RPCClient::StopRDMA() {
#ifdef VINEYARD_WITH_RDMA
  if (!rdma_connected_) {
    return Status::OK();
  }
  rdma_connected_ = false;

  void* msg;
  RETURN_ON_ERROR(rdma_client_->GetTXFreeMsgBuffer(msg));
  VineyardMsg* vmsg = reinterpret_cast<VineyardMsg*>(msg);
  vmsg->type = VINEYARD_MSG_CLOSE;
  RETURN_ON_ERROR(rdma_client_->Send(msg, sizeof(VineyardMsg), nullptr));
  RETURN_ON_ERROR(rdma_client_->GetTXCompletion(-1, nullptr));

  RETURN_ON_ERROR(rdma_client_->Stop());
  RETURN_ON_ERROR(rdma_client_->Close());

  {
    std::lock_guard<std::mutex> guard(local_infos_mutex_);
    local_infos_.clear();
  }
#endif

  return Status::OK();
}

Status RPCClient::Fork(RPCClient& client) {
  RETURN_ON_ASSERT(!client.Connected(),
                   "The client has already been connected to vineyard server");
  return client.Connect(rpc_endpoint_, session_id_, rdma_endpoint_);
}

Status RPCClient::GetMetaData(const ObjectID id, ObjectMeta& meta,
                              const bool sync_remote) {
  ENSURE_CONNECTED(this);
  json tree;
  RETURN_ON_ERROR(GetData(id, tree, sync_remote));
  meta.Reset();
  meta.SetMetaData(this, tree);
  return Status::OK();
}

Status RPCClient::GetMetaData(const std::vector<ObjectID>& ids,
                              std::vector<ObjectMeta>& metas,
                              const bool sync_remote) {
  ENSURE_CONNECTED(this);
  std::vector<json> trees;
  RETURN_ON_ERROR(GetData(ids, trees, sync_remote));
  metas.resize(trees.size());

  for (size_t idx = 0; idx < trees.size(); ++idx) {
    metas[idx].Reset();
    metas[idx].SetMetaData(this, trees[idx]);
  }
  return Status::OK();
}

std::shared_ptr<Object> RPCClient::GetObject(const ObjectID id) {
  ObjectMeta meta;
  RETURN_NULL_ON_ERROR(this->GetMetaData(id, meta, true));
  RETURN_NULL_ON_ASSERT(!meta.MetaData().empty());
  auto object = ObjectFactory::Create(meta.GetTypeName());
  if (object == nullptr) {
    object = std::unique_ptr<Object>(new Object());
  }
  object->Construct(meta);
  return object;
}

Status RPCClient::GetObject(const ObjectID id,
                            std::shared_ptr<Object>& object) {
  ObjectMeta meta;
  RETURN_ON_ERROR(this->GetMetaData(id, meta, true));
  RETURN_ON_ASSERT(!meta.MetaData().empty());

  std::map<ObjectID, std::shared_ptr<RemoteBlob>> remote_blobs;
  RETURN_ON_ERROR(
      GetRemoteBlobs(meta.buffer_set_->AllBufferIds(), remote_blobs));
  for (auto& item : remote_blobs) {
    RETURN_ON_ERROR(
        meta.buffer_set_->EmplaceBuffer(item.first, item.second->Buffer()));
  }

  // as we now has these buffers, using force local to let `Blob` aware
  // the buffer in `Construct()`.
  meta.ForceLocal();

  object = ObjectFactory::Create(meta.GetTypeName());
  if (object == nullptr) {
    object = std::unique_ptr<Object>(new Object());
  }
  object->Construct(meta);
  return Status::OK();
}

std::vector<std::shared_ptr<Object>> RPCClient::GetObjects(
    const std::vector<ObjectID>& ids) {
  std::vector<std::shared_ptr<Object>> objects(ids.size());
  std::vector<ObjectMeta> metas;
  if (!this->GetMetaData(ids, metas, true).ok()) {
    for (size_t index = 0; index < ids.size(); ++index) {
      objects[index] = nullptr;
    }
    return objects;
  }
  for (size_t index = 0; index < metas.size(); ++index) {
    if (metas[index].MetaData().empty()) {
      objects[index] = nullptr;
    } else {
      auto object = ObjectFactory::Create(metas[index].GetTypeName());
      if (object == nullptr) {
        object = std::unique_ptr<Object>(new Object());
      }
      object->Construct(metas[index]);
      objects[index] = std::shared_ptr<Object>(object.release());
    }
  }
  return objects;
}

std::vector<ObjectMeta> RPCClient::ListObjectMeta(std::string const& pattern,
                                                  const bool regex,
                                                  size_t const limit, bool) {
  std::unordered_map<ObjectID, json> meta_trees;
  VINEYARD_CHECK_OK(ListData(pattern, regex, limit, meta_trees));

  // construct object metadatas
  std::vector<ObjectMeta> objects;
  objects.reserve(meta_trees.size());
  for (auto const& kv : meta_trees) {
    ObjectMeta meta;
    meta.SetMetaData(this, kv.second);
    objects.emplace_back(meta);
  }
  return objects;
}

std::vector<std::shared_ptr<Object>> RPCClient::ListObjects(
    std::string const& pattern, const bool regex, size_t const limit) {
  std::unordered_map<ObjectID, json> meta_trees;
  VINEYARD_CHECK_OK(ListData(pattern, regex, limit, meta_trees));

  // construct objects
  std::vector<std::shared_ptr<Object>> objects;
  objects.reserve(meta_trees.size());
  for (auto const& kv : meta_trees) {
    ObjectMeta meta;
    meta.SetMetaData(this, kv.second);
    auto object = ObjectFactory::Create(meta.GetTypeName());
    if (object == nullptr) {
      object = std::unique_ptr<Object>(new Object());
    }
    object->Construct(meta);
    objects.emplace_back(std::shared_ptr<Object>(object.release()));
  }
  return objects;
}

namespace detail {

Status compress_and_send(std::shared_ptr<Compressor> const& compressor, int fd,
                         const char* buffer, const size_t buffer_size) {
  RETURN_ON_ERROR(compressor->Compress(buffer, buffer_size));
  void* chunk = nullptr;
  size_t chunk_size = 0;
  while (compressor->Pull(chunk, chunk_size).ok()) {
    if (chunk_size == 0) {
      continue;
    }
    RETURN_ON_ERROR(send_bytes(fd, &chunk_size, sizeof(size_t)));
    RETURN_ON_ERROR(send_bytes(fd, chunk, chunk_size));
  }
  return Status::OK();
}

Status recv_and_decompress(std::shared_ptr<Decompressor> const& decompressor,
                           int fd, char* buffer, const size_t buffer_size) {
  size_t decompressed_offset = 0;
  void* incoming_buffer = nullptr;
  size_t incoming_buffer_size = 0;
  while (true) {
    RETURN_ON_ERROR(
        decompressor->Buffer(incoming_buffer, incoming_buffer_size));
    size_t nbytes = 0;
    RETURN_ON_ERROR(recv_bytes(fd, &nbytes, sizeof(size_t)));
    assert(nbytes <= incoming_buffer_size);
    RETURN_ON_ERROR(recv_bytes(fd, incoming_buffer, nbytes));
    RETURN_ON_ERROR(decompressor->Decompress(nbytes));
    size_t chunk_size = 0;
    while (decompressor
               ->Pull(buffer + decompressed_offset,
                      buffer_size - decompressed_offset, chunk_size)
               .ok()) {
      decompressed_offset += chunk_size;
      if (decompressed_offset == buffer_size) {
        break;
      }
    }
    // the decompressor is expected to be "finished"
    while (true) {
      char data;
      size_t size = 0;
      auto s = decompressor->Pull(&data, 1, size);
      if (s.IsStreamDrained()) {
        break;
      }
      assert(s.ok() && size == 0);
    }
    if (decompressed_offset == buffer_size) {
      break;
    }
  }
  return Status::OK();
}

}  // namespace detail

bool RPCClient::IsFetchable(const ObjectMeta& meta) {
  auto instance_id = meta.meta_["instance_id"];
  if (instance_id.is_null()) {
    // it is a newly created metadata
    return true;
  }
  return remote_instance_id_ == instance_id.get<InstanceID>();
}

Status RPCClient::CreateRemoteBlob(
    std::shared_ptr<RemoteBlobWriter> const& buffer, ObjectMeta& meta) {
  ENSURE_CONNECTED(this);
  VINEYARD_ASSERT(buffer != nullptr, "Expects a non-null remote blob rewriter");
  std::shared_ptr<Compressor> compressor;
  if (compression_enabled()) {
    compressor = std::make_shared<Compressor>();
  }

  ObjectID id;
  Payload payload;
  int fd_sent = -1;

  std::string message_out;
  if (rdma_connected_) {
    WriteCreateRemoteBufferRequest(buffer->size(), false, true, message_out);
  } else {
    WriteCreateRemoteBufferRequest(buffer->size(), !!compressor, false,
                                   message_out);
  }
  RETURN_ON_ERROR(doWrite(message_out));

  // receive a confirm to continue
  {
    json message_in;
    RETURN_ON_ERROR(doRead(message_in));
    RETURN_ON_ERROR(ReadCreateBufferReply(message_in, id, payload, fd_sent));
  }

  // send the actual payload
  if (rdma_connected_) {
#ifdef VINEYARD_WITH_RDMA
    RegisterMemInfo info;
    info.address = reinterpret_cast<uint64_t>(buffer->data());
    info.size = buffer->size();
    RETURN_ON_ERROR(RegisterMem(info));
    RETURN_ON_ERROR(
        rdma_client_->Write(buffer->data(), buffer->size(),
                            reinterpret_cast<uint64_t>(payload.pointer),
                            remote_info_.rkey, info.mr_desc, nullptr));
    RETURN_ON_ERROR(rdma_client_->GetTXCompletion(-1, nullptr));
#endif
  } else {
    if (compressor && buffer->size() > 0) {
      RETURN_ON_ERROR(detail::compress_and_send(
          compressor, vineyard_conn_, buffer->data(), buffer->size()));
    } else if (buffer->size() > 0) {
      RETURN_ON_ERROR(
          send_bytes(vineyard_conn_, buffer->data(), buffer->size()));
    }
    json message_in;
    RETURN_ON_ERROR(doRead(message_in));
    RETURN_ON_ERROR(ReadCreateBufferReply(message_in, id, payload, fd_sent));
    RETURN_ON_ASSERT(
        static_cast<size_t>(payload.data_size) == buffer->size(),
        "The result blob size doesn't match with the requested size, " +
            std::to_string(payload.data_size) + " vs. " +
            std::to_string(buffer->size()));
  }
  // add the metadata to allow adding the returned remote blob as member
  // without an extra get operation.
  meta.SetId(id);
  meta.SetTypeName("vineyard::Blob");
  meta.SetNBytes(payload.data_size);
  meta.SetInstanceId(this->remote_instance_id_);
  meta.AddKeyValue("length", 0);
  meta.AddKeyValue("transient", true);

  return Status::OK();
}

Status RPCClient::CreateRemoteBlobs(
    std::vector<std::shared_ptr<RemoteBlobWriter>> const& buffers,
    std::vector<ObjectMeta>& metas) {
  ENSURE_CONNECTED(this);
  std::vector<size_t> sizes;
  for (auto const& buffer : buffers) {
    VINEYARD_ASSERT(buffer != nullptr,
                    "Expects a non-null remote blob rewriter");
    sizes.emplace_back(buffer->size());
  }
  std::shared_ptr<Compressor> compressor;
  if (compression_enabled()) {
    compressor = std::make_shared<Compressor>();
  }

  std::vector<ObjectID> ids;
  std::vector<Payload> payloads;
  std::vector<int> fds_sent;

  std::string message_out;
  if (rdma_connected_) {
    WriteCreateRemoteBuffersRequest(sizes, false, true, message_out);
  } else {
    WriteCreateRemoteBuffersRequest(sizes, !!compressor, false, message_out);
  }
  RETURN_ON_ERROR(doWrite(message_out));

  // receive a confirm to continue
  {
    json message_in;
    RETURN_ON_ERROR(doRead(message_in));
    RETURN_ON_ERROR(
        ReadCreateBuffersReply(message_in, ids, payloads, fds_sent));
  }

  // send the actual payload
  if (rdma_connected_) {
#ifdef VINEYARD_WITH_RDMA
    for (size_t i = 0; i < payloads.size(); ++i) {
      RegisterMemInfo info;
      info.address = reinterpret_cast<uint64_t>(buffers[i]->data());
      info.size = buffers[i]->size();
      RETURN_ON_ERROR(RegisterMem(info));
      RETURN_ON_ERROR(
          rdma_client_->Write(buffers[i]->data(), buffers[i]->size(),
                              reinterpret_cast<uint64_t>(payloads[i].pointer),
                              remote_info_.rkey, info.mr_desc, nullptr));
      RETURN_ON_ERROR(rdma_client_->GetTXCompletion(-1, nullptr));
    }
#endif
  } else {
    for (auto const& buffer : buffers) {
      if (compressor && buffer->size() > 0) {
        RETURN_ON_ERROR(detail::compress_and_send(
            compressor, vineyard_conn_, buffer->data(), buffer->size()));
      } else if (buffer->size() > 0) {
        RETURN_ON_ERROR(
            send_bytes(vineyard_conn_, buffer->data(), buffer->size()));
      }
    }

    json message_in;
    RETURN_ON_ERROR(doRead(message_in));
    ids.clear();
    payloads.clear();
    fds_sent.clear();
    RETURN_ON_ERROR(
        ReadCreateBuffersReply(message_in, ids, payloads, fds_sent));
    RETURN_ON_ASSERT(
        payloads.size() == buffers.size(),
        "The result size doesn't match with the requested sizes: " +
            std::to_string(payloads.size()) + " vs. " +
            std::to_string(buffers.size()));
    for (size_t i = 0; i < payloads.size(); ++i) {
      RETURN_ON_ASSERT(
          static_cast<size_t>(payloads[i].data_size) == buffers[i]->size(),
          "The result blob size doesn't match with the requested size: " +
              std::to_string(payloads[i].data_size) + " vs. " +
              std::to_string(buffers[i]->size()));
    }
  }

  // add the metadata to allow adding the returned remote blob as member
  // without an extra get operation.
  for (size_t i = 0; i < payloads.size(); ++i) {
    ObjectMeta meta;
    meta.SetId(ids[i]);
    meta.SetTypeName("vineyard::Blob");
    meta.SetNBytes(payloads[i].data_size);
    meta.SetInstanceId(this->remote_instance_id_);
    meta.AddKeyValue("length", 0);
    meta.AddKeyValue("transient", true);
    metas.emplace_back(meta);
  }

  return Status::OK();
}

#ifdef VINEYARD_WITH_RDMA
Status RPCClient::RegisterMem(RegisterMemInfo& info) {
  RETURN_ON_ERROR(rdma_client_->RegisterMemory(info));

  std::lock_guard<std::mutex> guard(local_infos_mutex_);
  local_infos_.emplace_back(info);
  return Status::OK();
}
#endif

Status RPCClient::GetRemoteBlob(const ObjectID& id,
                                std::shared_ptr<RemoteBlob>& buffer) {
  return this->GetRemoteBlob(id, false, buffer);
}

Status RPCClient::GetRemoteBlob(const ObjectID& id, const bool unsafe,
                                std::shared_ptr<RemoteBlob>& buffer) {
  ENSURE_CONNECTED(this);
  std::shared_ptr<Decompressor> decompressor;
  if (compression_enabled()) {
    decompressor = std::make_shared<Decompressor>();
  }

  std::vector<Payload> payloads;
  std::vector<int> fd_sent;

  std::string message_out;
  if (rdma_connected_) {
    WriteGetRemoteBuffersRequest(std::set<ObjectID>{id}, unsafe, false, true,
                                 message_out);
  } else {
    WriteGetRemoteBuffersRequest(std::set<ObjectID>{id}, unsafe, !!decompressor,
                                 false, message_out);
  }
  RETURN_ON_ERROR(doWrite(message_out));
  json message_in;
  RETURN_ON_ERROR(doRead(message_in));
  RETURN_ON_ERROR(ReadGetBuffersReply(message_in, payloads, fd_sent));
  RETURN_ON_ASSERT(payloads.size() == 1, "Expects only one payload");

  buffer = std::shared_ptr<RemoteBlob>(new RemoteBlob(
      payloads[0].object_id, remote_instance_id_, payloads[0].data_size));
  // read the actual payload
  if (rdma_connected_) {
#ifdef VINEYARD_WITH_RDMA
    RegisterMemInfo info;
    info.address = reinterpret_cast<uint64_t>(buffer->mutable_data());
    info.size = buffer->size();
    RETURN_ON_ERROR(RegisterMem(info));
    RETURN_ON_ERROR(
        rdma_client_->Read(buffer->mutable_data(), payloads[0].data_size,
                           reinterpret_cast<uint64_t>(payloads[0].pointer),
                           remote_info_.rkey, info.mr_desc, nullptr));
    RETURN_ON_ERROR(rdma_client_->GetTXCompletion(-1, nullptr));
#endif
  } else {
    if (decompressor && payloads[0].data_size > 0) {
      RETURN_ON_ERROR(detail::recv_and_decompress(decompressor, vineyard_conn_,
                                                  buffer->mutable_data(),
                                                  payloads[0].data_size));
    } else if (payloads[0].data_size > 0) {
      RETURN_ON_ERROR(recv_bytes(vineyard_conn_, buffer->mutable_data(),
                                 payloads[0].data_size));
    }
  }
  return Status::OK();
}

Status RPCClient::GetRemoteBlobs(
    std::vector<ObjectID> const& ids,
    std::vector<std::shared_ptr<RemoteBlob>>& remote_blobs) {
  return this->GetRemoteBlobs(ids, false, remote_blobs);
}

Status RPCClient::GetRemoteBlobs(
    std::set<ObjectID> const& ids,
    std::map<ObjectID, std::shared_ptr<RemoteBlob>>& remote_blobs) {
  return this->GetRemoteBlobs(ids, false, remote_blobs);
}

Status RPCClient::GetRemoteBlobs(
    std::vector<ObjectID> const& ids, const bool unsafe,
    std::vector<std::shared_ptr<RemoteBlob>>& remote_blobs) {
  ENSURE_CONNECTED(this);
  std::shared_ptr<Decompressor> decompressor;
  if (compression_enabled()) {
    decompressor = std::make_shared<Decompressor>();
  }

  std::unordered_set<ObjectID> id_set(ids.begin(), ids.end());
  std::vector<Payload> payloads;
  std::vector<int> fd_sent;

  std::string message_out;
  if (rdma_connected_) {
    WriteGetRemoteBuffersRequest(id_set, unsafe, false, true, message_out);
  } else {
    WriteGetRemoteBuffersRequest(id_set, unsafe, !!decompressor, false,
                                 message_out);
  }
  RETURN_ON_ERROR(doWrite(message_out));
  json message_in;
  RETURN_ON_ERROR(doRead(message_in));
  RETURN_ON_ERROR(ReadGetBuffersReply(message_in, payloads, fd_sent));
  RETURN_ON_ASSERT(payloads.size() == id_set.size(),
                   "The result size doesn't match with the requested sizes: " +
                       std::to_string(payloads.size()) + " vs. " +
                       std::to_string(id_set.size()));

  std::unordered_map<ObjectID, std::shared_ptr<RemoteBlob>> id_payload_map;
  if (rdma_connected_) {
#ifdef VINEYARD_WITH_RDMA
    for (auto const& payload : payloads) {
      auto remote_blob = std::shared_ptr<RemoteBlob>(new RemoteBlob(
          payload.object_id, remote_instance_id_, payload.data_size));
      RegisterMemInfo info;
      info.address = reinterpret_cast<uint64_t>(remote_blob->mutable_data());
      info.size = remote_blob->size();
      RETURN_ON_ERROR(RegisterMem(info));
      RETURN_ON_ERROR(
          rdma_client_->Read(remote_blob->mutable_data(), payload.data_size,
                             reinterpret_cast<uint64_t>(payload.pointer),
                             remote_info_.rkey, info.mr_desc, nullptr));
      RETURN_ON_ERROR(rdma_client_->GetTXCompletion(-1, nullptr));
      id_payload_map[payload.object_id] = remote_blob;
    }
#endif
  } else {
    for (auto const& payload : payloads) {
      auto remote_blob = std::shared_ptr<RemoteBlob>(new RemoteBlob(
          payload.object_id, remote_instance_id_, payload.data_size));
      if (decompressor && payload.data_size > 0) {
        RETURN_ON_ERROR(detail::recv_and_decompress(
            decompressor, vineyard_conn_, remote_blob->mutable_data(),
            payload.data_size));
      } else {
        RETURN_ON_ERROR(recv_bytes(vineyard_conn_, remote_blob->mutable_data(),
                                   payload.data_size));
      }
      id_payload_map[payload.object_id] = remote_blob;
    }
  }
  // clear the result container
  remote_blobs.clear();

  for (auto const& id : ids) {
    auto it = id_payload_map.find(id);
    if (it == id_payload_map.end()) {
      remote_blobs.emplace_back(nullptr);
    } else {
      remote_blobs.emplace_back(it->second);
    }
  }
  return Status::OK();
}

Status RPCClient::GetRemoteBlobs(
    std::set<ObjectID> const& ids, const bool unsafe,
    std::map<ObjectID, std::shared_ptr<RemoteBlob>>& remote_blobs) {
  std::vector<ObjectID> ids_vec(ids.begin(), ids.end());
  std::vector<std::shared_ptr<RemoteBlob>> remote_blobs_vec;
  RETURN_ON_ERROR(GetRemoteBlobs(ids_vec, unsafe, remote_blobs_vec));
  for (size_t i = 0; i < ids_vec.size(); ++i) {
    remote_blobs[ids_vec[i]] = std::move(remote_blobs_vec[i]);
  }
  return Status::OK();
}

}  // namespace vineyard
