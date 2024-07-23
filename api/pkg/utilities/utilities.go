package utilities

import (
	"api/pkg/logger"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// encode - кодирует  в base64
func encode(bin []byte) []byte {
	e64 := base64.StdEncoding

	maxEncLen := e64.EncodedLen(len(bin))
	encBuf := make([]byte, maxEncLen)

	e64.Encode(encBuf, bin)
	return encBuf
}

// Exists - проверяет существует ли файл
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func format(enc []byte, mime string) string {
	switch mime {
	case "image/gif", "image/jpeg", "image/pjpeg", "image/png", "image/tiff":
		return fmt.Sprintf("data:%s;charset=utf-8;base64,%s", mime, enc)
	default:
	}

	return fmt.Sprintf("data:image/png;base64,%s", enc)
}

// executeExt - извлекает расширение файла
func executeExt(str string) string {
	switch strings.TrimSuffix(str[5:strings.Index(str, ",")], ";base64") {
	case "image/png":
		return "png"
	case "image/jpeg":
		return "jpg"
	default:
		return ""
	}
}

// trimBase64 - отрезает лишнее и возвращает только base64
func trimBase64(str string) string {
	base64str := str[strings.Index(str, ",")+1:]
	return base64str
}

// decodeBase64 - декодирует base64 и возвращает []byte
func decodeBase64(b64 string) []byte {
	b, _ := base64.StdEncoding.DecodeString(b64)
	return b
}

// Base64ToFile - конвертирует base64 в файл и возвращает путь к нему
func Base64ToFile(str string) (string, error) {
	//расширение
	ext := executeExt(str)
	//байты
	trimedBase64 := trimBase64(str)
	unbased := decodeBase64(trimedBase64)
	//Путь к текущей директории
	pwd, _ := os.Getwd()
	//Путь к папке с images (полный путь)
	pathDir := filepath.Join(pwd, "images")
	//Случайное имя файла с расширением
	fileName := fmt.Sprintf("%s.%s", GenerateRandomString(24), ext)

	//Создание директории
	{
		err := os.Mkdir(pathDir, os.ModePerm)
		if err != nil && !os.IsExist(err) {
			logger.Error("Ошибка при создании директории: %s", err.Error())
			return "", err
		}
	}

	//Создание файла
	file, err := os.Create(filepath.Join(pathDir, fileName))
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.Write(unbased)
	if err != nil {
		return "", err
	}

	return filepath.Join(pathDir, fileName), nil
}

// fromBuffer принимает набор байтов в буфере
// возвращает base64 строку.
func fromBuffer(buf bytes.Buffer) string {
	enc := encode(buf.Bytes())
	mime := http.DetectContentType(buf.Bytes())
	logger.Info("mime: %s", mime)

	return format(enc, mime)
}

// fileFromLocal - файл с локального хранилища
func fileFromLocal(fname string) (string, error) {
	var b bytes.Buffer

	fileExists, _ := exists(fname)
	if !fileExists {
		return "", fmt.Errorf("File does not exist\n")
	}

	file, err := os.Open(fname)
	if err != nil {
		return "", fmt.Errorf("Error opening file\n")
	}

	_, err = b.ReadFrom(file)
	if err != nil {
		return "", fmt.Errorf("Error reading file to buffer\n")
	}

	return fromBuffer(b), nil
}

// FileToBase64 - конвертирует файл в base64
func FileToBase64(filepath string) string {
	result, err := fileFromLocal(filepath)
	if err != nil {
		logger.Error("Ошибка при попытке кодирования файла в base64: %s", err.Error())
		result = ""
	}

	return result
}

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "   ")
	if err != nil {
		return in
	}
	return out.String()
}

// ToJSON - конвертирует объект в JSON строку
func ToJSON(object any) string {
	jsonByte, err := json.Marshal(object)
	if err != nil {
		logger.Error("Ошибка при получении JSON: ", err.Error())
	}
	n := len(jsonByte)             //Find the length of the byte array
	result := string(jsonByte[:n]) //convert to string

	return jsonPrettyPrint(result)
}

// Compare - сравнивает строку с зашифрованной строкой
func Compare(str, cryptStr string) bool {
	logger.Info("crypted password: %s", MD5(str))
	b := MD5(str) == cryptStr
	logger.Info("сравнение: %t", b)
	return MD5(str) == cryptStr
}

// StrToUint - Конвертирует строку в uint
func StrToUint(s string) uint {
	i, err := strconv.Atoi(s)
	if err != nil {
		logger.Error("%s", err.Error())
		return 0
	}
	return uint(i)
}

// RandInt - возвращает случайное число от min до max
func RandInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

// GenerateRandomString - генерирует случайный набор символов (англ алфавит, case uppercase + символ _ и цифры от 0 до 9)
func GenerateRandomString(leng int) string {
	alph := "Q o S 4 r T 0 8 D m 7 d Z V O P w u c f M 2 h a F i N y E j B K 3 U t C 9 I q Y _ l z" +
		" v 6 X p g W s A J e G 5 H 1 R x n L b k"
	arrayStr := strings.Split(alph, " ")
	rand.Seed(time.Now().UnixNano())
	result := ""
	if leng == 0 {
		return ""
	}

	for i := 0; i < leng; i++ {
		result += arrayStr[RandInt(0, len(arrayStr))]
	}

	return result
}

// ChangeEnvAttribute - заменяет значение в attr на prop в .env файле
func ChangeEnvAttribute(attr, prop string) error {
	// Путь к файлу .env
	envFile := ".env"

	// Чтение файла .env
	data, err := os.ReadFile(envFile)
	if err != nil {
		return err
	}

	attrFormat := fmt.Sprintf("%s=", attr)
	newAttrProp := fmt.Sprintf("%s=%s", attr, prop)

	// Разделение файла .env на строки
	lines := strings.Split(string(data), "\n")

	// Флаг для определения замены
	replaced := false

	// Итерация по строкам файла .env
	for i, line := range lines {
		// Проверка, содержит ли строка "attrFormat"
		if strings.HasPrefix(line, attrFormat) {
			// Замена значения attr
			lines[i] = newAttrProp
			replaced = true
			break
		}
	}

	// Если значение BotToken не было заменено, добавьте его в конец файла
	if !replaced {
		lines = append(lines, newAttrProp)
	}

	// Запись измененных строк обратно в файл .env
	err = os.WriteFile(envFile, []byte(strings.Join(lines, "\n")), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// Percentage - рассчитывает процент от числа
func Percentage(percent, num float64) float64 {
	return (percent / 100) * num
}

/*
// GetCryptoKey - создает ключ шифрования и возвращает его хэш
func GetCryptoKey() string {
	apiToken := strings.Split(config.Cfg.CryptKey, ":")[1]
	hash := MD5(apiToken)

	if len(hash) < 32 {
		for len(hash) < 32 {
			hash += "1"
		}
	}

	return hash[:32]
}
*/

// GetBotAvatarPath - возвращает локальный путь к файлу аватара бота
func GetBotAvatarPath() string {
	//Путь к текущей директории
	pwd, _ := os.Getwd()
	//Путь к папке с profile (полный путь)
	pathDir := filepath.Join(pwd, "profile")
	//Имя файла
	fileName := "r_avatar.jpg"

	fullPath := path.Join(pathDir, fileName)

	return fullPath
}
