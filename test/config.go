package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Configs 定义了一个嵌套的JSON配置结构
type Configs map[string]interface{}

// LoadConfig 从配置文件中载入json字符串
func LoadConfig(path string) (Configs, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("load config conf failed: %w", err)
	}
	var allConfigs Configs
	err = json.Unmarshal(buf, &allConfigs)
	if err != nil {
		return nil, fmt.Errorf("decode config file failed: %w", err)
	}
	return allConfigs, nil
}

// LoadConfigFromEnv 从环境变量中读取文件名并加载配置文件
func LoadConfigFromEnv() (Configs, error) {
	path := os.Getenv("go_config_path")
	if path == "" {
		return nil, fmt.Errorf("environment variable CONFIG_FILE_PATH is not set")
	}
	return LoadConfig(path)
}

// PrintConfigs 打印配置信息
func PrintConfigs(configs Configs) {
	printConfigsRecursive(configs, 0)
}

func printConfigsRecursive(configs map[string]interface{}, depth int) {
	for key, value := range configs {
		fmt.Printf("%s%s: ", getIndentation(depth), key)
		switch v := value.(type) {
		case map[string]interface{}:
			fmt.Println()
			printConfigsRecursive(v, depth+1)
		default:
			fmt.Println(v)
		}
	}
}

func getIndentation(depth int) string {
	return "  " // 两个空格作为缩进
}

// func main() {
// 	confs, err := LoadConfigFromEnv()
// 	// confs, err := LoadConfig("para.json")
// 	if err != nil {
// 		log.Fatalf("failed to load config: %v", err)
// 	}

// 	// 强制转换才能使用多级键值。因为空接口不支持直接访问多级键值
// 	fmt.Printf("type is %T", confs["etcd"].(map[string]interface{})["host"])
// }
