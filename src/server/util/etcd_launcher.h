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

#ifndef SRC_SERVER_UTIL_ETCD_LAUNCHER_H_
#define SRC_SERVER_UTIL_ETCD_LAUNCHER_H_

#include <netdb.h>

#include <memory>
#include <set>
#include <string>
#include <vector>

#if defined(BUILD_VINEYARDD_ETCD)

#include "boost/process/child.hpp"  // IWYU pragma: keep
#include "etcd/Client.hpp"

#include "common/util/status.h"
#include "server/util/etcd_member.h"

namespace vineyard {

class EtcdLauncher {
 public:
  explicit EtcdLauncher(const json& etcd_spec, const uint32_t& rpc_socket_port,
                        const bool create_new_instance);
  ~EtcdLauncher();

  Status LaunchEtcdServer(std::unique_ptr<etcd::Client>& etcd_client,
                          std::string& sync_lock);

  // Check if the etcd server available, return True if succeed, otherwise
  // False.
  static bool probeEtcdServer(std::unique_ptr<etcd::Client>& etcd_client,
                              std::string const& key);

 private:
  Status handleEtcdFailure(std::unique_ptr<etcd::Client>& etcd_client,
                           const std::string& member_name,
                           const std::string& errMessage);

  Status parseEndpoint();

  std::string generateMemberName(
      const std::vector<std::string>& existing_members_name);

  const uint64_t GetMemberID() { return etcd_member_id_; }

  Status RemoveMember(std::unique_ptr<etcd::Client>& etcd_client,
                      const uint64_t& member_id) {
    return removeMember(etcd_client, member_id);
  }

  Status UpdateEndpoint(std::unique_ptr<etcd::Client>& etcd_client);

  Status initHostInfo();

  const json etcd_spec_;
  const uint32_t rpc_socket_port_;
  const bool create_new_instance_;
  std::string endpoint_host_;
  std::string etcd_data_dir_;
  uint32_t endpoint_port_;
  std::set<std::string> local_hostnames_;
  std::set<std::string> local_ip_addresses_;

  uint64_t etcd_member_id_;
  std::string etcd_endpoints_;

  std::unique_ptr<boost::process::child> etcd_proc_;

  friend class EtcdMetaService;
};

}  // namespace vineyard

#endif  // BUILD_VINEYARDD_ETCD

#endif  // SRC_SERVER_UTIL_ETCD_LAUNCHER_H_
