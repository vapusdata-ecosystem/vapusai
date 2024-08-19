package utils

import (
	"archive/tar"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	os "os"
	"path"
	filepath "path/filepath"
	"reflect"
	"strconv"
	"strings"

	toml "github.com/pelletier/go-toml/v2"
	dmerrors "github.com/vapusdata-ecosystem/vapusai-studio/internals/errors"
	"gopkg.in/yaml.v2"
)

func AStructToAMap(data interface{}) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, dmerrors.DMError(dmerrors.ErrStruct2Json, err)
	}

	err = json.Unmarshal(bytes, &result)

	if err != nil {
		return nil, dmerrors.DMError(dmerrors.ErrStruct2Json, err)
	}
	return result, nil
}

func StructToMap(data interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, dmerrors.DMError(dmerrors.ErrStruct2Json, err)
	}

	err = json.Unmarshal(bytes, &result)

	if err != nil {
		return nil, dmerrors.DMError(dmerrors.ErrStruct2Json, err)
	}
	return result, nil
}

func AStructToAString(data any) string {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Println("Error converting struct to string: ", err)
		return EMPTYSTR
	}

	return string(bytes)
}

func StrToStruct(data any, resp interface{}) error {
	log.Println("data to be marshelled - ", data)
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Println("Error converting interface to string: ", err)
		return err
	}

	err = json.Unmarshal(bytes, resp)

	if err != nil {
		return dmerrors.DMError(dmerrors.ErrStruct2Json, err)
	}
	return nil
}

func Int32Ptr(x int32) *int32 { return &x }

func Int64Ptr(x int64) *int64 { return &x }

func IntPtr(x int) *int { return &x }

func Bool2Ptr(x bool) *bool { return &x }

func Str2Ptr(x string) *string { return &x }

func GetObjectTypeName(obj any) string {
	return reflect.TypeOf(obj).Name()
}

func GetConfFileType(fileName string) string {
	return strings.Replace(path.Ext(fileName), ".", "", -1)
}

func TrailingSlash(path string, verify, add bool) (bool, string) {
	if verify && strings.HasSuffix(path, "/") {
		return true, path
	}
	if add && !strings.HasSuffix(path, "/") {
		return true, path + "/"
	}
	return false, path
}

func CreateSecretName(opts ...string) string {
	return strings.Join(opts, "::")
}

func IsInt(s string) (bool, int) {
	val, err := strconv.Atoi(s)
	return err == nil, val
}

func IsFloat(s string) (bool, float64) {
	val, err := strconv.ParseFloat(s, 64)
	return err == nil, val
}

func GenerateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}

func CreateTarFile(tarFile string, files2Add []string, fileDest string) error {
	tarFileHandle, err := os.Create(tarFile)
	if err != nil {
		return err
	}
	defer tarFileHandle.Close()

	tw := tar.NewWriter(tarFileHandle)
	defer tw.Close()

	for _, fl := range files2Add {
		file, err := os.Open(fl)
		if err != nil {
			return err
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			return err
		}

		hdr := &tar.Header{
			Name: filepath.Join(fileDest, filepath.Base(file.Name())),
			Mode: 0644,
			Size: fileInfo.Size(),
		}
		log.Println("File to be added to tar - ", filepath.Join(fileDest, filepath.Base(file.Name())))
		err = tw.WriteHeader(hdr)
		if err != nil {
			return err
		}

		_, err = io.Copy(tw, file)
		if err != nil {
			return err
		}
	}

	return err
}

func SetCtxValue(ctx context.Context, key ContextKeys, value any) context.Context {
	return context.WithValue(ctx, key, value)
}

func GetCtxValue(ctx context.Context, key ContextKeys) any {
	return ctx.Value(key)
}

func WriteYAMLFile(data any, fileName string) error {
	bytes, err := yaml.Marshal(data)
	if err != nil {
		return dmerrors.DMError(dmerrors.ErrStruct2Json, err)
	}
	err = os.WriteFile(fileName, bytes, 0644)
	if err != nil {
		return dmerrors.DMError(dmerrors.ErrWriteYAMLFile, err)
	}
	return nil
}

func WriteJSONFile(data any, fileName string) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return dmerrors.DMError(dmerrors.ErrStruct2Json, err)
	}
	err = os.WriteFile(fileName, bytes, 0644)
	if err != nil {
		return dmerrors.DMError(dmerrors.ErrWriteYAMLFile, err)
	}
	return nil
}

func WriteTOMLFile(data any, fileName string) error {
	bytes, err := toml.Marshal(data)
	if err != nil {
		return dmerrors.DMError(dmerrors.ErrStruct2Json, err)
	}
	err = os.WriteFile(fileName, bytes, 0644)
	if err != nil {
		return dmerrors.DMError(dmerrors.ErrWriteYAMLFile, err)
	}
	return nil
}

func ReadFile(fileName string) ([]byte, error) {
	return os.ReadFile(fileName)
}
func GenericUnMarshaler(bytes []byte, result any, format string) error {
	var err error
	switch format {
	case "yaml":
		err = yaml.Unmarshal(bytes, result)
		if err != nil {
			return err
		}
	case "json":
		err = json.Unmarshal(bytes, result)
		if err != nil {
			return err
		}
	case "toml":
		err = toml.Unmarshal(bytes, result)
		if err != nil {
			return err
		}
	default:
		return dmerrors.ErrInvalidArgs
	}
	return err
}

func GenericMarshaler(object any, format string) ([]byte, error) {
	format = strings.ToLower(format)
	switch format {
	case "yaml":
		return yaml.Marshal(object)
	case "json":
		return json.Marshal(object)
	case "toml":
		return toml.Marshal(object)
	default:
		return nil, dmerrors.ErrInvalidArgs
	}
}

func CreateFile(filename, filePath string, data any) error {
	if filePath == "" {
		curPath, err := os.Getwd()
		if err != nil {
			curPath = os.TempDir()
		}
		filePath = curPath
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if _, err := os.Create(filePath); err != nil {
			return err
		}
	}
	fType := GetConfFileType(filename)
	switch strings.ToLower(fType) {
	case "yaml":
		return WriteYAMLFile(data, filepath.Join(filePath, filename))
	case "json":
		return WriteJSONFile(data, filepath.Join(filePath, filename))
	case "toml":
		return WriteTOMLFile(data, filepath.Join(filePath, filename))
	default:
		return dmerrors.ErrInvalidArgs
	}
}

func DeleteFile(filePath string) error {
	return os.Remove(filePath)
}
