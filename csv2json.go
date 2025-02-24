package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"io"
)

// Структура для записи данных CSV -> JSON
type Record map[string]string

func main() {
	// Открываем CSV-файл
	file, err := os.Open("input.csv")
	if err != nil {
		log.Fatalf("Ошибка при открытии CSV: %v", err)
	}
	defer file.Close()

	// Получаем размер файла для прогресса
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("Ошибка получения информации о файле: %v", err)
	}
	fileSize := fileInfo.Size()

	// Создаём CSV-ридер
	reader := csv.NewReader(file)
	reader.LazyQuotes = true         // Игнорирует ошибки с кавычками
	reader.TrimLeadingSpace = true   // Убирает пробелы в начале
	reader.FieldsPerRecord = -1      // Разрешает разное число колонок

	// Читаем заголовки (первую строку CSV)
	headers, err := reader.Read()
	if err == io.EOF {
		log.Fatal("Ошибка: Файл CSV пустой")
	} else if err != nil {
		log.Fatalf("Ошибка чтения заголовков: %v", err)
	}

	// Открываем JSON-файл для записи
	outFile, err := os.Create("output.json")
	if err != nil {
		log.Fatalf("Ошибка создания JSON-файла: %v", err)
	}
	defer outFile.Close()

	// Записываем открывающий символ JSON-массива
	outFile.WriteString("[\n")

	firstRecord := true
	for {
		// Читаем строку
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("Ошибка чтения строки CSV: %v", err)
		}

		// Создаём JSON-объект
		item := make(Record)
		for i, value := range record {
			if i < len(headers) {
				item[headers[i]] = value
			}
		}

		// Если это не первая запись, добавляем запятую перед новой строкой
		if !firstRecord {
			outFile.WriteString(",\n")
		}
		firstRecord = false

		// Кодируем в JSON и записываем в файл
		jsonData, err := json.MarshalIndent(item, "  ", "  ")
		if err != nil {
			log.Fatalf("Ошибка кодирования JSON: %v", err)
		}
		outFile.Write(jsonData)

		// Прогресс
		currentPosition, _ := file.Seek(0, io.SeekCurrent)
		progress := float64(currentPosition) / float64(fileSize) * 100
		fmt.Printf("\rПрогресс: %.2f%%", progress)
	}

	// Закрываем JSON-массив
	outFile.WriteString("\n]\n")

	fmt.Println("\n✅ Конвертация завершена. Данные записаны в output.json")
}
