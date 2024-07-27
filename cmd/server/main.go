package main

import (
	"PandoraFuclaudePlusHelper/cmd/server/wire"
	commonConfig "PandoraFuclaudePlusHelper/config"
	"PandoraFuclaudePlusHelper/pkg/log"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"math/big"
)

// @title           Nunu Example API
// @version         1.0.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:8000
// @securityDefinitions.apiKey Bearer
// @in header
// @name Authorization
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {

	commonConfig.InitConfig()

	logger := log.NewLog()

	app, cleanup, err := wire.NewWire(logger)
	adminPassword := commonConfig.GetConfig().AdminPassword
	if adminPassword == "" {
		var err error
		// 生成至少12位的密码
		adminPassword, err = GenerateSecurePassword(12)
		if err != nil {
			// 如果生成密码出错，则抛出异常
			panic(errors.New("生成密码出错"))
		}
		fmt.Println("Generated admin password:", adminPassword)
		commonConfig.GetConfig().AdminPassword = adminPassword
	} else if len(adminPassword) < 8 {
		panic(errors.New("未设置密码或密码长度小于8"))
	}

	defer cleanup()
	if err != nil {
		panic(err)
	}

	logger.Info("server start", zap.Any("host", fmt.Sprintf("http://%s:%d", commonConfig.GetConfig().HttpHost, commonConfig.GetConfig().HttpPort)))
	if err = app.Run(context.Background()); err != nil {
		panic(err)
	}
}

// GenerateSecurePassword 生成一个安全的密码
func GenerateSecurePassword(length int) (string, error) {
	// 可以在密码中使用的字符集
	lowerLetters := "abcdefghijklmnopqrstuvwxyz"
	upperLetters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"
	symbols := "!@#$%^&*()_-+=[]{}|;:',.<>/?`~"

	allChars := lowerLetters + upperLetters + digits + symbols
	var password []byte

	// 确保密码中至少有一个字符来自每个所需的集合
	requiredSets := []string{lowerLetters, upperLetters, digits, symbols}
	for _, set := range requiredSets {
		char, err := randomCharFromSet(set)
		if err != nil {
			return "", err
		}
		password = append(password, char)
	}

	// 填充剩余的密码长度
	remainingLength := length - len(requiredSets)
	for i := 0; i < remainingLength; i++ {
		char, err := randomCharFromSet(allChars)
		if err != nil {
			return "", err
		}
		password = append(password, char)
	}

	// 打乱密码字符顺序以增加随机性
	shuffledPassword, err := shuffleBytes(password)
	if err != nil {
		return "", err
	}

	return string(shuffledPassword), nil
}

func randomCharFromSet(set string) (byte, error) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(set))))
	if err != nil {
		return 0, err
	}
	return set[nBig.Int64()], nil
}

func shuffleBytes(slice []byte) ([]byte, error) {
	shuffled := make([]byte, len(slice))
	copy(shuffled, slice)

	for i := range shuffled {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(len(shuffled))))
		if err != nil {
			return nil, err
		}
		shuffled[i], shuffled[j.Int64()] = shuffled[j.Int64()], shuffled[i]
	}

	return shuffled, nil
}
