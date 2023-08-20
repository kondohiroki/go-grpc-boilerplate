package conv

import (
	"fmt"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/golang-module/carbon/v2"
)

func ConvertStringToInt(name string, strValue string) int {
	number, err := strconv.Atoi(strValue)
	if err != nil {
		panic(fmt.Sprintf("error while converting %s to int %s", name, strValue))
	}

	return number
}

func ConvertStringToInt64(name string, strValue string) int64 {
	number, err := strconv.ParseInt(strValue, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("error while converting %s to int64 %s", name, strValue))
	}

	return number
}

func ConvertStringToPointerInt64(name string, strValue string) *int64 {
	number, err := strconv.ParseInt(strValue, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("error while converting %s to int64 %s", name, strValue))
	}

	return &number
}

func ConvertStringToPointerString(strValue string) *string {
	if strValue == "" {
		return nil
	}
	return &strValue
}

func ConvertStringToFloat64(name string, strValue string) float64 {
	number, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		panic(fmt.Sprintf("error while converting %s to float64 %s", name, strValue))
	}

	return number
}
func ConvertStringToPointerFloat64(name string, strValue string) *float64 {
	number, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		panic(fmt.Sprintf("error while converting %s to float64 %s", name, strValue))
	}

	return &number
}

// format => "d/m/Y"
func ConvertStringUnixTimeToDatetime(name string, strValue string, format string) string {
	number, err := strconv.Atoi(strValue)
	if err != nil {
		panic(fmt.Sprintf("error while converting %s to int %s", name, strValue))
	}

	birthDateFormatted := carbon.CreateFromTimestamp(int64(number)).Format(format)

	return birthDateFormatted
}

func ConvertJsonStringToPointerIntSlice(name string, strValue string) []*int {
	var strSlice []string
	err := sonic.Unmarshal([]byte(strValue), &strSlice)
	if err != nil {
		panic(fmt.Sprintf("error while converting %s to string slice %s", name, strValue))
	}

	var result []*int
	for _, s := range strSlice {
		i, err := strconv.Atoi(s)
		if err != nil {
			panic(fmt.Sprintf("error while converting %s to int slice %s", name, strValue))
		}
		result = append(result, &i)
	}

	return result

}

func ConvertJsonStringToStringSlice(name string, strValue string) []string {
	var result []string
	err := sonic.Unmarshal([]byte(strValue), &result)
	if err != nil {
		panic(fmt.Sprintf("error while converting %s to string slice %s", name, strValue))
	}

	return result
}

func ConvertJsonStringToFloat64Slice(name string, strValue string) []float64 {
	var strSlice []string
	err := sonic.Unmarshal([]byte(strValue), &strSlice)
	if err != nil {
		panic(fmt.Sprintf("error while converting %s to string slice %s", name, strValue))
	}

	var result []float64
	for _, s := range strSlice {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			panic(fmt.Sprintf("error while converting %s to float64 slice %s", name, strValue))
		}
		result = append(result, f)
	}

	return result
}
