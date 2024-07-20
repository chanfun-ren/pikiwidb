// Copyright (c) 2015-present, Qihoo, Inc.  All rights reserved.
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <algorithm>

#include "pstd_util.h"

namespace pstd {

std::mt19937 gen;

void InitRandom() {
  std::random_device rd;
  gen.seed(rd());
}

int RandomInt(int max) { return RandomInt(0, max); }

int RandomInt(int min, int max) {
  std::uniform_int_distribution<> dist(min, max);
  return dist(gen);
}

double RandomDouble() {
  std::uniform_real_distribution<double> dis(0.0, 1.0);
  return dis(gen);
}

std::string NowTimeStr() {
  auto now = std::chrono::system_clock::now();
  auto time_t_now = std::chrono::system_clock::to_time_t(now);
  std::ostringstream oss;
  oss << std::put_time(std::localtime(&time_t_now), "%Y-%m-%d_%H-%M-%S");
  return oss.str();
}

}  // namespace pstd
