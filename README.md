# RTC
project/
├── cmd/                # Entry points for the application
│   └── main.go         # Main application file
├── config/             # Configuration-related code
│   └── config.go       # Load and manage .env and other settings
├── internal/           # Non-exported application code
│   ├── handlers/       # HTTP handlers
│   │   └── user.go     # User-related HTTP endpoints
│   ├── models/         # Database models
│   │   └── user.go     # User table structure and DB interactions
│   ├── middlewares/    # Middleware functions
│   │   └── auth.go     # Authentication middleware
│   └── validators/     # Validation-related code
│       └── validate.go # Custom validation logic
├── pkg/                # Reusable utilities or libraries
│   └── hashing/        # Password hashing utilities
│       └── bcrypt.go
└── go.mod              # Go module file
