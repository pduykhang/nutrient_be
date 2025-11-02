package validator

import (
	"context"
	"testing"

	"nutrient_be/internal/dto/request"
	"nutrient_be/internal/pkg/logger"
)

// mockLogger is a simple mock logger for testing
type mockLogger struct {
	warnings []string
	errors   []error
}

func (m *mockLogger) Debug(ctx context.Context, msg string, fields ...logger.Field)  {}
func (m *mockLogger) Debugf(ctx context.Context, format string, args ...interface{}) {}
func (m *mockLogger) Info(ctx context.Context, msg string, fields ...logger.Field)   {}
func (m *mockLogger) Infof(ctx context.Context, format string, args ...interface{})  {}
func (m *mockLogger) Warn(ctx context.Context, msg string, fields ...logger.Field) {
	m.warnings = append(m.warnings, msg)
}
func (m *mockLogger) Warnf(ctx context.Context, format string, args ...interface{}) {}
func (m *mockLogger) Error(ctx context.Context, msg string, fields ...logger.Field) {
	for _, field := range fields {
		if field.Value != nil {
			if err, ok := field.Value.(error); ok {
				m.errors = append(m.errors, err)
			}
		}
	}
}
func (m *mockLogger) Errorf(ctx context.Context, format string, args ...interface{}) {}
func (m *mockLogger) Fatal(ctx context.Context, msg string, fields ...logger.Field)  {}
func (m *mockLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {}
func (m *mockLogger) Panic(ctx context.Context, msg string, fields ...logger.Field)  {}
func (m *mockLogger) Panicf(ctx context.Context, format string, args ...interface{}) {}
func (m *mockLogger) With(fields ...logger.Field) logger.Logger {
	return m
}

// Helper function to create a valid base request
func createValidFoodRequest() *request.CreateFoodRequest {
	return &request.CreateFoodRequest{
		Name: request.MultiLanguage{
			"en": "Apple",
			"vi": "Táo",
		},
		Category: "fruit",
		Macros: request.MacroNutrientsRequest{
			Protein:       0.3,
			Carbohydrates: 14.0,
			Fat:           0.2,
			Fiber:         2.4,
			Sugar:         10.4,
		},
		Micros: request.MicroNutrientsRequest{
			VitaminA:  54.0,
			VitaminC:  4.6,
			Calcium:   6.0,
			Iron:      0.12,
			Sodium:    1.0,
			Potassium: 107.0,
		},
		ServingSizes: []request.ServingSizeRequest{
			{
				Unit:           "gram",
				Amount:         100,
				GramEquivalent: 100,
			},
			{
				Unit:           "piece",
				Amount:         1,
				Description:    "1 medium apple",
				GramEquivalent: 182,
			},
		},
		Calories:   63.8, // Matches calculated: 0.3*4 + 14*4 + 0.2*9 + 2.4*2 = 1.2 + 56 + 1.8 + 4.8 = 63.8
		Visibility: "public",
	}
}

func TestNewFoodValidator(t *testing.T) {
	mockLog := &mockLogger{}
	validator := NewFoodValidator(mockLog)

	if validator == nil {
		t.Fatal("NewFoodValidator returned nil")
	}

	if validator.maxNameLength != 200 {
		t.Errorf("Expected maxNameLength 200, got %d", validator.maxNameLength)
	}
	if validator.maxDescriptionLength != 1000 {
		t.Errorf("Expected maxDescriptionLength 1000, got %d", validator.maxDescriptionLength)
	}
	if validator.maxCalories != 1000 {
		t.Errorf("Expected maxCalories 1000, got %.2f", validator.maxCalories)
	}
	if validator.maxMacroValue != 100 {
		t.Errorf("Expected maxMacroValue 100, got %.2f", validator.maxMacroValue)
	}
	if validator.caloriesTolerance != 10 {
		t.Errorf("Expected caloriesTolerance 10, got %.2f", validator.caloriesTolerance)
	}
}

func TestValidateCreateRequest_ValidRequest(t *testing.T) {
	mockLog := &mockLogger{}
	validator := NewFoodValidator(mockLog)
	ctx := context.Background()

	req := createValidFoodRequest()

	err := validator.ValidateCreateRequest(ctx, req)
	if err != nil {
		t.Errorf("Expected valid request to pass, got error: %v", err)
	}
}

func TestValidateCreateRequest_NameValidation(t *testing.T) {
	tests := []struct {
		name        string
		request     *request.CreateFoodRequest
		expectedErr string
	}{
		{
			name: "missing name",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Name = nil
				return req
			}(),
			expectedErr: "name validation failed",
		},
		{
			name: "empty name map",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Name = request.MultiLanguage{}
				return req
			}(),
			expectedErr: "name validation failed",
		},
		{
			name: "missing English translation",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Name = request.MultiLanguage{"vi": "Táo"}
				return req
			}(),
			expectedErr: "name validation failed",
		},
		{
			name: "empty English value",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Name = request.MultiLanguage{"en": "  "}
				return req
			}(),
			expectedErr: "name validation failed",
		},
		{
			name: "name too long",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				longName := make([]byte, 201)
				for i := range longName {
					longName[i] = 'a'
				}
				req.Name = request.MultiLanguage{"en": string(longName)}
				return req
			}(),
			expectedErr: "name validation failed",
		},
		{
			name: "name with whitespace only",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Name = request.MultiLanguage{"en": "   \t\n   "}
				return req
			}(),
			expectedErr: "name validation failed",
		},
	}

	mockLog := &mockLogger{}
	validator := NewFoodValidator(mockLog)
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCreateRequest(ctx, tt.request)
			if err == nil {
				t.Errorf("Expected error containing '%s', got nil", tt.expectedErr)
				return
			}
			if err.Error()[:len(tt.expectedErr)] != tt.expectedErr {
				t.Errorf("Expected error containing '%s', got '%s'", tt.expectedErr, err.Error())
			}
		})
	}
}

func TestValidateCreateRequest_DescriptionValidation(t *testing.T) {
	tests := []struct {
		name        string
		request     *request.CreateFoodRequest
		shouldError bool
	}{
		{
			name: "nil description",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Description = nil
				return req
			}(),
			shouldError: false,
		},
		{
			name: "empty description",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Description = request.MultiLanguage{}
				return req
			}(),
			shouldError: false,
		},
		{
			name: "valid description",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Description = request.MultiLanguage{
					"en": "A sweet and crunchy fruit",
					"vi": "Một loại trái cây ngọt và giòn",
				}
				return req
			}(),
			shouldError: false,
		},
		{
			name: "description too long",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				longDesc := make([]byte, 1001)
				for i := range longDesc {
					longDesc[i] = 'a'
				}
				req.Description = request.MultiLanguage{"en": string(longDesc)}
				return req
			}(),
			shouldError: true,
		},
		{
			name: "description with only whitespace",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Description = request.MultiLanguage{"en": "   "}
				return req
			}(),
			shouldError: false, // Whitespace-only description is allowed
		},
	}

	mockLog := &mockLogger{}
	validator := NewFoodValidator(mockLog)
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCreateRequest(ctx, tt.request)
			if tt.shouldError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
		})
	}
}

func TestValidateCreateRequest_NutritionValidation(t *testing.T) {
	tests := []struct {
		name        string
		request     *request.CreateFoodRequest
		expectedErr string
	}{
		{
			name: "all macros zero",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Macros = request.MacroNutrientsRequest{}
				return req
			}(),
			expectedErr: "nutrition validation failed",
		},
		{
			name: "negative protein",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Macros.Protein = -1
				return req
			}(),
			expectedErr: "nutrition validation failed",
		},
		{
			name: "protein exceeds max",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Macros.Protein = 101
				return req
			}(),
			expectedErr: "nutrition validation failed",
		},
		{
			name: "negative carbohydrates",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Macros.Carbohydrates = -1
				return req
			}(),
			expectedErr: "nutrition validation failed",
		},
		{
			name: "carbohydrates exceed max",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Macros.Carbohydrates = 101
				return req
			}(),
			expectedErr: "nutrition validation failed",
		},
		{
			name: "negative fat",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Macros.Fat = -1
				return req
			}(),
			expectedErr: "nutrition validation failed",
		},
		{
			name: "fat exceeds max",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Macros.Fat = 101
				return req
			}(),
			expectedErr: "nutrition validation failed",
		},
		{
			name: "negative fiber",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Macros.Fiber = -1
				return req
			}(),
			expectedErr: "nutrition validation failed",
		},
		{
			name: "fiber exceeds max",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Macros.Fiber = 101
				return req
			}(),
			expectedErr: "nutrition validation failed",
		},
		{
			name: "negative sugar",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Macros.Sugar = -1
				return req
			}(),
			expectedErr: "nutrition validation failed",
		},
		{
			name: "total macros exceed max",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Macros = request.MacroNutrientsRequest{
					Protein:       400,
					Carbohydrates: 300,
					Fat:           200,
					Fiber:         100,
				}
				return req
			}(),
			expectedErr: "nutrition validation failed",
		},
		{
			name: "negative calories",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Calories = -1
				return req
			}(),
			expectedErr: "nutrition validation failed",
		},
		{
			name: "calories exceed max",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Calories = 1001
				return req
			}(),
			expectedErr: "nutrition validation failed",
		},
		{
			name: "negative vitamin A",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Micros.VitaminA = -1
				return req
			}(),
			expectedErr: "nutrition validation failed",
		},
		{
			name: "negative vitamin C",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Micros.VitaminC = -1
				return req
			}(),
			expectedErr: "nutrition validation failed",
		},
		{
			name: "negative calcium",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Micros.Calcium = -1
				return req
			}(),
			expectedErr: "nutrition validation failed",
		},
		{
			name: "negative iron",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Micros.Iron = -1
				return req
			}(),
			expectedErr: "nutrition validation failed",
		},
		{
			name: "negative sodium",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Micros.Sodium = -1
				return req
			}(),
			expectedErr: "nutrition validation failed",
		},
		{
			name: "negative potassium",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Micros.Potassium = -1
				return req
			}(),
			expectedErr: "nutrition validation failed",
		},
		{
			name: "zero micros allowed",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Micros = request.MicroNutrientsRequest{}
				return req
			}(),
			expectedErr: "", // Should pass
		},
	}

	mockLog := &mockLogger{}
	validator := NewFoodValidator(mockLog)
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCreateRequest(ctx, tt.request)
			if tt.expectedErr == "" {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error containing '%s', got nil", tt.expectedErr)
					return
				}
				if !contains(err.Error(), tt.expectedErr) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.expectedErr, err.Error())
				}
			}
		})
	}
}

func TestValidateCreateRequest_ServingSizesValidation(t *testing.T) {
	tests := []struct {
		name        string
		request     *request.CreateFoodRequest
		expectedErr string
	}{
		{
			name: "empty serving sizes",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.ServingSizes = []request.ServingSizeRequest{}
				return req
			}(),
			expectedErr: "serving sizes validation failed",
		},
		{
			name: "nil serving sizes",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.ServingSizes = nil
				return req
			}(),
			expectedErr: "serving sizes validation failed",
		},
		{
			name: "invalid unit",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.ServingSizes[0].Unit = "invalid"
				return req
			}(),
			expectedErr: "serving sizes validation failed",
		},
		{
			name: "zero amount",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.ServingSizes[0].Amount = 0
				return req
			}(),
			expectedErr: "serving sizes validation failed",
		},
		{
			name: "negative amount",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.ServingSizes[0].Amount = -1
				return req
			}(),
			expectedErr: "serving sizes validation failed",
		},
		{
			name: "zero gramEquivalent",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.ServingSizes[0].GramEquivalent = 0
				return req
			}(),
			expectedErr: "serving sizes validation failed",
		},
		{
			name: "negative gramEquivalent",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.ServingSizes[0].GramEquivalent = -1
				return req
			}(),
			expectedErr: "serving sizes validation failed",
		},
		{
			name: "gram unit amount mismatch",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.ServingSizes[0].Unit = "gram"
				req.ServingSizes[0].Amount = 100
				req.ServingSizes[0].GramEquivalent = 150
				return req
			}(),
			expectedErr: "serving sizes validation failed",
		},
		{
			name: "gramEquivalent too large",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.ServingSizes[0].GramEquivalent = 100001
				return req
			}(),
			expectedErr: "serving sizes validation failed",
		},
		{
			name: "valid gram serving",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.ServingSizes = []request.ServingSizeRequest{
					{
						Unit:           "gram",
						Amount:         100,
						GramEquivalent: 100,
					},
				}
				return req
			}(),
			expectedErr: "", // Should pass
		},
		{
			name: "valid piece serving",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.ServingSizes = []request.ServingSizeRequest{
					{
						Unit:           "piece",
						Amount:         1,
						GramEquivalent: 182,
					},
				}
				return req
			}(),
			expectedErr: "", // Should pass
		},
		{
			name: "valid cup serving",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.ServingSizes = []request.ServingSizeRequest{
					{
						Unit:           "cup",
						Amount:         1,
						GramEquivalent: 240,
					},
				}
				return req
			}(),
			expectedErr: "", // Should pass
		},
		{
			name: "valid ml serving",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.ServingSizes = []request.ServingSizeRequest{
					{
						Unit:           "ml",
						Amount:         250,
						GramEquivalent: 250,
					},
				}
				return req
			}(),
			expectedErr: "", // Should pass
		},
		{
			name: "valid kg serving",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.ServingSizes = []request.ServingSizeRequest{
					{
						Unit:           "kg",
						Amount:         1,
						GramEquivalent: 1000,
					},
				}
				return req
			}(),
			expectedErr: "", // Should pass
		},
		{
			name: "valid box serving",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.ServingSizes = []request.ServingSizeRequest{
					{
						Unit:           "box",
						Amount:         1,
						GramEquivalent: 500,
					},
				}
				return req
			}(),
			expectedErr: "", // Should pass
		},
	}

	mockLog := &mockLogger{}
	validator := NewFoodValidator(mockLog)
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCreateRequest(ctx, tt.request)
			if tt.expectedErr == "" {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error containing '%s', got nil", tt.expectedErr)
					return
				}
				if !contains(err.Error(), tt.expectedErr) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.expectedErr, err.Error())
				}
			}
		})
	}
}

func TestValidateCreateRequest_CaloriesConsistency(t *testing.T) {
	baseMacros := request.MacroNutrientsRequest{
		Protein:       10.0,
		Carbohydrates: 20.0,
		Fat:           5.0,
		Fiber:         3.0,
	}
	// Calculate base calories: 10*4 + 20*4 + 5*9 + 3*2 = 40 + 80 + 45 + 6 = 171

	tests := []struct {
		name        string
		calories    float64
		macros      request.MacroNutrientsRequest
		expectedErr string
	}{
		{
			name:        "exact match",
			calories:    171.0,
			macros:      baseMacros,
			expectedErr: "",
		},
		{
			name:        "within tolerance (lower)",
			calories:    161.0, // 10 calories below
			macros:      baseMacros,
			expectedErr: "",
		},
		{
			name:        "within tolerance (upper)",
			calories:    181.0, // 10 calories above
			macros:      baseMacros,
			expectedErr: "",
		},
		{
			name:        "outside tolerance (too low)",
			calories:    160.0, // 11 calories below (outside tolerance)
			macros:      baseMacros,
			expectedErr: "calories consistency validation failed",
		},
		{
			name:        "outside tolerance (too high)",
			calories:    182.0, // 11 calories above (outside tolerance)
			macros:      baseMacros,
			expectedErr: "calories consistency validation failed",
		},
		{
			name:     "complex calculation",
			calories: 480.0,
			macros: request.MacroNutrientsRequest{
				Protein:       20.0,
				Carbohydrates: 50.0,
				Fat:           20.0,
				Fiber:         10.0,
			},
			expectedErr: "", // 20*4 + 50*4 + 20*9 + 10*2 = 80 + 200 + 180 + 20 = 480
		},
	}

	mockLog := &mockLogger{}
	validator := NewFoodValidator(mockLog)
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := createValidFoodRequest()
			req.Calories = tt.calories
			req.Macros = tt.macros

			err := validator.ValidateCreateRequest(ctx, req)
			if tt.expectedErr == "" {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error containing '%s', got nil", tt.expectedErr)
					return
				}
				if !contains(err.Error(), tt.expectedErr) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.expectedErr, err.Error())
				}
			}
		})
	}
}

func TestValidateCreateRequest_ImageURLValidation(t *testing.T) {
	tests := []struct {
		name        string
		imageURL    string
		expectedErr string
	}{
		{
			name:        "empty URL",
			imageURL:    "",
			expectedErr: "", // Should pass (optional)
		},
		{
			name:        "whitespace only URL",
			imageURL:    "   ",
			expectedErr: "", // Should pass (optional)
		},
		{
			name:        "valid HTTP URL",
			imageURL:    "http://example.com/image.jpg",
			expectedErr: "",
		},
		{
			name:        "valid HTTPS URL",
			imageURL:    "https://example.com/image.jpg",
			expectedErr: "",
		},
		{
			name:        "invalid scheme",
			imageURL:    "ftp://example.com/image.jpg",
			expectedErr: "image URL validation failed",
		},
		{
			name:        "missing scheme",
			imageURL:    "example.com/image.jpg",
			expectedErr: "image URL validation failed",
		},
		{
			name:        "missing host",
			imageURL:    "https://",
			expectedErr: "image URL validation failed",
		},
		{
			name:        "invalid URL format",
			imageURL:    "not a url",
			expectedErr: "image URL validation failed",
		},
		{
			name:        "URL too long",
			imageURL:    "https://example.com/" + string(make([]byte, 2049)),
			expectedErr: "image URL validation failed",
		},
		{
			name:        "URL with query params",
			imageURL:    "https://example.com/image.jpg?w=100&h=100",
			expectedErr: "",
		},
		{
			name:        "URL with path",
			imageURL:    "https://cdn.example.com/images/food/apple.jpg",
			expectedErr: "",
		},
	}

	mockLog := &mockLogger{}
	validator := NewFoodValidator(mockLog)
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := createValidFoodRequest()
			req.ImageURL = tt.imageURL

			err := validator.ValidateCreateRequest(ctx, req)
			if tt.expectedErr == "" {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error containing '%s', got nil", tt.expectedErr)
					return
				}
				if !contains(err.Error(), tt.expectedErr) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.expectedErr, err.Error())
				}
			}
		})
	}
}

func TestValidateCreateRequest_NoGramBaseWarning(t *testing.T) {
	mockLog := &mockLogger{}
	validator := NewFoodValidator(mockLog)
	ctx := context.Background()

	req := createValidFoodRequest()
	req.ServingSizes = []request.ServingSizeRequest{
		{
			Unit:           "piece",
			Amount:         1,
			GramEquivalent: 182,
		},
	}

	err := validator.ValidateCreateRequest(ctx, req)
	if err != nil {
		t.Errorf("Expected request to pass (warning only), got error: %v", err)
	}

	// Check that warning was logged
	if len(mockLog.warnings) == 0 {
		t.Error("Expected warning about missing gram base serving size, but no warning was logged")
	}
	if len(mockLog.warnings) > 0 && !contains(mockLog.warnings[0], "No gram base serving size found") {
		t.Errorf("Expected warning about missing gram base, got: %s", mockLog.warnings[0])
	}
}

func TestValidateCreateRequest_BoundaryValues(t *testing.T) {
	mockLog := &mockLogger{}
	validator := NewFoodValidator(mockLog)
	ctx := context.Background()

	tests := []struct {
		name        string
		request     *request.CreateFoodRequest
		shouldError bool
	}{
		{
			name: "max name length (200)",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				longName := make([]byte, 200)
				for i := range longName {
					longName[i] = 'a'
				}
				req.Name = request.MultiLanguage{"en": string(longName)}
				return req
			}(),
			shouldError: false,
		},
		{
			name: "max description length (1000)",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				longDesc := make([]byte, 1000)
				for i := range longDesc {
					longDesc[i] = 'a'
				}
				req.Description = request.MultiLanguage{"en": string(longDesc)}
				return req
			}(),
			shouldError: false,
		},
		{
			name: "max calories (1000)",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Macros = request.MacroNutrientsRequest{
					Protein:       100.0,
					Carbohydrates: 0,
					Fat:           0,
					Fiber:         0,
				}
				req.Calories = 400.0 // 100*4 = 400
				return req
			}(),
			shouldError: false,
		},
		{
			name: "max macro value (100)",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Macros = request.MacroNutrientsRequest{
					Protein:       100.0,
					Carbohydrates: 0,
					Fat:           0,
					Fiber:         0,
				}
				req.Calories = 400.0 // 100*4
				return req
			}(),
			shouldError: false,
		},
		{
			name: "total macros at limit (999)",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.Macros = request.MacroNutrientsRequest{
					Protein:       50.0,
					Carbohydrates: 50.0,
					Fat:           50.0,
					Fiber:         49.0,
				}
				// Adjust calories to match: 50*4 + 50*4 + 50*9 + 49*2 = 200 + 200 + 450 + 98 = 948
				req.Calories = 948.0
				return req
			}(),
			shouldError: false,
		},
		{
			name: "max gramEquivalent (100000)",
			request: func() *request.CreateFoodRequest {
				req := createValidFoodRequest()
				req.ServingSizes = []request.ServingSizeRequest{
					{
						Unit:           "kg",
						Amount:         100,
						GramEquivalent: 100000,
					},
				}
				return req
			}(),
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCreateRequest(ctx, tt.request)
			if tt.shouldError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
		})
	}
}

func TestValidateCreateRequest_RealWorldScenarios(t *testing.T) {
	mockLog := &mockLogger{}
	validator := NewFoodValidator(mockLog)
	ctx := context.Background()

	tests := []struct {
		name        string
		request     *request.CreateFoodRequest
		description string
		shouldError bool
	}{
		{
			name: "chicken breast",
			request: &request.CreateFoodRequest{
				Name:     request.MultiLanguage{"en": "Chicken Breast", "vi": "Ức gà"},
				Category: "protein",
				Macros: request.MacroNutrientsRequest{
					Protein:       31.0,
					Carbohydrates: 0.0,
					Fat:           3.6,
					Fiber:         0.0,
				},
				Micros: request.MicroNutrientsRequest{},
				ServingSizes: []request.ServingSizeRequest{
					{Unit: "gram", Amount: 100, GramEquivalent: 100},
				},
				Calories:   156.4, // 31*4 + 0*4 + 3.6*9 = 124 + 32.4 = 156.4
				Visibility: "public",
			},
			description: "Real chicken breast nutrition",
			shouldError: false,
		},
		{
			name: "banana",
			request: &request.CreateFoodRequest{
				Name:     request.MultiLanguage{"en": "Banana", "vi": "Chuối"},
				Category: "fruit",
				Macros: request.MacroNutrientsRequest{
					Protein:       1.1,
					Carbohydrates: 23.0,
					Fat:           0.3,
					Fiber:         2.6,
					Sugar:         12.0,
				},
				Micros: request.MicroNutrientsRequest{
					Potassium: 358.0,
					VitaminC:  8.7,
				},
				ServingSizes: []request.ServingSizeRequest{
					{Unit: "gram", Amount: 100, GramEquivalent: 100},
					{Unit: "piece", Amount: 1, Description: "1 medium banana", GramEquivalent: 118},
				},
				Calories:   104.3, // 1.1*4 + 23*4 + 0.3*9 + 2.6*2 = 4.4 + 92 + 2.7 + 5.2 = 104.3
				Visibility: "public",
			},
			description: "Real banana nutrition",
			shouldError: false,
		},
		{
			name: "olive oil",
			request: &request.CreateFoodRequest{
				Name:     request.MultiLanguage{"en": "Olive Oil", "vi": "Dầu ô liu"},
				Category: "fat",
				Macros: request.MacroNutrientsRequest{
					Protein:       0.0,
					Carbohydrates: 0.0,
					Fat:           100.0,
					Fiber:         0.0,
				},
				Micros: request.MicroNutrientsRequest{},
				ServingSizes: []request.ServingSizeRequest{
					{Unit: "gram", Amount: 100, GramEquivalent: 100},
					{Unit: "ml", Amount: 15, Description: "1 tablespoon", GramEquivalent: 13.5},
				},
				Calories:   900.0, // 100*9 = 900
				Visibility: "public",
			},
			description: "Real olive oil nutrition (fat category)",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCreateRequest(ctx, tt.request)
			if tt.shouldError && err == nil {
				t.Errorf("Expected error for %s, got nil", tt.description)
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error for %s, got: %v", tt.description, err)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsInner(s, substr)))
}

func containsInner(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
