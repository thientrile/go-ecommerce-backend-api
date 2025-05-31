package cli

import "go.uber.org/zap"


func main(){
	sugar := zap.NewExample().Sugar()
	sugar.Infof("hello name: %s, age: %d", "thientriel", 18)
	// logger

	logger := zap.NewExample()
	logger.Info("hello world", zap.String("name", "thientriel"), zap.Int("age", 18))
}