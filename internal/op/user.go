package op

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fengyuanluo/go-backend-template/internal/conf"
	"github.com/fengyuanluo/go-backend-template/internal/db"
	"github.com/fengyuanluo/go-backend-template/internal/model"
	"github.com/fengyuanluo/go-backend-template/internal/utils/log"
	"gorm.io/gorm"
)

var (
	userCache        model.User
	userMu           sync.RWMutex
	initialCredsPath = "data/initial-user.json"
	randomPassword   = func() (string, error) {
		buf := make([]byte, 16)
		if _, err := rand.Read(buf); err != nil {
			return "", err
		}
		return hex.EncodeToString(buf), nil
	}
)

var (
	ErrInvalidOldPassword = errors.New("invalid old password")
	ErrInvalidUsername    = errors.New("invalid username")
)

func initialCredentials() (string, string, bool, error) {
	prefix := strings.ToUpper(conf.APP_NAME)
	user := strings.TrimSpace(os.Getenv(prefix + "_INIT_USER"))
	pass := os.Getenv(prefix + "_INIT_PASSWORD")
	if user != "" && pass != "" {
		return user, pass, false, nil
	}
	pass, err := randomPassword()
	if err != nil {
		return "", "", false, err
	}
	user = "admin"
	if err := os.MkdirAll(filepath.Dir(initialCredsPath), 0755); err != nil {
		return "", "", false, err
	}
	payload := fmt.Sprintf(`{"username":"%s","password":"%s"}`, user, pass)
	if err := os.WriteFile(initialCredsPath, []byte(payload), 0600); err != nil {
		return "", "", false, err
	}
	return user, pass, true, nil
}

func UserInit() error {
	var user model.User
	if err := db.GetDB().First(&user).Error; err == nil {
		userMu.Lock()
		userCache = user
		userMu.Unlock()
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	username, password, wrote, err := initialCredentials()
	if err != nil {
		return err
	}
	user = model.User{Username: username, Password: password}
	if err := user.HashPassword(); err != nil {
		return err
	}
	if err := db.GetDB().Create(&user).Error; err != nil {
		return err
	}
	userMu.Lock()
	userCache = user
	userMu.Unlock()
	if wrote {
		log.Infof("initial user created, credentials saved to %s", initialCredsPath)
	} else {
		log.Infof("initial user created from environment variables")
	}
	return nil
}

func UserChangePassword(oldPassword, newPassword string) error {
	userMu.Lock()
	defer userMu.Unlock()
	if err := userCache.ComparePassword(oldPassword); err != nil {
		return fmt.Errorf("incorrect old password: %w", ErrInvalidOldPassword)
	}

	userCache.Password = newPassword
	if err := userCache.HashPassword(); err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	if err := db.GetDB().Model(&userCache).Update("password", userCache.Password).Error; err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func UserChangeUsername(newUsername string) error {
	userMu.Lock()
	defer userMu.Unlock()
	if userCache.Username == newUsername {
		return fmt.Errorf("new username is the same as the old username: %w", ErrInvalidUsername)
	}
	userCache.Username = newUsername
	if err := db.GetDB().Model(&userCache).Update("username", userCache.Username).Error; err != nil {
		return fmt.Errorf("failed to update username: %w", err)
	}
	return nil
}

func UserVerify(username, password string) error {
	userMu.RLock()
	defer userMu.RUnlock()
	if username != userCache.Username {
		return fmt.Errorf("incorrect username")
	}
	if err := userCache.ComparePassword(password); err != nil {
		return fmt.Errorf("incorrect password")
	}
	return nil
}

func UserGet() model.User {
	userMu.RLock()
	defer userMu.RUnlock()
	return userCache
}
