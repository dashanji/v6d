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

#include <chrono>
#include <cstddef>
#include <cstring>
#include <iomanip>
#include <memory>
#include <string>
#include <thread>
#include <algorithm>
#include <vector>
#include <iostream>

#include "arrow/api.h"
#include "arrow/io/api.h"

#include "client/ds/blob.h"
#include "client/client.h"
#include "client/ds/object_meta.h"
#include "client/rpc_client.h"
#include "common/util/logging.h"

using namespace vineyard;  // NOLINT(build/namespaces)
using namespace std;

size_t parseDataSize(const std::string &sizeStr) {
    std::istringstream is(sizeStr);
    double size;
    is >> size;
    std::string unit;
    if (is.fail()) {
        throw std::invalid_argument("Invalid data size format");
    }
    // Skip leading whitespace
    is >> std::ws;
    // Read unit, if any
    if (is.peek() != std::istringstream::traits_type::eof()) {
        is >> unit;
    }
    if (unit.empty() || unit == "B" || unit == "b") {
        return static_cast<size_t>(size);
    } else if (unit == "K" || unit == "k" || unit == "KB" || unit == "kb" || unit == "KILOBYTE" || unit == "kilobyte") {
        return static_cast<size_t>(size * pow(1024, 1));
    } else if (unit == "M" || unit == "m" || unit == "MB" || unit == "mb" || unit == "MEGABYTE" || unit == "megabyte") {
        return static_cast<size_t>(size * pow(1024, 2));
    } else if (unit == "G" || unit == "g" || unit == "GB" || unit == "gb" || unit == "GIGABYTE" || unit == "gigabyte") {
        return static_cast<size_t>(size * pow(1024, 3));
    } else if (unit == "T" || unit == "t" || unit == "TB" || unit == "tb" || unit == "TERABYTE" || unit == "terabyte") {
        return static_cast<size_t>(size * pow(1024, 4));
    } else {
        throw std::invalid_argument("Unsupported data size unit");
    }
}

std::vector<std::string> generateRandomData(int requests_num, int data_size) {
  const char charset[] = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
  std::vector<std::string> value_list;
  for (int i = 0; i < requests_num; i++) {
    std::string value;
    value.resize(data_size);
    for (int i = 0; i < data_size; i++) {
      value[i] = charset[rand() % (sizeof(charset) - 1)];
    }
    value_list.push_back(value);
  }
  return value_list;
}

void PutBlob(
  std::string &ipc_socket,
  std::vector<std::string> &value_list,
  size_t &data_size,
  std::vector<double> &latencies
) {
  Client client;
  VINEYARD_CHECK_OK(client.Connect(ipc_socket));
  for (int i = 0; i < value_list.size(); i++) {
    auto start = std::chrono::high_resolution_clock::now();
    std::unique_ptr<BlobWriter> blob_writer;
    VINEYARD_CHECK_OK(client.CreateBlob(data_size, blob_writer));
    std::memcpy(blob_writer->data(), value_list[i].c_str(), data_size);
    blob_writer->Seal(client);
    auto end = std::chrono::high_resolution_clock::now();
    std::chrono::duration<double, std::micro> duration = end - start;
    latencies.push_back(duration.count());
  }
  client.Disconnect();
}

void GetBlob(
  std::string& ipc_socket,
  ObjectID& blob_id,
  int requests_num,
  std::vector<double> &latencies
) {
  Client client;
  VINEYARD_CHECK_OK(client.Connect(ipc_socket));
  for (int i = 0; i < requests_num; i++) {
    auto start = std::chrono::high_resolution_clock::now();
    std::shared_ptr<Blob> blob;
    VINEYARD_CHECK_OK(client.GetBlob(blob_id, blob));
    auto end = std::chrono::high_resolution_clock::now();
    std::chrono::duration<double, std::micro> duration = end - start;
    latencies.push_back(duration.count());
  }
  client.Disconnect();
}

void PutRemoteBlob(
  std::string& rpc_endpoint,
  std::string& rdma_endpoint,
  std::vector<std::shared_ptr<RemoteBlobWriter>>& remote_blob_writers,
  std::vector<double> &latencies
) {
  std::shared_ptr<RPCClient> rpc_client = std::make_shared<RPCClient>();
  VINEYARD_CHECK_OK(rpc_client->Connect(rpc_endpoint, "", "", rdma_endpoint));
  for (int i = 0; i < remote_blob_writers.size(); i++) {
    auto start = std::chrono::high_resolution_clock::now();
    ObjectMeta meta;
    VINEYARD_CHECK_OK(rpc_client->CreateRemoteBlob(remote_blob_writers[i], meta));
    auto end = std::chrono::high_resolution_clock::now();
    std::chrono::duration<double, std::micro> duration = end - start;
    latencies.push_back(duration.count());
  }
  rpc_client->Disconnect();
}

void GetRemoteBlob(
  std::string& rpc_endpoint,
  std::string& rdma_endpoint,
  ObjectID& id,
  int& requests_num,
  std::vector<double> &latencies
) {
  std::shared_ptr<RPCClient> rpc_client = std::make_shared<RPCClient>();
  VINEYARD_CHECK_OK(rpc_client->Connect(rpc_endpoint, "", "", rdma_endpoint));
  for (int i = 0; i < requests_num; i++) {
    auto start = std::chrono::high_resolution_clock::now();
    std::shared_ptr<RemoteBlob> remote_blob;
    VINEYARD_CHECK_OK(rpc_client->GetRemoteBlob(id, remote_blob));
    auto end = std::chrono::high_resolution_clock::now();
    std::chrono::duration<double, std::micro> duration = end - start;
    latencies.push_back(duration.count());
  }
  rpc_client->Disconnect();
}

void TestPutBlob(
  int clients_num,
  std::string ipc_socket,
  std::vector<std::vector<std::string>> &value_lists,
  size_t data_size,
  std::vector<double>& put_blob_latencies
){
  std::vector<std::thread> threads;
  std::vector<std::vector<double>> local_latencies(clients_num);
  for (int i = 0; i < clients_num; i++) {
    threads.push_back(std::thread(PutBlob, std::ref(ipc_socket),
                                    std::ref(value_lists[i]),
                                    std::ref(data_size),
                                    std::ref(local_latencies[i])));
  }
  for (int i = 0; i < clients_num; i++) {
    threads[i].join();
  }
  for (auto &latencies : local_latencies) {
    put_blob_latencies.insert(put_blob_latencies.end(), latencies.begin(), latencies.end());
  }
}

void TestGetBlob(
  int clients_num,
  std::string ipc_socket,
  std::vector<int>& requests_num,
  size_t data_size,
  std::vector<double>& get_blob_latencies
){
  Client client;
  VINEYARD_CHECK_OK(client.Connect(ipc_socket));
  std::unique_ptr<BlobWriter> blob_writer;
  VINEYARD_CHECK_OK(client.CreateBlob(data_size, blob_writer));
  std::vector<std::string> values = generateRandomData(1, data_size);
  std::memcpy(blob_writer->data(), values[0].c_str(), data_size);
  std::shared_ptr<Object> blob = blob_writer->Seal(client);
  ObjectID blob_id = blob->id();
  std::vector<std::thread> threads;
  std::vector<std::vector<double>> local_latencies(clients_num);
  for (int i = 0; i < clients_num; i++) {
    threads.push_back(std::thread(GetBlob, std::ref(ipc_socket),
                                    std::ref(blob_id),
                                    std::ref(requests_num[i]),
                                    std::ref(local_latencies[i])
                                  ));
  }
  for (int i = 0; i < clients_num; i++) {
    threads[i].join();
  }
  client.Disconnect();
  for (auto &latencies : local_latencies) {
    get_blob_latencies.insert(get_blob_latencies.end(), latencies.begin(), latencies.end());
  }
}

void TestPutRemoteBlob(
  int clients_num,
  std::string rpc_endpoint,
  std::string rdma_endpoint,
  std::vector<std::vector<std::shared_ptr<RemoteBlobWriter>>> &remote_blob_writers,
  std::vector<double>& put_remote_blob_latencies
){
  std::vector<std::thread> threads;
  std::vector<std::vector<double>> local_latencies(clients_num);
  for (int i = 0; i < clients_num; i++) {
    threads.push_back(std::thread(PutRemoteBlob, std::ref(rpc_endpoint), std::ref(rdma_endpoint),
                                    std::ref(remote_blob_writers[i]),
                                    std::ref(local_latencies[i])));
  }
  for (int i = 0; i < clients_num; i++) {
    threads[i].join();
  }
  for (auto &latencies : local_latencies) {
    put_remote_blob_latencies.insert(put_remote_blob_latencies.end(), latencies.begin(), latencies.end());
  }
}

void TestGetRemoteBlob(
  int clients_num,
  std::string rpc_endpoint,
  std::string rdma_endpoint,
  std::vector<int> requests_num,
  size_t data_size,
  std::vector<double>& get_remote_blob_latencies
){
  std::shared_ptr<RPCClient> rpc_client= std::make_shared<RPCClient>();
  VINEYARD_CHECK_OK(rpc_client->Connect(rpc_endpoint, "", "", rdma_endpoint));
  std::shared_ptr<RemoteBlobWriter> remote_blob_writer(new RemoteBlobWriter(data_size));
  std::vector<std::string> values = generateRandomData(1, data_size);
  std::memcpy(remote_blob_writer->data(), values[0].c_str(), data_size);
  ObjectMeta meta;
  VINEYARD_CHECK_OK(rpc_client->CreateRemoteBlob(remote_blob_writer, meta));
  ObjectID blob_id = meta.GetId();
  std::vector<std::thread> threads;
  std::vector<std::vector<double>> local_latencies(clients_num);
  for (int i = 0; i < clients_num; i++) {
    threads.push_back(std::thread(GetRemoteBlob, std::ref(rpc_endpoint), std::ref(rdma_endpoint),
                                    std::ref(blob_id),
                                    std::ref(requests_num[i]),
                                    std::ref(local_latencies[i])
                                  ));
  }
  for (int i = 0; i < clients_num; i++) {
    threads[i].join();
  }
  rpc_client->Disconnect();
  for (auto &latencies : local_latencies) {
    get_remote_blob_latencies.insert(get_remote_blob_latencies.end(), latencies.begin(), latencies.end());
  }
}

void printStats(const std::string &op_name, int requests_num, int clients_num, size_t data_size, const std::vector<double> &latencies) {
    double total_time = std::accumulate(latencies.begin(), latencies.end(), 0.0);
    double min_time = *std::min_element(latencies.begin(), latencies.end());
    double max_time = *std::max_element(latencies.begin(), latencies.end());
    double avg_time = total_time / latencies.size();
    double throughput = requests_num / (total_time / 1e6);
    std::vector<double> sorted_latencies = latencies;
    std::sort(sorted_latencies.begin(), sorted_latencies.end());
    auto percentile = [&](int p) {
        return sorted_latencies[p * latencies.size() / 100];
    };
    std::cout << "====== " << op_name << " ======" << std::endl;
    std::cout << std::fixed << std::setprecision(2);
    std::cout << "  " << requests_num << " requests completed in " << (total_time / 1e6) << " seconds" << std::endl;
    std::cout << "  " << clients_num << " parallel clients" << std::endl;
    std::cout << "  " << data_size << " bytes payload" << std::endl;
    std::cout << "  min / avg / max latencies: " << min_time << " / " << avg_time << " / " << max_time << " μs" << std::endl;
    std::cout << "  throughput: " << throughput << " requests per second" << std::endl;
    std::cout << "  latencies percentiles:" << std::endl;
    std::cout << "    p50: " << percentile(50) << " μs" << std::endl;
    std::cout << "    p95: " << percentile(95) << " μs" << std::endl;
    std::cout << "    p99: " << percentile(99) << " μs" << std::endl;
}

int main(int argc, const char** argv) {
    if (argc < 7) {
        LOG(ERROR) << "usage: " << argv[0] << " <ipc_socket>"
                   << " <rpc_endpoint>"
                   << " <rdma_endpoint>"
                   << " <clients_num>"
                   << " <data_size>"
                   << " <requests_num>";
        return -1;
    }
    std::string ipc_socket = std::string(argv[1]);
    std::string rpc_endpoint = std::string(argv[2]);
    std::string rdma_endpoint = std::string(argv[3]);
    int clients_num = std::stoi(argv[4]);
    std::string data_size_str = std::string(argv[5]);
    int requests_num = std::stoi(argv[6]);
    int requests_num_per_client = requests_num / clients_num;
    int left_requests_num = requests_num % clients_num;
    size_t data_size = parseDataSize(data_size_str);
    std::vector<std::vector<std::string>> value_lists;
    std::vector<std::vector<std::shared_ptr<RemoteBlobWriter>>> remote_blob_writers;
    std::vector<int> requests_num_list;
    
    // Initialize lists for each client
    for (int i = 0; i < clients_num; i++) {
        if (i == clients_num - 1) {
            value_lists.push_back(generateRandomData(requests_num_per_client + left_requests_num, data_size));
            requests_num_list.push_back(requests_num_per_client + left_requests_num);
            remote_blob_writers.push_back(std::vector<std::shared_ptr<RemoteBlobWriter>>(requests_num_per_client + left_requests_num));
        } else {
            value_lists.push_back(generateRandomData(requests_num_per_client, data_size));
            requests_num_list.push_back(requests_num_per_client);
            remote_blob_writers.push_back(std::vector<std::shared_ptr<RemoteBlobWriter>>(requests_num_per_client));
        }
    }
    for (auto &remote_blob_writer: remote_blob_writers) {
      for (auto &writer: remote_blob_writer) {
        writer = std::make_shared<RemoteBlobWriter>(data_size);
      }
    }
  
    std::vector<double> put_blob_latencies;
    std::vector<double> get_blob_latencies;
    std::vector<double> put_remote_blob_latencies;
    std::vector<double> get_remote_blob_latencies;
    try {
        // Test PutBlob
        auto start = std::chrono::high_resolution_clock::now();
        TestPutBlob(clients_num, ipc_socket, value_lists, data_size, put_blob_latencies);
        auto end = std::chrono::high_resolution_clock::now();
        std::chrono::duration<double> duration = end - start;
        printStats("PutBlob", requests_num, clients_num, data_size, put_blob_latencies);

        // Test GetBlob
        start = std::chrono::high_resolution_clock::now();
        TestGetBlob(clients_num, ipc_socket, requests_num_list, data_size, get_blob_latencies);
        end = std::chrono::high_resolution_clock::now();
        duration = end - start;
        printStats("GetBlob", requests_num, clients_num, data_size, get_blob_latencies);

        // Test PutRemoteBlob
        start = std::chrono::high_resolution_clock::now();
        TestPutRemoteBlob(clients_num, rpc_endpoint, rdma_endpoint, remote_blob_writers, put_remote_blob_latencies);
        end = std::chrono::high_resolution_clock::now();
        duration = end - start;
        printStats("PutRemoteBlob", requests_num, clients_num, data_size, put_remote_blob_latencies);

        // Test GetRemoteBlob
        start = std::chrono::high_resolution_clock::now();
        TestGetRemoteBlob(clients_num, rpc_endpoint, rdma_endpoint, requests_num_list, data_size, get_remote_blob_latencies);
        end = std::chrono::high_resolution_clock::now();
        duration = end - start;
        printStats("GetRemoteBlob", requests_num, clients_num, data_size, get_remote_blob_latencies);

    } catch (std::exception &e) {
        LOG(ERROR) << "Caught exception: " << e.what();
        return -1;
    }
    LOG(INFO) << "Passed benchmark suite test.";
    return 0;
}