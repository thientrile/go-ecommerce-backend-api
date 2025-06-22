package setting	
// Config chứa cấu hình cho ứng dụng, bao gồm các thiết lập server, database, logger, Redis, Kafka, JWT và rate limiter.
type Config struct {
	Server  ServerSetting  `mapstructure:"server"`   // Cấu hình server.
	Mysql   MSQLSetting    `mapstructure:"mysql"`    // Cấu hình cơ sở dữ liệu MySQL.
	Logger  LoggerSetting  `mapstructure:"logger"`   // Cấu hình ghi log.
	Redis   RedisSetting   `mapstructure:"redis"`    // Cấu hình Redis.
	Kafka   KafkaSetting   `mapstructure:"kafka"`    // Cấu hình Kafka.
	JWT     JWTSetting     `mapstructure:"jwt"`      // Cấu hình xác thực JWT.
	Limiter LimiterSetting `mapstructure:"limiter"`  // Cấu hình giới hạn tốc độ (rate limiter).
}

// ServerSetting chứa cấu hình cho server.
type ServerSetting struct {
	Port    int    `mapstructure:"port"`    // Cổng mà server lắng nghe.
	Mode    string `mapstructure:"mode"`    // Chế độ server (ví dụ: "development", "production").
	Domain  string `mapstructure:"domain"`  // Tên miền của server.
	Version string `mapstructure:"version"` // Phiên bản API.
}

// RedisSetting chứa cấu hình cho Redis.
type RedisSetting struct {
	Host      string `mapstructure:"host"`       // Địa chỉ host của Redis.
	Port      int    `mapstructure:"port"`       // Cổng của Redis.
	Password  string `mapstructure:"password"`   // Mật khẩu xác thực Redis.
	Database  int    `mapstructure:"database"`   // Số thứ tự database Redis.
	Pool_size int    `mapstructure:"pool_size"`  // Kích thước pool kết nối Redis.
}

// MSQLSetting chứa cấu hình cho cơ sở dữ liệu MySQL.
type MSQLSetting struct {
	Host            string `mapstructure:"host"`             // Địa chỉ host của MySQL.
	Port            int    `mapstructure:"port"`             // Cổng của MySQL.
	Username        string `mapstructure:"username"`         // Tên đăng nhập MySQL.
	Password        string `mapstructure:"password"`         // Mật khẩu MySQL.
	Dbname          string `mapstructure:"dbname"`           // Tên cơ sở dữ liệu MySQL.
	MaxIdleConns    int    `mapstructure:"maxIdleConns"`     // Số kết nối nhàn rỗi tối đa trong pool.
	MaxOpenConns    int    `mapstructure:"maxOpenConns"`     // Số kết nối tối đa tới database.
	ConnMaxLifetime int    `mapstructure:"connMaxLifetime"`  // Thời gian sống tối đa của một kết nối (tính bằng giây).
	ConnMaxIdleTime int    `mapstructure:"connMaxIdleTime"`  // Thời gian nhàn rỗi tối đa của một kết nối (tính bằng giây).
}

// LoggerSetting chứa cấu hình ghi log.
type LoggerSetting struct {
	Log_Level     string `mapstructure:"log_level"`      // Mức độ log (ví dụ: "info", "debug", "error").
	File_log_name string `mapstructure:"file_log_name"`  // Tên file log.
	File_log_path string `mapstructure:"file_log_path"`  // Đường dẫn file log.
	Max_size      int    `mapstructure:"max_size"`       // Kích thước tối đa của file log trước khi xoay vòng (MB).
	Max_backups   int    `mapstructure:"max_backups"`    // Số lượng file log cũ được giữ lại tối đa.
	Max_age       int    `mapstructure:"max_age"`        // Số ngày giữ lại file log cũ tối đa.
	Compress      bool   `mapstructure:"compress"`       // Có nén file log sau khi xoay vòng hay không.
}

// KafkaSetting chứa cấu hình cho Kafka.
type KafkaSetting struct {
	Host  string            `mapstructure:"host"`   // Địa chỉ host của Kafka broker.
	Port  int               `mapstructure:"port"`   // Cổng của Kafka broker.
	Topic map[string]string `mapstructure:"topic"`  // Danh sách các topic.
}

// JWTSetting chứa cấu hình xác thực JWT.
type JWTSetting struct {
	TokenHourLifespan int    `mapstructure:"TOKEN_HOUR_LIFETIME"` // Thời gian sống của token (giờ).
	JwtExpiration     string `mapstructure:"JWT_EXPIRATION"`      // Thời gian hết hạn JWT (dạng chuỗi).
	ApiSecret         string `mapstructure:"API_SECRET"`          // Khóa bí mật ký JWT.
	Issuer            string `mapstructure:"ISSUER"`              // Nhà phát hành JWT.
}

// LimiterSetting chứa cấu hình giới hạn tốc độ (rate limiting).
type LimiterSetting struct {
	Store         int                    `mapstructure:"store" json:"store"`                       // Loại store cho rate limiter (ví dụ: in-memory, Redis).
	DefaultConfig map[string]interface{} `mapstructure:"default_config" json:"default_config"`     // Cấu hình mặc định cho rate limiter.
	Rules         map[string]RuleConfig  `mapstructure:"rules" json:"rules"`                       // Các rule giới hạn tốc độ tùy chỉnh.
	URLPath       struct {
		Public  []string `mapstructure:"public"`   // Danh sách các đường dẫn URL public.
		Private []string `mapstructure:"private"`  // Danh sách các đường dẫn URL private.
	} `mapstructure:"url_path"`
}

// RuleConfig chứa cấu hình cho một rule giới hạn tốc độ cụ thể.
type RuleConfig struct {
	Rate            string   `mapstructure:"rate" json:"rate"`                         // Tốc độ giới hạn (ví dụ: "10-S" là 10 request/giây).
	Description     string   `mapstructure:"description" json:"description"`           // Mô tả rule.
	Enabled         bool     `mapstructure:"enabled" json:"enabled"`                   // Có bật rule này hay không.
	BurstMultiplier int      `mapstructure:"burst_multiplier" json:"burst_multiplier"` // Hệ số burst cho rate limiting.
	StrictMode      bool     `mapstructure:"strict_mode" json:"strict_mode"`           // Nếu bật, mọi vi phạm nhỏ đều bị chặn ngay lập tức, không có ngoại lệ.
	IPWhitelist     []string `mapstructure:"ip_whitelist" json:"ip_whitelist"`         // Danh sách IP được phép.
	MaxFileSize     string   `mapstructure:"max_file_size" json:"max_file_size"`       // Kích thước file tối đa cho phép (dạng chuỗi, ví dụ: "10MB").
}
