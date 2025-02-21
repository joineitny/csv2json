package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Структура для записи данных CSV -> JSON
type Record map[string]string

func main() {
	// Открываем CSV файл
	file, err := os.Open("input.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Получаем размер файла для расчета прогресса
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fileSize := fileInfo.Size()

	// Создаем CSV-ридер
	reader := csv.NewReader(file)

	// Читаем заголовки (первую строку CSV)
	headers, err := reader.Read()
	if err != nil {
		log.Fatal(err)
	}

	// Открываем JSON файл для записи
	outFile, err := os.Create("output.json")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	// Создаем JSON-энкодер
	encoder := json.NewEncoder(outFile)

	// Записываем открывающий символ JSON-массива
	outFile.WriteString("[\n")

	firstRecord := true

	// Построчное чтение CSV
	for {
		// Читаем строку
		record, err := reader.Read()
		if err != nil {
			break // Достигнут конец файла
		}

		// Создаем JSON-объект
		item := make(Record)
		for i, value := range record {
			item[headers[i]] = value
		}

		// Если это не первая строка, добавляем запятую перед записью
		if !firstRecord {
			outFile.WriteString(",\n")
		}
		firstRecord = false

		// Кодируем в JSON и записываем в файл через encoder
		encoder.Encode(item)

		// Получаем текущую позицию в файле
		currentPosition, _ := file.Seek(0, os.SEEK_CUR)
		progress := float64(currentPosition) / float64(fileSize) * 100
		fmt.Printf("\rПрогресс: %.2f%%", progress)
	}

	// Закрываем JSON-массив
	outFile.WriteString("\n]\n")

	fmt.Println("\nКонвертация завершена. Данные записаны в output.json")
}
