package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

// Тип для записи CSV-данных
type Record map[string]string

func main() {
	// Открываем CSV-файл
	file, err := os.Open("input.csv")
	if err != nil {
		log.Fatalf("Ошибка при открытии CSV: %v", err)
	}
	defer file.Close()

	// Получаем размер файла (для прогресса)
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("Ошибка получения информации о файле: %v", err)
	}
	fileSize := fileInfo.Size()

	// Создаём CSV-ридер
	reader := csv.NewReader(bufio.NewReader(file))
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1

	// Читаем заголовки
	headers, err := reader.Read()
	if err == io.EOF {
		log.Fatal("Ошибка: CSV-файл пустой")
	} else if err != nil {
		log.Fatalf("Ошибка чтения заголовков: %v", err)
	}

	// Открываем JSON-файл с буферизированной записью
	outFile, err := os.Create("output.json")
	if err != nil {
		log.Fatalf("Ошибка создания JSON-файла: %v", err)
	}
	defer outFile.Close()
	jsonWriter := bufio.NewWriter(outFile)

	// Записываем начало JSON-массива
	jsonWriter.WriteString("[\n")

	var processedBytes int64
	var processedLines int64
	firstRecord := true

	// Используем буферизированный вывод прогресса
	progressWriter := bufio.NewWriter(os.Stdout)

	// Читаем и обрабатываем CSV построчно
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("Ошибка чтения CSV: %v", err)
		}

		// Создаём JSON-объект
		item := make(Record)
		for i, value := range record {
			if i < len(headers) {
				item[headers[i]] = value
			}
		}

		// Сериализуем JSON в строку (без лишних `\n`)
		jsonData, err := json.Marshal(item)
		if err != nil {
			log.Fatalf("Ошибка кодирования JSON: %v", err)
		}

		// Добавляем запятую между объектами (кроме первого)
		if !firstRecord {
			jsonWriter.WriteString(",\n")
		}
		firstRecord = false

		// Записываем JSON-объект в файл
		jsonWriter.Write(jsonData)

		// Обновляем статистику
		processedLines++
		processedBytes += int64(len(record))

		// Выводим прогресс каждые 10 000 строк
		if processedLines%10000 == 0 || processedBytes >= fileSize {
			progress := float64(processedBytes) / float64(fileSize) * 100
			fmt.Fprintf(progressWriter, "\rПрогресс: %.2f%% (%d строк обработано)", progress, processedLines)
			progressWriter.Flush() // Принудительный вывод в консоль
		}
	}

	// Завершаем JSON-массив
	jsonWriter.WriteString("\n]\n")
	jsonWriter.Flush() // Сбрасываем буфер в файл

	fmt.Println("\n✅ Конвертация завершена! Данные записаны в output.json")
}
