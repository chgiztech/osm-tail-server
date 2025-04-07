package http_handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"osm-tail/dto"
	"osm-tail/utils/envconf"
	"osm-tail/utils/validation"
	"path/filepath"
)

func validateBody(c *gin.Context, db *gorm.DB) (dto.CoordinateDto, error) {
	var requestData dto.CoordinateDto

	if err := c.ShouldBindJSON(&requestData); err != nil {
		log.Println("Invalid JSON input, err: ", err)
		return requestData, err
	}

	err := validation.Validate.Struct(requestData)

	if err != nil {
		log.Println("Invalid JSON input, err: ", err)
		return requestData, err
	}

	return requestData, nil
}

func downloadOSM(dto dto.CoordinateDto) error {
	url := fmt.Sprintf(OSM_API_URL, dto.MinLon, dto.MinLat, dto.MaxLon, dto.MaxLat)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	OSM_FILE_PATH := fmt.Sprintf("%s/%s", TMP_DIR, OSM_FILE)

	// Create the tmp directory if it doesn't exist
	file, err := os.Create(OSM_FILE_PATH)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func importToPostGIS() error {
	absPath, err := filepath.Abs(fmt.Sprintf("%s/%s", TMP_DIR, OSM_FILE))

	if err != nil {
		log.Printf("Ошибка при получении абсолютного пути: %v", err)
	}

	// Используем osm2pgsql для импорта OSM-файла в PostGIS
	cmd := exec.Command("osm2pgsql", "-H", envconf.App.PostgreSQL.Host, "-U", envconf.App.PostgreSQL.User, "-d", envconf.App.PostgreSQL.Database, "--create", absPath)

	cmd.Env = append(os.Environ(), "PGPASSWORD="+envconf.App.PostgreSQL.Password)
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Println("Import error:", string(output))
		return err
	}

	log.Println("Import successful:", string(output))

	// Удаление временного файла после успешного импорта
	err = os.Remove(absPath)
	if err != nil {
		log.Printf("Removing file error: %v", err)
	}

	return nil
}
