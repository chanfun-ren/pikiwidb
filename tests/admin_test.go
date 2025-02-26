/*
 * Copyright (c) 2023-present, Qihoo, Inc.  All rights reserved.
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

package pikiwidb_test

import (
	"context"
	"log"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redis/go-redis/v9"

	"github.com/OpenAtomFoundation/pikiwidb/tests/util"
)

var _ = Describe("Admin", Ordered, func() {
	var (
		ctx    = context.TODO()
		s      *util.Server
		client *redis.Client
	)

	// BeforeAll closures will run exactly once before any of the specs
	// within the Ordered container.
	BeforeAll(func() {
		config := util.GetConfPath(false, 0)

		s = util.StartServer(config, map[string]string{"port": strconv.Itoa(7777)}, true)
		Expect(s).NotTo(Equal(nil))
	})

	// AfterAll closures will run exactly once after the last spec has
	// finished running.
	AfterAll(func() {
		err := s.Close()
		if err != nil {
			log.Println("Close Server fail.", err.Error())
			return
		}
	})

	// When running each spec Ginkgo will first run the BeforeEach
	// closure and then the subject closure.Doing so ensures that
	// each spec has a pristine, correctly initialized, copy of the
	// shared variable.
	BeforeEach(func() {
		client = s.NewClient()
	})

	// nodes that run after the spec's subject(It).
	AfterEach(func() {
		err := client.Close()
		if err != nil {
			log.Println("Close client conn fail.", err.Error())
			return
		}
	})

	//TODO(dingxiaoshuai) Add more test cases.
	It("Cmd INFO", func() {
		log.Println("Cmd INFO Begin")
		Expect(client.Info(ctx).Val()).NotTo(Equal("FooBar"))
	})

	It("Cmd Shutdown", func() {
		Expect(client.Shutdown(ctx).Err()).NotTo(HaveOccurred())

		// PikiwiDB does not support the Ping command right now
		// wait for 5 seconds and then ping server
		// time.Sleep(5 * time.Second)
		// Expect(client.Ping(ctx).Err()).To(HaveOccurred())

		// restart server
		config := util.GetConfPath(false, 0)
		s = util.StartServer(config, map[string]string{"port": strconv.Itoa(7777)}, true)
		Expect(s).NotTo(Equal(nil))

		// PikiwiDB does not support the Ping command right now
		// wait for 5 seconds and then ping server
		// time.Sleep(5 * time.Second)
		// client = s.NewClient()
		// Expect(client.Ping(ctx).Err()).NotTo(HaveOccurred())
	})

	It("Cmd Select", func() {
		var outRangeNumber = 100

		r, e := client.Set(ctx, DefaultKey, DefaultValue, 0).Result()
		Expect(e).NotTo(HaveOccurred())
		Expect(r).To(Equal(OK))

		r, e = client.Get(ctx, DefaultKey).Result()
		Expect(e).NotTo(HaveOccurred())
		Expect(r).To(Equal(DefaultValue))

		rDo, eDo := client.Do(ctx, kCmdSelect, outRangeNumber).Result()
		Expect(eDo).To(MatchError(kInvalidIndex))

		r, e = client.Get(ctx, DefaultKey).Result()
		Expect(e).NotTo(HaveOccurred())
		Expect(r).To(Equal(DefaultValue))

		rDo, eDo = client.Do(ctx, kCmdSelect, 1).Result()
		Expect(eDo).NotTo(HaveOccurred())
		Expect(rDo).To(Equal(OK))

		r, e = client.Get(ctx, DefaultKey).Result()
		Expect(e).To(MatchError(redis.Nil))
		Expect(r).To(Equal(Nil))

		rDo, eDo = client.Do(ctx, kCmdSelect, 0).Result()
		Expect(eDo).NotTo(HaveOccurred())
		Expect(rDo).To(Equal(OK))

		rDel, eDel := client.Del(ctx, DefaultKey).Result()
		Expect(eDel).NotTo(HaveOccurred())
		Expect(rDel).To(Equal(int64(1)))
	})

	It("Cmd Config", func() {
		res := client.ConfigGet(ctx, "timeout")
		Expect(res.Err()).NotTo(HaveOccurred())
		Expect(res.Val()).To(Equal(map[string]string{"timeout": "0"}))

		res = client.ConfigGet(ctx, "daemonize")
		Expect(res.Err()).NotTo(HaveOccurred())
		Expect(res.Val()).To(Equal(map[string]string{"daemonize": "no"}))

		resSet := client.ConfigSet(ctx, "timeout", "60")
		Expect(resSet.Err()).NotTo(HaveOccurred())
		Expect(resSet.Val()).To(Equal("OK"))

		resSet = client.ConfigSet(ctx, "daemonize", "yes")
		Expect(resSet.Err()).To(MatchError("ERR Invalid Argument"))

		res = client.ConfigGet(ctx, "timeout")
		Expect(res.Err()).NotTo(HaveOccurred())
		Expect(res.Val()).To(Equal(map[string]string{"timeout": "60"}))

		res = client.ConfigGet(ctx, "time*")
		Expect(res.Err()).NotTo(HaveOccurred())
		Expect(res.Val()).To(Equal(map[string]string{"timeout": "60"}))
	})

	It("PING", func() {
		ping := client.Ping(ctx)
		Expect(ping.Err()).NotTo(HaveOccurred())
	})

	It("Cmd DBSize", func() {
		res := client.DBSize(ctx)
		Expect(res.Err()).NotTo(HaveOccurred())
		Expect(res.Val()).To(Equal(int64(0)))

		client.Set(ctx, "key1", "value1", 0)
		Expect(res.Err()).NotTo(HaveOccurred())
		Expect(client.DBSize(ctx).Val()).To(Equal(int64(1)))

		client.Set(ctx, "key2", "value2", 0)
		Expect(res.Err()).NotTo(HaveOccurred())
		Expect(client.DBSize(ctx).Val()).To(Equal(int64(2)))

		client.Del(ctx, "key1")
		Expect(res.Err()).NotTo(HaveOccurred())
		Expect(client.DBSize(ctx).Val()).To(Equal(int64(1)))

		client.Del(ctx, "key2")
		Expect(res.Err()).NotTo(HaveOccurred())
		Expect(client.DBSize(ctx).Val()).To(Equal(int64(0)))
	})

	It("Cmd Lastsave", func() {
		lastSave := client.LastSave(ctx)
		Expect(lastSave.Err()).NotTo(HaveOccurred())
		Expect(lastSave.Val()).To(Equal(int64(0)))

		bgSaveTime1 := time.Now().Unix()
		bgSave, err := client.BgSave(ctx).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(bgSave).To(ContainSubstring("Background saving started"))
		time.Sleep(1 * time.Second)
		bgSaveTime2 := time.Now().Unix()

		lastSave = client.LastSave(ctx)
		Expect(lastSave.Err()).NotTo(HaveOccurred())
		Expect(lastSave.Val()).To(BeNumerically(">=", bgSaveTime1))
		Expect(lastSave.Val()).To(BeNumerically("<=", bgSaveTime2))
	})

	It("Cmd Debug", func() {
		// TODO: enable test after implementing DebugObject
		// res := client.DebugObject(ctx, "timeout")
		// Expect(res.Err()).NotTo(HaveOccurred())
		// Expect(res.Val()).To(Equal(map[string]string{"timeout": "0"}))
	})
})
