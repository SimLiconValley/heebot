package config

import (
	"encoding/json"
	"errors" // 에러 생성을 위해 추가
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Token             string   `json:"token"`
	Domain            string   `json:"domain"`
	Port              int      `json:"port"`
	TempDirectoryName string   `json:"temp_directory_name"`
	ExpirySeconds     int      `json:"expiry_seconds"`
	AdminPassword     string   `json:"admin_password"`
	AdminWhiteListIds []string `json:"admin_whitelist_ids"`
}

var (
	AppConfig         Config
	adminWhiteListMap = make(map[string]bool)
)

func LoadConfig() error {
	file, err := os.Open("config.json")
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&AppConfig)
	if err != nil {
		return err
	}

	// 봇 토큰이 없을때
	if AppConfig.Token == "discord_bot_token" || strings.TrimSpace(AppConfig.Token) == "" {
		return errors.New("디스코드 봇 토큰을 입력해주세요")
	}

	// 1. Domain에서 혹시 모를 http://, https:// 제거
	AppConfig.Domain = strings.TrimPrefix(AppConfig.Domain, "http://")
	AppConfig.Domain = strings.TrimPrefix(AppConfig.Domain, "https://")

	// 2. TempDirectoryName에서 순수 폴더명만 떼어내기 (예: "./temp/" -> "temp")
	pureDir := filepath.Base(AppConfig.TempDirectoryName)
	AppConfig.TempDirectoryName = strings.Trim(pureDir, "./ ")

	// 💡 불러온 슬라이스(배열) 데이터를 조회용 맵으로 변환
	// 기존 맵 초기화 후 데이터 삽입
	adminWhiteListMap = make(map[string]bool)
	for _, id := range AppConfig.AdminWhiteListIds {
		trimmedID := strings.TrimSpace(id)
		if trimmedID != "" {
			adminWhiteListMap[trimmedID] = true
		}
	}
	return nil
}

func IsAdminWhiteList(userID string) bool {
	return adminWhiteListMap[userID]
}
