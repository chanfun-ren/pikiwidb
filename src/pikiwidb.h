/*
 * Copyright (c) 2023-present, Qihoo, Inc.  All rights reserved.
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "cmd_table_manager.h"
#include "cmd_thread_pool.h"
#include "common.h"
#include "io_thread_pool.h"
#include "net/tcp_connection.h"

#define KPIKIWIDB_VERSION "4.0.0"

#ifdef BUILD_DEBUG
#  define KPIKIWIDB_BUILD_TYPE "DEBUG"
#else
#  define KPIKIWIDB_BUILD_TYPE "RELEASE"
#endif

namespace pikiwidb {
class PRaft;
}  // namespace pikiwidb

class PikiwiDB final {
 public:
  PikiwiDB() = default;
  ~PikiwiDB() = default;

  bool ParseArgs(int ac, char* av[]);
  const PString& GetConfigName() const { return cfg_file_; }

  bool Init();
  void Run();
  //  void Recycle();
  void Stop();

  void OnNewConnection(pikiwidb::TcpConnection* obj);

  bool IsBgSaving() const { return is_bgsaving_; }
  void StartBgsave();
  void FinishBgsave();
  int64_t GetLastSave() const;

  //  pikiwidb::CmdTableManager& GetCmdTableManager();
  uint32_t GetCmdID() { return ++cmd_id_; };

  void SubmitFast(const std::shared_ptr<pikiwidb::CmdThreadPoolTask>& runner) { cmd_threads_.SubmitFast(runner); }

  void PushWriteTask(const std::shared_ptr<pikiwidb::PClient>& client) { worker_threads_.PushWriteTask(client); }

 public:
  PString cfg_file_;
  uint16_t port_{0};
  PString log_level_;

  PString master_;
  uint16_t master_port_{0};

  static const uint32_t kRunidSize;

 private:
  pikiwidb::WorkIOThreadPool worker_threads_;
  pikiwidb::IOThreadPool slave_threads_;
  pikiwidb::CmdThreadPool cmd_threads_;
  //  pikiwidb::CmdTableManager cmd_table_manager_;

  uint32_t cmd_id_ = 0;

  std::atomic<bool> is_bgsaving_ = false;
  int64_t lastsave_ = 0;
};

extern std::unique_ptr<PikiwiDB> g_pikiwidb;
