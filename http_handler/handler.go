package http_handler

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Handler struct {
	Tracer trace.Tracer
	Db     *gorm.DB
}

func (h *Handler) HandleCoordinates(c *gin.Context) {
	ctx := c.Request.Context()
	_, stage1Span := h.Tracer.Start(ctx, "Stage 1 (Validate Body)")
	defer stage1Span.End()

	requestData, err := validateBody(c, h.Db)
	if err != nil {
		log.Printf("Error validating body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error validating body", "details": err.Error()})
		return

	}

	err = downloadOSM(requestData)

	if err != nil {
		log.Printf("Error downloading OSM data: %v", err)
	}

	err = importToPostGIS()

	if err != nil {
		log.Printf("Error importing to PostGIS: %v", err)
	}
}

func (h *Handler) HandleGenerate(c *gin.Context) {
	z := c.Param("z")
	x := c.Param("x")
	y := c.Param("y")

	query := `
		SELECT COALESCE(ST_AsMVT(tile, 'osm_points', 4096, 'geom'), '')
		FROM (
			SELECT
				id,
				ST_AsMVTGeom(
					geom,
					ST_TileEnvelope(?, ?, ?),
					4096, 256, true
				) AS geom,
				tags::TEXT AS tags
			FROM osm_points
			WHERE geom && ST_TileEnvelope(?, ?, ?)
		) AS tile;
    `

	var tile sql.NullString

	err := h.Db.Raw(query, z, x, y, z, x, y).Scan(&tile).Error

	if err != nil {
		log.Printf("Ошибка SQL: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при генерации тайла"})
		return
	}

	if !tile.Valid || tile.String == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tile not found"})
		return
	}

	tileBytes := []byte(tile.String)

	c.Header("Content-Type", "application/x-protobuf")
	c.Data(http.StatusOK, "application/x-protobuf", tileBytes)
	// c.JSON(http.StatusOK, gin.H{"message": "generate"})
}
