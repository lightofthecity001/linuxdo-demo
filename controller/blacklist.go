package controller

import (
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	// BlacklistCache 存储黑名单ID的切片
	BlacklistCache []string
	// mutex 读写锁，保护 BlacklistCache 的并发安全
	mutex sync.RWMutex
)

func StartBlacklistSync() {
	// 立即执行一次，防止服务刚启动时是空的
	refreshBlacklist()

	// 创建一个定时器，每 10 分钟触发一次
	ticker := time.NewTicker(10 * time.Minute)

	// 开启一个后台协程，默默干活
	go func() {
		for range ticker.C {
			log.Println("[定时任务] 开始同步黑名单...")
			refreshBlacklist()
		}
	}()
}

func refreshBlacklist() {
	// 替换为你的 Git Raw 地址
	url := "https://raw.githubusercontent.com/lightofthecity001/linuxdo-demo/refs/heads/main/blacklist.txt"

	client := http.Client{Timeout: 5 * time.Second} // 设置超时，防止卡死
	resp, err := client.Get(url)
	if err != nil {
		log.Println("[错误] 同步黑名单失败:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("[错误] Git 返回状态码异常:", resp.StatusCode)
		return
	}

	body, _ := io.ReadAll(resp.Body)
	content := string(body)

	// 处理数据：去除换行，按逗号分割
	// 假设文件内容是: 1,2,3
	parts := strings.Split(content, ",")
	newBlacklist := make([]string, 0, len(parts))
	for _, id := range parts {
		cleanId := strings.TrimSpace(id)
		if cleanId != "" {
			newBlacklist = append(newBlacklist, cleanId)
		}
	}

	// 【关键步骤】加写锁，更新全局变量
	mutex.Lock()
	BlacklistCache = newBlacklist
	mutex.Unlock()

	log.Printf("[成功] 黑名单已更新，当前共 %d 个用户被封禁\n", len(BlacklistCache))
}

func IsUserBanned(userId string) bool {
	// 【关键步骤】加读锁，允许多个用户同时查，但不允许写的插队
	mutex.RLock()
	defer mutex.RUnlock()

	for _, id := range BlacklistCache {
		if id == userId {
			return true
		}
	}
	return false
}
