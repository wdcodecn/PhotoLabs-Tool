// fileName: util_test.go
package main

import (
	"fmt"
	"testing"
	"time"
)

func TestGetPhotoTakenTime(t *testing.T) {
	tests := []struct {
		filePath         string
		noShotTimeType   int
		expectedError    bool
		expectedTimeFunc func() time.Time
	}{
		{
			filePath:       "D:\\DCIM\\sorted1\\20240928\\P1010586.MOV",
			noShotTimeType: 2,
			expectedError:  false,
			expectedTimeFunc: func() time.Time {
				return time.Date(2013, 7, 7, 16, 36, 17, 0, time.UTC) // 示例 EXIF 时间
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.filePath, func(t *testing.T) {
			got, err := getPhotoTakenTime(tt.filePath, tt.noShotTimeType)
			//got
			fmt.Printf("getPhotoTakenTime: %+v \n", got)

			if (err != nil) != tt.expectedError {
				t.Errorf("getPhotoTakenTime() error = %v, expectedError %v", err, tt.expectedError)
				return
			}

			if !tt.expectedError && !got.Equal(tt.expectedTimeFunc()) {
				t.Errorf("getPhotoTakenTime() = %v, want %v", got, tt.expectedTimeFunc())
			}
		})
	}
}
