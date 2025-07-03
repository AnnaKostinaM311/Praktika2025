package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	//"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

// HealthData структура для хранения медицинских данных
type HealthData struct {
	UID    string  `json:"uid"`
	Age    int     `json:"age"`
	Gender int     `json:"gender"`
	RDW    float64 `json:"rdw"`
	WBC    float64 `json:"wbc"`
	RBC    float64 `json:"rbc"`
	HGB    float64 `json:"hgb"`
	HCT    float64 `json:"hct"`
	MCV    float64 `json:"mcv"`
	MCH    float64 `json:"mch"`
	MCHC   float64 `json:"mchc"`
	PLT    float64 `json:"plt"`
	NEU    float64 `json:"neu"`
	EOS    float64 `json:"eos"`
	BAS    float64 `json:"bas"`
	LYM    float64 `json:"lym"`
	MON    float64 `json:"mon"`
	SOE    float64 `json:"soe"`
	CHOL   float64 `json:"chol"`
	GLU    float64 `json:"glu"`
}

func main() {
	// Загрузка конфигурации
	apiURL := "https://apiml.labhub.online/api/v1/predict/hba1c"
	authToken := getEnv("API_AUTH_TOKEN", "Bearer 0l62<EJi/zJx]a?")
	port := getEnv("PORT", "8080")

	// Настройка HTTP-роутера
	http.HandleFunc("/api/forward", func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем только GET-запросы
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Парсинг параметров запроса
		query := r.URL.Query()
		data := parseHealthData(query)

		// Отправка данных в целевое API
		response, err := sendHealthData(apiURL, authToken, data)
		if err != nil {
			log.Printf("Forwarding error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Возвращаем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	})

	// Запуск сервера
	log.Printf("Starting API gateway on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// parseHealthData преобразует query-параметры в структуру HealthData
func parseHealthData(query map[string][]string) HealthData {
	data := HealthData{
		UID:    "web-client",
		Age:    0,
		Gender: 0,
		RDW:    0,
		WBC:    0,
		RBC:    0,
		HGB:    0,
		HCT:    0,
		MCV:    0,
		MCH:    0,
		MCHC:   0,
		PLT:    0,
		NEU:    0,
		EOS:    0,
		BAS:    0,
		LYM:    0,
		MON:    0,
		SOE:    0,
		CHOL:   0,
		GLU:    0,
	}

	for key, values := range query {
		if len(values) == 0 {
			continue
		}
		value := values[0]

		switch key {
		case "uid":
			data.UID = value
		case "age":
			if age, err := strconv.Atoi(value); err == nil {
				data.Age = age
			}
		case "gender":
			if gender, err := strconv.Atoi(value); err == nil {
				data.Gender = gender
			}
		case "rdw":
			data.RDW = parseFloat(value)
		case "wbc":
			data.WBC = parseFloat(value)
		case "rbc":
			data.RBC = parseFloat(value)
		case "hgb":
			data.HGB = parseFloat(value)
		case "hct":
			data.HCT = parseFloat(value)
		case "mcv":
			data.MCV = parseFloat(value)
		case "mch":
			data.MCH = parseFloat(value)
		case "mchc":
			data.MCHC = parseFloat(value)
		case "plt":
			data.PLT = parseFloat(value)
		case "neu":
			data.NEU = parseFloat(value)
		case "eos":
			data.EOS = parseFloat(value)
		case "bas":
			data.BAS = parseFloat(value)
		case "lym":
			data.LYM = parseFloat(value)
		case "mon":
			data.MON = parseFloat(value)
		case "soe":
			data.SOE = parseFloat(value)
		case "chol":
			data.CHOL = parseFloat(value)
		case "glu":
			data.GLU = parseFloat(value)
		}
	}

	return data
}

// parseFloat безопасно преобразует строку в float64
func parseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

// sendHealthData отправляет данные в целевое API
func sendHealthData(apiURL, authToken string, data HealthData) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json marshal error: %v", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("request creation error: %v", err)
	}

	req.Header.Set("Authorization", authToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response read error: %v", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("api error: %s", string(body))
	}

	return body, nil
}

// Вспомогательные функции для работы с окружением
func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is required", key)
	}
	return value
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
