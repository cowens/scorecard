// Copyright 2021 Security Scorecard Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"strconv"

	"github.com/ossf/scorecard/cron/bq"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/gcsblob"
)

func main() {
	ctx := context.Background()
	bucket, err := blob.OpenBucket(ctx, "gs://ossf-scorecard-data")
	if err != nil {
		panic(err)
	}
	defer bucket.Close()

	latestShard = nil
	shardNums := GetShardNumFiles()
	for _, shardNum := range shardNums {
		shardPrefix = GetShardPrefix(shardNum.Key)
		if CompletedBQTransfer(shardPrefix) {
			continue
		}
		data, err := bucket.ReadAll(ctx, shardNum.Key)
		if err != nil {
			panic(err)
		}
		num, err := strconv.Atoi(data)
		if err != nil {
			panic(err)
		}
		if GetNumShards(shardPrefix) == num {
			latest = shardPrefix
		}
	}
	if latestShard != nil {
		err := bq.StartDataTransferJob()
	}
	WriteCompleted(latestShard)
}
